package main

import (
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/ArchDevs/radix/internal/config"
	"github.com/ArchDevs/radix/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

type application struct {
	config *config.Config

	db *database.DB

	logger *slog.Logger
	router *gin.Engine

	wg sync.WaitGroup
}

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := database.New(cfg.DB.DataSource)
	if err != nil {
		return err
	}

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		router: gin.Default(),
	}

	// Configure CORS to allow all origins (adjust for production)
	app.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	app.routes()

	return app.serve()
}
