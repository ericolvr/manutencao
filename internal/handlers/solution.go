package handlers

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type SolutionHandler struct {
	solutionService service.SolutionService
}

func NewSolutionHandler(solutionService service.SolutionService) *SolutionHandler {
	return &SolutionHandler{
		solutionService: solutionService,
	}
}

func (h *SolutionHandler) Create(c *gin.Context) {
	var req dto.SolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	solution := req.ToSolutionDomain()
	createdSolution, err := h.solutionService.Create(c.Request.Context(), solution)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToSolutionResponse(createdSolution)
	c.JSON(http.StatusCreated, response)
}

func (h *SolutionHandler) List(c *gin.Context) {
	problemIDStr := c.Query("problem_id")

	if problemIDStr != "" {
		problemID, err := strconv.Atoi(problemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
			return
		}

		solutions, err := h.solutionService.GetByProblem(c.Request.Context(), problemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := dto.ToSolutionResponseList(solutions)
		c.JSON(http.StatusOK, response)
		return
	}

	solutions, err := h.solutionService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToSolutionResponseList(solutions)
	c.JSON(http.StatusOK, response)
}

func (h *SolutionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid solution ID"})
		return
	}

	solution, err := h.solutionService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Solution not found"})
		return
	}

	response := dto.ToSolutionResponse(solution)
	c.JSON(http.StatusOK, response)
}

func (h *SolutionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid solution ID"})
		return
	}

	var req dto.SolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	solution := req.ToSolutionDomain()
	updatedSolution, err := h.solutionService.Update(c.Request.Context(), id, solution)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToSolutionResponse(updatedSolution)
	c.JSON(http.StatusOK, response)
}

func (h *SolutionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid solution ID"})
		return
	}

	err = h.solutionService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
