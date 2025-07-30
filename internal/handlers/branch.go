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

type BranchHandler struct {
	service service.BranchService
}

func NewBranchHandler(service service.BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

func (h *BranchHandler) Create(c *gin.Context) {
	var req dto.BranchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	branch := domain.Branch{
		Name:         req.Name,
		Client:       req.Client,
		Uniorg:       req.Uniorg,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	}

	id, err := h.service.Create(context.Background(), &branch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
		return
	}

	c.JSON(http.StatusCreated, dto.BranchResponse{
		ID:           id,
		Name:         req.Name,
		Client:       req.Client,
		Uniorg:       req.Uniorg,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	})
}

func (h *BranchHandler) List(c *gin.Context) {
	branchs, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list branchs"})
		return
	}

	response := make([]dto.BranchResponse, 0, len(branchs))
	for _, branch := range branchs {
		response = append(response, dto.BranchResponse{
			ID:           branch.ID,
			Name:         branch.Name,
			Client:       branch.Client,
			Uniorg:       branch.Uniorg,
			Zipcode:      branch.Zipcode,
			State:        branch.State,
			City:         branch.City,
			Neighborhood: branch.Neighborhood,
			Address:      branch.Address,
			Complement:   branch.Complement,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BranchHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	branch, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}

	c.JSON(http.StatusOK, dto.BranchResponse{
		ID:           branch.ID,
		Name:         branch.Name,
		Client:       branch.Client,
		Uniorg:       branch.Uniorg,
		Zipcode:      branch.Zipcode,
		State:        branch.State,
		City:         branch.City,
		Neighborhood: branch.Neighborhood,
		Address:      branch.Address,
		Complement:   branch.Complement,
	})
}

func (h *BranchHandler) FindByUniorg(c *gin.Context) {
	uniorg := c.Param("uniorg")

	branch, err := h.service.FindByUniorg(context.Background(), uniorg)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}

	c.JSON(http.StatusOK, dto.BranchResponse{
		ID:           branch.ID,
		Name:         branch.Name,
		Client:       branch.Client,
		Uniorg:       branch.Uniorg,
		Zipcode:      branch.Zipcode,
		State:        branch.State,
		City:         branch.City,
		Neighborhood: branch.Neighborhood,
		Address:      branch.Address,
		Complement:   branch.Complement,
	})
}

func (h *BranchHandler) GetByClient(c *gin.Context) {
	client := c.Param("client")

	branchs, err := h.service.GetByClient(context.Background(), client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branches by client"})
		return
	}

	response := make([]dto.BranchResponse, 0, len(branchs))
	for _, branch := range branchs {
		response = append(response, dto.BranchResponse{
			ID:           branch.ID,
			Name:         branch.Name,
			Client:       branch.Client,
			Uniorg:       branch.Uniorg,
			Zipcode:      branch.Zipcode,
			State:        branch.State,
			City:         branch.City,
			Neighborhood: branch.Neighborhood,
			Address:      branch.Address,
			Complement:   branch.Complement,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BranchHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var req dto.BranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	branch := domain.Branch{
		ID:           id,
		Name:         req.Name,
		Client:       req.Client,
		Uniorg:       req.Uniorg,
		Zipcode:      req.Zipcode,
		State:        req.State,
		City:         req.City,
		Neighborhood: req.Neighborhood,
		Address:      req.Address,
		Complement:   req.Complement,
	}

	err = h.service.Update(context.Background(), &branch)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch"})
		return
	}

	c.JSON(http.StatusOK, dto.BranchResponse{
		ID:           id,
		Name:         branch.Name,
		Client:       branch.Client,
		Uniorg:       branch.Uniorg,
		Zipcode:      branch.Zipcode,
		State:        branch.State,
		City:         branch.City,
		Neighborhood: branch.Neighborhood,
		Address:      branch.Address,
		Complement:   branch.Complement,
	})
}

func (h *BranchHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
}
