package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daung-digital/location-api/config"
	"github.com/daung-digital/location-api/internal/database"
	"github.com/daung-digital/location-api/internal/handler"
	"github.com/daung-digital/location-api/internal/middleware"
	"github.com/daung-digital/location-api/internal/repository"
	"github.com/daung-digital/location-api/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting location API server",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize database
	pool, err := database.NewPool(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Run migrations
	if err := database.RunMigrations(context.Background(), pool); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize repository
	locationRepo := repository.NewPostgresLocationRepository(pool)
	routePlanRepo := repository.NewPostgresRoutePlanRepository(pool)

	// Initialize service
	locationService := service.NewLocationService(locationRepo)
	routePlanService := service.NewRoutePlanService(routePlanRepo)

	// Initialize handlers
	locationHandler := handler.NewLocationHandler(locationService)
	routePlanHandler := handler.NewRoutePlanHandler(routePlanService)

	// Initialize Gin router
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Routes
	v1 := router.Group("/v1")
	{
		locations := v1.Group("/locations")
		{
			locations.POST("", locationHandler.CreateLocation)
			locations.GET("", locationHandler.ListLocations)
			locations.GET("/nearby", locationHandler.FindNearby)
			locations.GET("/:id", locationHandler.GetLocation)
			locations.PATCH("/:id", locationHandler.UpdateLocation)
			locations.DELETE("/:id", locationHandler.DeleteLocation)
		}

		routePlans := v1.Group("/route-plans")
		{
			routePlans.POST("", routePlanHandler.CreateRoutePlan)
			routePlans.GET("", routePlanHandler.ListRoutePlans)
			routePlans.GET("/:id", routePlanHandler.GetRoutePlan)
			routePlans.PATCH("/:id", routePlanHandler.UpdateRoutePlan)
			routePlans.DELETE("/:id", routePlanHandler.DeleteRoutePlan)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
