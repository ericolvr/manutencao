package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
)

type TicketSolutionHandler struct {
	ticketService service.TicketService
}

func NewTicketSolutionHandler(ticketService service.TicketService) *TicketSolutionHandler {
	return &TicketSolutionHandler{
		ticketService: ticketService,
	}
}

// AddSolutionToTicket associa uma solution do catálogo a um ticket
// POST /api/v1/tickets/:id/solutions
func (h *TicketSolutionHandler) AddSolutionToTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req dto.TicketSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.ticketService.AddSolutionToTicket(c.Request.Context(), ticketID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Solution added to ticket successfully"})
}

// GetTicketSolutions lista todas as solutions associadas a um ticket
// GET /api/v1/tickets/:id/solutions
func (h *TicketSolutionHandler) GetTicketSolutions(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	solutions, err := h.ticketService.GetTicketSolutions(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, solutions)
}

// RemoveSolutionFromTicket remove a associação entre uma solution e um ticket
// DELETE /api/v1/tickets/:id/solutions/:solution_id
func (h *TicketSolutionHandler) RemoveSolutionFromTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	solutionIDStr := c.Param("solution_id")
	solutionID, err := strconv.Atoi(solutionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid solution ID"})
		return
	}

	err = h.ticketService.RemoveSolutionFromTicket(c.Request.Context(), ticketID, solutionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Solution removed from ticket successfully"})
}
