package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// CreateTicket cria um novo ticket
// @Summary Criar ticket
// @Description Cria um novo ticket de manutenção
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body dto.TicketRequest true "Dados do ticket"
// @Success 201 {object} dto.TicketResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req dto.TicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.ticketService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// ListTickets lista todos os tickets com paginação
// @Summary Listar tickets
// @Description Lista todos os tickets com paginação
// @Tags tickets
// @Produce json
// @Param limit query int false "Limite de registros por página" default(10)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets [get]
func (h *TicketHandler) ListTickets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tickets, total, err := h.ticketService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tickets,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// FindTicketByID busca um ticket por ID
// @Summary Buscar ticket por ID
// @Description Busca um ticket específico pelo seu ID
// @Tags tickets
// @Produce json
// @Param id path int true "ID do ticket"
// @Success 200 {object} dto.TicketResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id} [get]
func (h *TicketHandler) FindTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := h.ticketService.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateTicket atualiza um ticket existente
// @Summary Atualizar ticket
// @Description Atualiza um ticket existente
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "ID do ticket"
// @Param ticket body dto.UpdateTicketRequest true "Dados atualizados do ticket"
// @Success 200 {object} dto.TicketResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id} [put]
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.ticketService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket remove um ticket
// @Summary Deletar ticket
// @Description Remove um ticket do sistema
// @Tags tickets
// @Produce json
// @Param id path int true "ID do ticket"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id} [delete]
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	err = h.ticketService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTicketNumber obtém o próximo número de ticket disponível
// @Summary Obter próximo número de ticket
// @Description Retorna o próximo número sequencial disponível para um novo ticket
// @Tags tickets
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/number [get]
func (h *TicketHandler) GetTicketNumber(c *gin.Context) {
	number, err := h.ticketService.GetTicketNumber(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"number": number})
}

// AddProviderToTicket associa um prestador ao ticket
// @Summary Associar prestador ao ticket
// @Description Adiciona um prestador a um ticket específico
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "ID do ticket"
// @Param provider body dto.AddProviderRequest true "Dados do prestador"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id}/providers [post]
func (h *TicketHandler) AddProviderToTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req dto.AddProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.ticketService.AddProvider(c.Request.Context(), ticketID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider added successfully"})
}

// RemoveProviderFromTicket remove o prestador do ticket
// @Summary Remover prestador do ticket
// @Description Remove a associação de prestador de um ticket
// @Tags tickets
// @Produce json
// @Param id path int true "ID do ticket"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id}/providers [delete]
func (h *TicketHandler) RemoveProviderFromTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	err = h.ticketService.RemoveProvider(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider removed successfully"})
}

// GetProviderOnTicket obtém o prestador associado ao ticket
// @Summary Obter prestador do ticket
// @Description Retorna o prestador associado a um ticket específico
// @Tags tickets
// @Produce json
// @Param id path int true "ID do ticket"
// @Success 200 {object} dto.ProviderSummaryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id}/providers [get]
func (h *TicketHandler) GetProviderOnTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	provider, err := h.ticketService.GetProviderOnTicket(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}
