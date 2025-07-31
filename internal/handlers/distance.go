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

type DistanceHandler struct {
	service service.DistanceService
}

func NewDistanceHandler(service service.DistanceService) *DistanceHandler {
	return &DistanceHandler{service: service}
}

func (h *DistanceHandler) Create(c *gin.Context) {
	var req dto.DistanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	distance := domain.Distance{
		Distance:     req.Distance,
		TicketNumber: req.TicketNumber,
		ProviderId:   req.ProviderId,
		ProviderName: req.ProviderName,
	}

	id, err := h.service.Create(context.Background(), &distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create distance"})
		return
	}

	c.JSON(http.StatusCreated, dto.DistanceResponse{
		ID:           id,
		Distance:     req.Distance,
		TicketNumber: req.TicketNumber,
		ProviderId:   req.ProviderId,
		ProviderName: req.ProviderName,
	})
}

func (h *DistanceHandler) List(c *gin.Context) {
	distances, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list distances"})
		return
	}

	response := make([]dto.DistanceResponse, 0, len(distances))
	for _, distance := range distances {
		response = append(response, dto.DistanceResponse{
			ID:           distance.ID,
			Distance:     distance.Distance,
			TicketNumber: distance.TicketNumber,
			ProviderId:   distance.ProviderId,
			ProviderName: distance.ProviderName,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *DistanceHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid distance ID"})
		return
	}

	distance, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrDistanceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Distance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get distance"})
		return
	}

	c.JSON(http.StatusOK, dto.DistanceResponse{
		ID:           distance.ID,
		Distance:     distance.Distance,
		TicketNumber: distance.TicketNumber,
		ProviderId:   distance.ProviderId,
		ProviderName: distance.ProviderName,
	})
}

func (h *DistanceHandler) FindByNumber(c *gin.Context) {
	number := c.Param("number")

	distance, err := h.service.FindByNumber(context.Background(), number)
	if err != nil {
		if err == repository.ErrDistanceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Distance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get distance"})
		return
	}

	c.JSON(http.StatusOK, dto.DistanceResponse{
		ID:           distance.ID,
		Distance:     distance.Distance,
		TicketNumber: distance.TicketNumber,
		ProviderId:   distance.ProviderId,
		ProviderName: distance.ProviderName,
	})
}

func (h *DistanceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid distance ID"})
		return
	}

	var req dto.DistanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	distance := domain.Distance{
		ID:           id,
		Distance:     req.Distance,
		TicketNumber: req.TicketNumber,
		ProviderId:   req.ProviderId,
		ProviderName: req.ProviderName,
	}

	err = h.service.Update(context.Background(), &distance)
	if err != nil {
		if err == repository.ErrDistanceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Distance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update distance"})
		return
	}

	c.JSON(http.StatusOK, dto.DistanceResponse{
		ID:           id,
		Distance:     distance.Distance,
		TicketNumber: distance.TicketNumber,
		ProviderId:   distance.ProviderId,
		ProviderName: distance.ProviderName,
	})
}

func (h *DistanceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid distance ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrDistanceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Distance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete distance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Distance deleted successfully"})
}
