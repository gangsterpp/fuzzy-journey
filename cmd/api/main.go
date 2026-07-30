package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	auth "github.com/gangsterpp/fuzzy-journey/internal/auth"
	database "github.com/gangsterpp/fuzzy-journey/internal/database"
	middleware "github.com/gangsterpp/fuzzy-journey/internal/middleware"
	"github.com/gangsterpp/fuzzy-journey/internal/response"
	"github.com/gangsterpp/fuzzy-journey/internal/token"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	JWTSecret := os.Getenv("JWT_SECRET")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if JWTSecret == "" {
		log.Fatal("JWTSecret is required")
	}

	jwtManager := token.NewManager(
		JWTSecret,
		15*time.Minute,
	)

	db, err := database.NewPgPool(ctx, databaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(middleware.DatabaseMiddleware(db))

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Use(middleware.SetupCors())
	router.Use(middleware.SecurityHeaders())

	apiV1 := router.Group("api/v1")

	authService := auth.CreateAuthService(auth.NewAuthRepository(db), jwtManager)
	authHandler := auth.CreateAuthHandler(authService)

	if authHandler != nil {
		apiV1.POST("/register", authHandler.Register)
		apiV1.POST("/login", authHandler.Login)
		apiV1.DELETE("/login", authHandler.Delete)

	}
	apiV1.GET("/me", func(c *gin.Context) {
		value, exists := c.Get(middleware.UserIDKey)
		if !exists {
			response.Fail(
				c,
				http.StatusUnauthorized,
				response.CodeUnauthorized,
				response.ErrInvalidCredentials.Error(),
			)
			return
		}
		_, ok := value.(string)
		if !ok {
			response.Fail(
				c,
				http.StatusInternalServerError,
				response.CodeInternal,
				response.ErrInvalidCredentials.Error(),
			)
			return
		}
		response.OK(c, "OK")
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
