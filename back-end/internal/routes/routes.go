package routes

import (
	"entry-project/back-end/internal/handler"
	"entry-project/back-end/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	r.Use(middleware.CORsMiddleware())
	api := r.Group("/api/v1")
	{
		api.POST("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Xin chao! Day la API")
		})
		api.POST("/login", func(c *gin.Context) {
			userHandler.Login(c)
		})
		api.POST("/logout", middleware.AuthMiddleware(), func(c *gin.Context) {
			userHandler.Logout(c)
		})
		api.POST("/register", func(c *gin.Context) {
			userHandler.Register(c)
		})
		api.POST("/refresh", func(c *gin.Context) {
			userHandler.Refresh(c)
		})
		api.GET("/profile/:username", middleware.AuthMiddleware(), func(c *gin.Context) {
			userHandler.GetUser(c)
		})
		api.PUT("/profile/:username", middleware.AuthMiddleware(), func(c *gin.Context) {
			userHandler.UpdateUser(c)
		})
		api.GET("/me", middleware.AuthMiddleware(), func(c *gin.Context) {
			username, _ := c.Get("username")
			if username != "" {
				c.JSON(http.StatusOK, gin.H{
					"username": username,
				})
			}
		})

	}
}
