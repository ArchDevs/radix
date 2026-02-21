package main

import (
	"github.com/ArchDevs/radix/internal/handler"
	"github.com/ArchDevs/radix/internal/middleware"
	"github.com/ArchDevs/radix/internal/repository"
	"github.com/ArchDevs/radix/internal/service"
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

	// Websocket
	hub := wsocket.NewHub()

	// Repositories
	userRepo := repository.NewUserRepository(app.db)
	challengeRepo := repository.NewChallengeRepository(app.db)

	// Services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userSvc)
	challengeSvc := service.NewChallengeService(challengeRepo, userSvc)
	jwtSvc := service.NewJWTService(app.config.Security.JwtSecret, app.config.Security.JwtTTL)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	challengeHandler := handler.NewChallengeHandler(challengeSvc, jwtSvc)
	wsHandler := handler.NewWsHandler(hub, jwtSvc)

	// Routes
	v1.POST("/auth/register", authHandler.Register)
	v1.GET("/challenge", challengeHandler.CreateChallenge)
	v1.POST("/challenge/verify", challengeHandler.Verify)
	v1.GET("/ws", wsHandler.Handle)
}
