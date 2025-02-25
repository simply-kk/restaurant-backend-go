package middleware

import (
	"fmt"
	"golang-restaurant-management/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authentication Middleware
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization token provided"})
			return
		}

		// Validate Token
		claims, err := helpers.ValidateToken(clientToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			return
		}

		// Store Claims in Context
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)   
		c.Set("uid", claims.Uid)

		// Continue to next handler
		c.Next()
	}
}
