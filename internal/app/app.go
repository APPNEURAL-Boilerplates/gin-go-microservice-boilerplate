package app

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice/internal/common"
	"github.com/example/gin-go-microservice/internal/config"
	"github.com/example/gin-go-microservice/internal/middleware"
	"github.com/example/gin-go-microservice/internal/modules/health"
	"github.com/example/gin-go-microservice/internal/modules/items"
	"github.com/example/gin-go-microservice/internal/modules/root"
)

func NewRouter(cfg config.Config, logger *slog.Logger) *gin.Engine {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	router := gin.New()
	router.HandleMethodNotAllowed = true

	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		logger.Warn("failed to configure trusted proxies", "error", err)
	}

	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.SecurityHeaders())

	router.NoRoute(func(c *gin.Context) {
		common.Error(c, common.NotFound("route not found"))
	})

	router.NoMethod(func(c *gin.Context) {
		common.Error(c, common.MethodNotAllowed("method not allowed"))
	})

	rootHandler := root.NewHandler(cfg)
	healthHandler := health.NewHandler(cfg)

	itemRepository := items.NewMemoryRepository()
	itemService := items.NewService(itemRepository)
	itemHandler := items.NewHandler(itemService)

	router.GET("/", rootHandler.Get)

	v1 := router.Group("/api/v1")
	v1.GET("/health", healthHandler.Health)
	v1.GET("/ready", healthHandler.Ready)

	items.RegisterRoutes(v1.Group("/items"), itemHandler)

	return router
}
