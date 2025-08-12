package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/repository"
	"github.com/ericolvr/maintenance-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type SlaHandler struct {
	service service.SlaService
}

func NewSlaHandler(service service.SlaService) *SlaHandler {
	return &SlaHandler{service: service}
}

func (h *SlaHandler) Create(c *gin.Context) {
	var req dto.SlaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sla := domain.Sla{
		ClientID: req.ClientID,
		Priority: req.Priority,
		Hours:    req.Hours,
	}

	id, err := h.service.Create(context.Background(), &sla)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
		return
	}

	c.JSON(http.StatusCreated, dto.SlaResponse{
		ID:       id,
		ClientID: req.ClientID,
		Priority: req.Priority,
		Hours:    req.Hours,
	})
}

func (h *SlaHandler) List(c *gin.Context) {
	slas, err := h.service.ListWithClientNames(context.Background())
	if err != nil {
		// Log do erro específico para debug
		fmt.Printf("Error in SLA List handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list slas"})
		return
	}

	c.JSON(http.StatusOK, slas)
}

func (h *SlaHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sla ID"})
		return
	}

	sla, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sla not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sla"})
		return
	}

	c.JSON(http.StatusOK, dto.SlaResponse{
		ID:       sla.ID,
		ClientID: sla.ClientID,
		Priority: sla.Priority,
		Hours:    sla.Hours,
	})
}

func (h *SlaHandler) FindByParams(c *gin.Context) {
	client := c.Param("client")
	client_id, err := strconv.Atoi(client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	priorityStr := c.Param("priority")
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority ID"})
		return
	}

	sla, err := h.service.FindByParams(context.Background(), client_id, priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sla by params"})
		return
	}

	c.JSON(http.StatusOK, sla)
}

func (h *SlaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sla ID"})
		return
	}

	var req dto.SlaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sla := domain.Sla{
		ID:       id,
		ClientID: req.ClientID,
		Priority: req.Priority,
		Hours:    req.Hours,
	}

	err = h.service.Update(context.Background(), &sla)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sla not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch"})
		return
	}

	c.JSON(http.StatusOK, dto.SlaResponse{
		ID:       id,
		ClientID: req.ClientID,
		Priority: req.Priority,
		Hours:    req.Hours,
	})
}

func (h *SlaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sla ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sla not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sla"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sla deleted successfully"})
}
