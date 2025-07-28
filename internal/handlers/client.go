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

type ClientHandler struct {
	service service.ClientService
}

func NewClientHandler(service service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.ClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client := domain.Client{
		Name: req.Name,
	}

	id, err := h.service.Create(context.Background(), &client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, dto.ClientResponse{
		ID:   strconv.Itoa(id),
		Name: client.Name,
	})
}

func (h *ClientHandler) List(c *gin.Context) {
	clients, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list clients"})
		return
	}

	response := make([]dto.ClientResponse, 0, len(clients))
	for _, client := range clients {
		response = append(response, dto.ClientResponse{
			ID:   strconv.Itoa(client.ID),
			Name: client.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	client, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client"})
		return
	}

	c.JSON(http.StatusOK, dto.ClientResponse{
		ID:   strconv.Itoa(client.ID),
		Name: client.Name,
	})
}

func (h *ClientHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var req dto.ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client := domain.Client{
		ID:   id,
		Name: req.Name,
	}

	err = h.service.Update(context.Background(), &client)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}

	c.JSON(http.StatusOK, dto.ClientResponse{
		ID:   strconv.Itoa(id),
		Name: client.Name,
	})
}

func (h *ClientHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
