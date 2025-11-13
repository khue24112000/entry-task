package repository

import (
	"context"
	"entry-project/back-end/internal/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	UpdateUser(ctx context.Context, username, nickname, avatarUrl string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.
		Preload("Login").
		Where("username = ?", username).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, username, nickname, avatarURL string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("username= ?", username).Updates(model.User{Nickname: nickname, AvatarURL: avatarURL})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("khong co ban luu duoc cap nhat")
	}

	return nil
}
