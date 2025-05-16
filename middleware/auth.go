package middlewares

import (
	"manufacture_API/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleBasedAuth(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validate token and extract claims
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check role authorization
		authorized := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				authorized = true
				break
			}
		}
		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Store user info in context for downstream handlers
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next() // Continue to next handler
	}
}
