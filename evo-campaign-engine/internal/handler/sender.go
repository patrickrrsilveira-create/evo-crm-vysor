package handler

import (
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SenderHandler struct {
	repo *repository.SenderRepo
}

func NewSenderHandler(repo *repository.SenderRepo) *SenderHandler {
	return &SenderHandler{repo: repo}
}

func (h *SenderHandler) RegisterRoutes(r gin.IRouter) {
	g := r.Group("/senders")
	{
		g.POST("", h.Upsert)
		g.GET("/:id", h.GetByID)
	}
}

func (h *SenderHandler) Upsert(c *gin.Context) {
	var s domain.SenderInstance
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Upsert(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *SenderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}
