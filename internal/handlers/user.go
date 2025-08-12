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

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user := domain.User{
		Name:     req.Name,
		Mobile:   req.Mobile,
		Password: req.Password,
		Role:     req.Role,
		Status:   req.Status,
	}

	id, err := h.service.Create(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:     id,
		Name:   req.Name,
		Mobile: req.Mobile,
		Role:   req.Role,
		Status: req.Status,
	})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Mobile: user.Mobile,
			Role:   user.Role,
			Status: user.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Mobile: user.Mobile,
		Role:   user.Role,
		Status: user.Status,
	})
}

func (h *UserHandler) FindByname(c *gin.Context) {
	name := c.Param("name")

	users, err := h.service.FindByName(context.Background(), name)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := users[0]
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Mobile: user.Mobile,
		Role:   user.Role,
		Status: user.Status,
	})
}

func (h *UserHandler) FindByMobile(c *gin.Context) {
	mobile := c.Param("mobile")

	users, err := h.service.FindByMobile(context.Background(), mobile)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User or mobile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := users[0]
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Mobile: user.Mobile,
		Role:   user.Role,
		Status: user.Status,
	})
}

func (h *UserHandler) FindUsersToTicket(c *gin.Context) {
	users, err := h.service.FindUsersToTicket(context.Background())
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User or role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Mobile:   req.Mobile,
		Password: req.Password, // Adicionando o campo Password
		Role:     req.Role,
		Status:   req.Status,
	}

	err = h.service.Update(context.Background(), &user)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:     id,
		Name:   user.Name,
		Mobile: user.Mobile,
		Role:   user.Role,
		Status: user.Status,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.Delete(context.Background(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) Authenticate(c *gin.Context) {
	var loginData dto.UserLogin

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	authResponse, err := h.service.Authenticate(context.Background(), loginData.Mobile, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}
