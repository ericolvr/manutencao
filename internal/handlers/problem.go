package handlers

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type ProblemHandler struct {
	problemService service.ProblemService
}

func NewProblemHandler(problemService service.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		problemService: problemService,
	}
}

func (h *ProblemHandler) Create(c *gin.Context) {
	var req dto.ProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	problem := req.ToProblemDomain()
	createdProblem, err := h.problemService.Create(c.Request.Context(), problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProblemResponse(createdProblem)
	c.JSON(http.StatusCreated, response)
}

func (h *ProblemHandler) List(c *gin.Context) {
	problems, err := h.problemService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProblemResponseList(problems)
	c.JSON(http.StatusOK, response)
}

func (h *ProblemHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := h.problemService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	response := dto.ToProblemResponse(problem)
	c.JSON(http.StatusOK, response)
}

func (h *ProblemHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	var req dto.ProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	problem := req.ToProblemDomain()
	updatedProblem, err := h.problemService.Update(c.Request.Context(), id, problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProblemResponse(updatedProblem)
	c.JSON(http.StatusOK, response)
}

func (h *ProblemHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	err = h.problemService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
