package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetAuthUser(c)
		if user == nil || !user.IsAdmin() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
			return
		}
		c.Next()
	}
}

func IsSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetAuthUser(c)
		if user == nil || !user.IsSuperAdmin() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Accès refusé"})
			return
		}
		c.Next()
	}
}
