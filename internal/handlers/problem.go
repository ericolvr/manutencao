package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
)

type ProblemHandler struct {
	problemService service.ProblemService
}

func NewProblemHandler(problemService service.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		problemService: problemService,
	}
}

// @Summary Create a new problem
// @Description Create a new problem in the system
// @Tags problems
// @Accept json
// @Produce json
// @Param problem body dto.ProblemRequest true "Problem data"
// @Success 201 {object} dto.ProblemResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/problems [post]
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

// @Summary Get all problems
// @Description Get a list of all problems
// @Tags problems
// @Produce json
// @Success 200 {array} dto.ProblemResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/problems [get]
func (h *ProblemHandler) List(c *gin.Context) {
	problems, err := h.problemService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProblemResponseList(problems)
	c.JSON(http.StatusOK, response)
}

// @Summary Get problem by ID
// @Description Get a specific problem by its ID
// @Tags problems
// @Produce json
// @Param id path int true "Problem ID"
// @Success 200 {object} dto.ProblemResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/problems/{id} [get]
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

// @Summary Update problem
// @Description Update an existing problem
// @Tags problems
// @Accept json
// @Produce json
// @Param id path int true "Problem ID"
// @Param problem body dto.ProblemRequest true "Problem data"
// @Success 200 {object} dto.ProblemResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/problems/{id} [put]
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

// @Summary Delete problem
// @Description Delete a problem by ID
// @Tags problems
// @Param id path int true "Problem ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/problems/{id} [delete]
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
