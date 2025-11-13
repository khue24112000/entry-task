package middleware

import (
	"entry-project/back-end/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookieName = "access_token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(AccessTokenCookieName)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"erorr": "missing access token",
			})
		}
		claims, err := utils.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}

		if c.Param("username") != "" && c.Param("username") != claims.Username {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
