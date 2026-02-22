package main

import (
	"github.com/ArchDevs/radix/internal/auth"
	"github.com/ArchDevs/radix/internal/challenge"
	"github.com/ArchDevs/radix/internal/message"
	"github.com/ArchDevs/radix/internal/middleware"
	"github.com/ArchDevs/radix/internal/service"
	"github.com/ArchDevs/radix/internal/user"
	"github.com/ArchDevs/radix/internal/wsocket"
	"golang.org/x/time/rate"
)

func (app *application) routes() {
	rateLimit := rate.Limit(app.config.Server.RateLimit)
	if rateLimit <= 0 {
		rateLimit = rate.Limit(5)
	}

	burst := app.config.Server.RateBurst
	if burst <= 0 {
		burst = 2
	}

	globalLimiter := middleware.NewRateLimiter(rateLimit, burst)

	v1 := app.router.Group("/v1")

	v1.Use(globalLimiter.Middleware())

	// Service
	jwtSvc := service.NewJWTService(app.config.Security.JwtSecret, app.config.Security.JwtTTL)

	// User
	userRepo := user.NewUserRepository(app.db)
	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	// Auth
	authSvc := auth.NewAuthService(userSvc)
	authHandler := auth.NewAuthHandler(authSvc)

	// Challenge
	challengeRepo := challenge.NewChallengeRepository(app.db)
	challengeSvc := challenge.NewChallengeService(challengeRepo, userSvc)
	challengeHandler := challenge.NewChallengeHandler(challengeSvc, jwtSvc)

	// Message
	msgRepo := message.NewMessageRepository(app.db)
	msgSvc := message.NewMessageService(msgRepo)
	msgHandler := message.NewHandler(msgSvc)

	// Websocket
	hub := wsocket.NewHub()
	wsHandler := wsocket.NewWsHandler(hub, jwtSvc, msgSvc)

	// Public routes
	v1.POST("/auth/register", authHandler.Register)
	v1.GET("/challenge", challengeHandler.CreateChallenge)
	v1.POST("/challenge/verify", challengeHandler.Verify)
	v1.GET("/ws", wsHandler.Handle)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(auth.Auth(jwtSvc))
	{
		// User
		protected.GET("/me", userHandler.Me)
		protected.POST("/me/username", userHandler.SetUsername)
		protected.GET("/search", userHandler.Search)

		// Messages
		protected.GET("/messages", msgHandler.GetHistory) // ?with=rad:xxx
		protected.GET("/messages/undelivered", msgHandler.GetUndelivered)
	}
}
