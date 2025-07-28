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

type ProviderHandler struct {
	service service.ProviderService
}

func NewProviderHandler(service service.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

func (h *ProviderHandler) Create(c *gin.Context) {
	var req dto.ProviderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	provider := domain.Provider{
		Name:         req.Name,
		Mobile:       req.Mobile,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	}

	id, err := h.service.Create(context.Background(), &provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider"})
		return
	}

	c.JSON(http.StatusCreated, dto.ProviderResponse{
		ID:           strconv.Itoa(id),
		Name:         req.Name,
		Mobile:       req.Mobile,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	})
}

func (h *ProviderHandler) List(c *gin.Context) {
	providers, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list providers"})
		return
	}

	response := make([]dto.ProviderResponse, 0, len(providers))
	for _, provider := range providers {
		response = append(response, dto.ProviderResponse{
			ID:           strconv.Itoa(provider.ID),
			Name:         provider.Name,
			Mobile:       provider.Mobile,
			Zipcode:      provider.Zipcode,
			State:        provider.State,
			City:         provider.City,
			Neighborhood: provider.Neighborhood,
			Address:      provider.Address,
			Complement:   provider.Complement,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProviderHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	provider, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider"})
		return
	}

	c.JSON(http.StatusOK, dto.ProviderResponse{
		ID:           strconv.Itoa(provider.ID),
		Name:         provider.Name,
		Mobile:       provider.Mobile,
		Zipcode:      provider.Zipcode,
		State:        provider.State,
		City:         provider.City,
		Neighborhood: provider.Neighborhood,
		Address:      provider.Address,
		Complement:   provider.Complement,
	})
}

func (h *ProviderHandler) FindByName(c *gin.Context) {
	name := c.Param("name")

	provider, err := h.service.FindByName(context.Background(), name)
	if err != nil {
		if err == repository.ErrProviderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provider"})
		return
	}

	c.JSON(http.StatusOK, dto.ProviderResponse{
		ID:           strconv.Itoa(provider.ID),
		Name:         provider.Name,
		Mobile:       provider.Mobile,
		Zipcode:      provider.Zipcode,
		State:        provider.State,
		City:         provider.City,
		Neighborhood: provider.Neighborhood,
		Address:      provider.Address,
		Complement:   provider.Complement,
	})
}

func (h *ProviderHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	var req dto.ProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	provider := domain.Provider{
		ID:           id,
		Name:         req.Name,
		Mobile:       req.Mobile,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	}

	err = h.service.Update(context.Background(), &provider)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update provider"})
		return
	}

	c.JSON(http.StatusOK, dto.ProviderResponse{
		ID:           strconv.Itoa(id),
		Name:         provider.Name,
		Mobile:       provider.Mobile,
		Zipcode:      provider.Zipcode,
		State:        provider.State,
		City:         provider.City,
		Neighborhood: provider.Neighborhood,
		Address:      provider.Address,
		Complement:   provider.Complement,
	})
}

func (h *ProviderHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted successfully"})
}
