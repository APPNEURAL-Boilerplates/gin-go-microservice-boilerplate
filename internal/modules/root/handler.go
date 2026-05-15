package root

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice/internal/common"
	"github.com/example/gin-go-microservice/internal/config"
)

type Handler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) Get(c *gin.Context) {
	common.JSON(c, http.StatusOK, gin.H{
		"service":     h.cfg.ServiceName,
		"environment": h.cfg.Environment,
		"version":     "0.1.0",
	})
}
