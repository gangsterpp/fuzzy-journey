package middleware

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gangsterpp/fuzzy-journey/internal/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		c.Next()
	}
}

func SetupCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func DatabaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

const UserIDKey = "userID"

func Auth(tokenManager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "missing authorization token"},
			)
			return
		}

		userID, err := tokenManager.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid authorization token"},
			)
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
