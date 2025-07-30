package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/repository"
	"github.com/ericolvr/maintenance-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type CostHandler struct {
	service service.CostService
}

func NewCostHandler(service service.CostService) *CostHandler {
	return &CostHandler{service: service}
}

func (h *CostHandler) List(c *gin.Context) {
	costs, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list costs"})
		return
	}

	response := make([]dto.CostResponse, 0, len(costs))
	for _, cost := range costs {
		response = append(response, dto.CostResponse{
			ID:           cost.ID,
			ValuePerKm:   cost.ValuePerKm,
			InitialValue: cost.InitialValue,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *CostHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost ID"})
		return
	}

	cost, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cost"})
		return
	}

	c.JSON(http.StatusOK, dto.CostResponse{
		ID:           cost.ID,
		ValuePerKm:   cost.ValuePerKm,
		InitialValue: cost.InitialValue,
	})
}

func (h *CostHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost ID"})
		return
	}

	var req dto.CostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cost := domain.Cost{
		ID:           id,
		ValuePerKm:   req.ValuePerKm,
		InitialValue: req.InitialValue,
	}

	err = h.service.Update(context.Background(), &cost)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cost"})
		return
	}

	c.JSON(http.StatusOK, dto.CostResponse{
		ID:           id,
		ValuePerKm:   cost.ValuePerKm,
		InitialValue: cost.InitialValue,
	})
}

func (h *CostHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cost"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cost deleted successfully"})
}
