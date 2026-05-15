package items

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice-boilerplate/internal/common"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.GET("", handler.List)
	router.POST("", handler.Create)
	router.GET("/:id", handler.Get)
}

func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		common.Error(c, err)
		return
	}

	common.JSON(c, http.StatusOK, gin.H{
		"items": items,
	})
}

func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		common.Error(c, err)
		return
	}

	common.JSON(c, http.StatusOK, gin.H{
		"item": item,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.Error(c, common.Validation("invalid request body", err))
		return
	}

	item, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		common.Error(c, err)
		return
	}

	common.JSON(c, http.StatusCreated, gin.H{
		"item": item,
	})
}
