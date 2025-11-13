package service

import (
	"context"
	"encoding/json"
	"entry-project/back-end/internal/model"
	"entry-project/back-end/internal/repository"
	"entry-project/back-end/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserExisted        = errors.New("user already exist")
)

type RegisterInput struct {
	Username  string
	Nickname  string
	Password  string
	AvatarUrl string
}

type UserService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redis *redis.Client) *UserService {
	return &UserService{userRepo: userRepo, redis: redis}
}

// Login
func (s *UserService) Login(ctx context.Context, username, password string) (string, string, error) {
	// user, err := s.GetUserByUsername(ctx, username)
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", ErrUserNotFound
	}
	fmt.Println(user.Username)

	check := utils.VerifyPassword(password, user.Login.Salt, user.Login.PasswordHash)
	if !check {
		return "", "", ErrInvalidCredentials
	}

	access, _, err := utils.NewAccessToken(username)

	if err != nil {
		return "", "", err
	}

	refresh, refreshExp, err := utils.NewRefreshToken(username)
	if err != nil {
		return "", "", err
	}

	refreshKey := fmt.Sprintf("refresh:%s", user.Username)
	if err := s.redis.Set(ctx, refreshKey, refresh, time.Until(refreshExp)).Err(); err != nil {
		return "", "", err
	}

	return access, refresh, nil

}

// Refresh Token
func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New(("invalid refresh token"))
	}

	refreshKey := fmt.Sprintf("refresh:%s", claims.Username)
	stored, err := s.redis.Get(ctx, refreshKey).Result()
	if err != nil || stored != refreshToken {
		return "", "", errors.New("refresh token revoked or not found")
	}

	newAccess, _, err := utils.NewAccessToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	newRefresh, refreshExp, err := utils.NewRefreshToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	// Update redis
	if err := s.redis.Set(ctx, refreshKey, newRefresh, time.Until(refreshExp)).Err(); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil

}

// Logout Service
func (s *UserService) Logout(ctx context.Context, username string) error {
	refreshKey := fmt.Sprintf("refresh:%s", username)
	err := s.redis.Del(ctx, refreshKey).Err()
	return err
}

// Register
func (s *UserService) Register(regInput RegisterInput) error {
	existing, _ := s.userRepo.FindByUsername(regInput.Username)
	if existing != nil && existing.ID != 0 {
		return ErrUserExisted
	}

	salt := utils.GenerateSalt()
	hash := utils.GeneratePassword(regInput.Password, salt)

	user := model.User{
		Username:  regInput.Username,
		Nickname:  regInput.Nickname,
		AvatarURL: regInput.AvatarUrl,
		Login: model.Login{
			PasswordHash: string(hash),
			Salt:         string(salt),
		},
	}

	return s.userRepo.Create(&user)

}

const userCacheTTL = 5 * time.Minute

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:username:%s", username)
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var user model.User

			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				// cache hit
				return &user, nil
			}
		}
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if data, err := json.Marshal(user); err == nil {
			_ = s.redis.Set(ctx, cacheKey, data, userCacheTTL).Err()
		}
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, username, nickname, avatarURL string) error {
	if err := s.userRepo.UpdateUser(ctx, username, nickname, avatarURL); err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("user:username:%s", username)
	if s.redis != nil {
		_ = s.redis.Del(ctx, cacheKey).Err()
	}

	return nil
}
