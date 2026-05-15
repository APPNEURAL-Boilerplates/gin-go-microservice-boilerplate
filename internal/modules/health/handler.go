package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice-boilerplate/internal/common"
	"github.com/example/gin-go-microservice-boilerplate/internal/config"
)

type Handler struct {
	cfg       config.Config
	startedAt time.Time
}

func NewHandler(cfg config.Config) *Handler {
	return &Handler{
		cfg:       cfg,
		startedAt: time.Now().UTC(),
	}
}

func (h *Handler) Health(c *gin.Context) {
	common.JSON(c, http.StatusOK, gin.H{
		"status":         "healthy",
		"service":        h.cfg.ServiceName,
		"environment":    h.cfg.Environment,
		"uptime_seconds": int(time.Since(h.startedAt).Seconds()),
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Ready(c *gin.Context) {
	common.JSON(c, http.StatusOK, gin.H{
		"status":  "ready",
		"service": h.cfg.ServiceName,
	})
}
