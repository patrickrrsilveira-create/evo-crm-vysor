package handler

import (
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ThrottleHandler struct {
	repo *repository.ThrottleRepo
}

func NewThrottleHandler(repo *repository.ThrottleRepo) *ThrottleHandler {
	return &ThrottleHandler{repo: repo}
}

func (h *ThrottleHandler) RegisterRoutes(r gin.IRouter) {
	g := r.Group("/throttle-profiles")
	{
		g.POST("", h.Create)
		g.GET("", h.List)
		g.GET("/:id", h.GetByID)
		g.PUT("/:id", h.Update)
	}
}

func (h *ThrottleHandler) Create(c *gin.Context) {
	var p domain.ThrottleProfile
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ThrottleHandler) List(c *gin.Context) {
	channelType := c.DefaultQuery("channel_type", "whatsapp")
	profiles, err := h.repo.ListByChannel(c.Request.Context(), channelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profiles})
}

func (h *ThrottleHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ThrottleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existing, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existing.ID = id

	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}
