package handler

import (
	"entry-project/back-end/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
}

type UpdateRequest struct {
	Nickname  string `json:"nickname" binding:"required"`
	AvatarUrl string `json:"avatarUrl" binding:"required"`
}

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
)

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	accessToken, refreshToken, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	accessMaxAge := int((15 * time.Minute).Seconds())
	refreshMaxAge := int((30 * 24 * time.Hour).Seconds())

	domain := ""

	c.SetCookie(accessCookieName, accessToken, accessMaxAge, "/", domain, false, true)
	c.SetCookie(refreshCookieName, refreshToken, refreshMaxAge, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
	})

}

// Refresh Token
func (h *UserHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil || refresh == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	newAccess, newRefresh, err := h.userService.RefreshTokens(c.Request.Context(), refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessMaxAge := int((15 * time.Minute).Seconds())
	refreshMaxAge := int((30 * 24 * time.Hour).Seconds())
	domain := ""

	c.SetCookie(accessCookieName, newAccess, accessMaxAge, "/", domain, false, true)
	c.SetCookie(refreshCookieName, newRefresh, refreshMaxAge, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "refreshed"})
}

// Log out
func (h *UserHandler) Logout(c *gin.Context) {
	type LogoutRequest struct {
		Username string
	}
	var logoutRequest LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid username",
		})
		return
	}
	usernameAny, _ := c.Get("username")
	username, _ := usernameAny.(string)

	if logoutRequest.Username != username {
		c.JSON(400, gin.H{
			"error": "invalid username",
		})
	}

	_ = h.userService.Logout(c.Request.Context(), username)

	// clear cookie phía client
	domain := ""
	c.SetCookie(accessCookieName, "", -1, "/", domain, false, true)
	c.SetCookie(refreshCookieName, "", -1, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Register
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	regInput := service.RegisterInput{
		Username:  req.Username,
		Password:  req.Password,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
	}

	if err := h.userService.Register(regInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "register success"})
}

// Get User
func (h *UserHandler) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
	})
}

// Update User
func (h *UserHandler) UpdateUser(c *gin.Context) {
	username := c.Param("username")
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.userService.UpdateUser(c.Request.Context(), username, req.Nickname, req.AvatarUrl); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update successful",
	})
}
