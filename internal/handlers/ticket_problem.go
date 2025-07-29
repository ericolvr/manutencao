package handlers

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type TicketProblemHandler struct {
	ticketService service.TicketService
}

func NewTicketProblemHandler(ticketService service.TicketService) *TicketProblemHandler {
	return &TicketProblemHandler{
		ticketService: ticketService,
	}
}

func (h *TicketProblemHandler) AddProblemToTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req dto.TicketProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.ticketService.AddProblemToTicket(c.Request.Context(), ticketID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Problem added to ticket successfully"})
}

func (h *TicketProblemHandler) GetTicketProblems(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	problems, err := h.ticketService.GetTicketProblems(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problems)
}

func (h *TicketProblemHandler) RemoveProblemFromTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	problemIDStr := c.Param("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	err = h.ticketService.RemoveProblemFromTicket(c.Request.Context(), ticketID, problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Problem removed from ticket successfully"})
}
