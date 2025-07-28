package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/handlers"
)

func TicketRoutes(router *gin.Engine, ticketHandler *handlers.TicketHandler, ticketProblemHandler *handlers.TicketProblemHandler, ticketSolutionHandler *handlers.TicketSolutionHandler) {
	tickets := router.Group("/api/v1/tickets")
	{
		// CRUD básico
		tickets.POST("", ticketHandler.CreateTicket)
		tickets.GET("", ticketHandler.ListTickets)
		tickets.GET("/:id", ticketHandler.FindTicketByID)
		tickets.PUT("/:id", ticketHandler.UpdateTicket)
		tickets.DELETE("/:id", ticketHandler.DeleteTicket)

		// Utilitários
		tickets.GET("/number", ticketHandler.GetTicketNumber)

		// Associações com prestadores
		tickets.POST("/:id/providers", ticketHandler.AddProviderToTicket)
		tickets.GET("/:id/providers", ticketHandler.GetProviderOnTicket)
		tickets.DELETE("/:id/providers", ticketHandler.RemoveProviderFromTicket)

		// Associações com problemas
		tickets.POST("/:id/problems", ticketProblemHandler.AddProblemToTicket)
		tickets.GET("/:id/problems", ticketProblemHandler.GetTicketProblems)
		tickets.DELETE("/:id/problems/:problem_id", ticketProblemHandler.RemoveProblemFromTicket)

		// Associações com solutions (custos)
		tickets.POST("/:id/solutions", ticketSolutionHandler.AddSolutionToTicket)
		tickets.GET("/:id/solutions", ticketSolutionHandler.GetTicketSolutions)
		tickets.DELETE("/:id/solutions/:solution_id", ticketSolutionHandler.RemoveSolutionFromTicket)
	}
}
