package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type TicketService interface {
	Create(ctx context.Context, req *dto.TicketRequest) (*dto.TicketResponse, error)
	List(ctx context.Context, limit, offset int) ([]dto.TicketResponse, int, error)
	ListWithDetails(ctx context.Context, limit, offset int) ([]dto.TicketResponse, int, error)
	FindByID(ctx context.Context, id int) (*dto.TicketResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error)
	Delete(ctx context.Context, id int) error
	GetTicketNumber(ctx context.Context) (int, error)
	AddProvider(ctx context.Context, ticketID int, req *dto.AddProviderRequest) error
	RemoveProvider(ctx context.Context, ticketID int) error
	GetProviderOnTicket(ctx context.Context, ticketID int) (*dto.ProviderSummaryResponse, error)

	// Ticket Problem methods
	AddProblemToTicket(ctx context.Context, ticketID int, req *dto.TicketProblemRequest) error
	GetTicketProblems(ctx context.Context, ticketID int) ([]dto.TicketProblemResponse, error)
	RemoveProblemFromTicket(ctx context.Context, ticketID int, problemID int) error

	// Ticket Solution methods
	AddSolutionToTicket(ctx context.Context, ticketID int, req *dto.TicketSolutionRequest) error
	GetTicketSolutions(ctx context.Context, ticketID int) ([]dto.TicketSolutionResponse, error)
	RemoveSolutionFromTicket(ctx context.Context, ticketID int, solutionID int) error

	// Ticket Distance methods
	AddKilometersValueToTicket(ctx context.Context, ticketID int, kilometers float64) error
	RemoveKilometersValueFromTicket(ctx context.Context, ticketID int) error
}

type ticketService struct {
	ticketRepo      repository.TicketRepository
	branchRepo      repository.BranchRepository
	providerRepo    repository.ProviderRepository
	problemRepo     repository.ProblemRepository
	solutionRepo    repository.SolutionRepository
	distanceService DistanceService
}

func NewTicketService(
	ticketRepo repository.TicketRepository,
	branchRepo repository.BranchRepository,
	providerRepo repository.ProviderRepository,
	problemRepo repository.ProblemRepository,
	solutionRepo repository.SolutionRepository,
	distanceService DistanceService,
) TicketService {
	return &ticketService{
		ticketRepo:      ticketRepo,
		branchRepo:      branchRepo,
		providerRepo:    providerRepo,
		problemRepo:     problemRepo,
		solutionRepo:    solutionRepo,
		distanceService: distanceService,
	}
}

func (s *ticketService) Create(ctx context.Context, req *dto.TicketRequest) (*dto.TicketResponse, error) {
	// Validar se branch existe
	_, err := s.branchRepo.FindByID(ctx, req.BranchID)
	if err != nil {
		return nil, fmt.Errorf("branch not found: %w", err)
	}

	// Parsear data de abertura
	openDate, err := time.Parse(time.RFC3339, req.OpenDate)
	if err != nil {
		return nil, fmt.Errorf("invalid open_date format: %w", err)
	}

	// Criar domínio do ticket (sem provider e sem custos iniciais)
	ticket := &domain.Ticket{
		Number:      req.Number,
		Status:      req.Status,
		Priority:    req.Priority,
		Description: req.Description,
		OpenDate:    openDate,
		BranchID:    req.BranchID,
		ProviderID:  nil, // Será associado posteriormente
	}

	// Criar ticket no repositório
	ticketID, err := s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Buscar ticket criado
	createdTicket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created ticket: %w", err)
	}

	// Retornar ticket sem custos (custos serão adicionados posteriormente)
	return dto.ToTicketResponse(createdTicket), nil
}

func (s *ticketService) List(ctx context.Context, limit, offset int) ([]dto.TicketResponse, int, error) {
	tickets, total, err := s.ticketRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets: %w", err)
	}

	var responses []dto.TicketResponse
	for _, ticket := range tickets {
		// Buscar informações do branch
		branch, err := s.branchRepo.FindByID(ctx, ticket.BranchID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to find branch for ticket %d: %w", ticket.ID, err)
		}

		// Buscar informações do provider (se existir)
		var provider *domain.Provider
		if ticket.ProviderID != nil {
			provider, err = s.providerRepo.FindByID(ctx, *ticket.ProviderID)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to find provider for ticket %d: %w", ticket.ID, err)
			}
		}

		// Buscar custos para cada ticket
		costs, err := s.ticketRepo.GetTicketCosts(ctx, ticket.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get ticket costs for ticket %d: %w", ticket.ID, err)
		}

		// Buscar distance do ticket pelo número
		var distanceValue *float64
		distance, err := s.distanceService.FindByNumber(ctx, ticket.Number)
		if err == nil && distance != nil {
			distanceValue = &distance.Distance
		}

		responses = append(responses, *dto.ToTicketResponseWithBranchProviderDistanceAndCosts(&ticket, branch, provider, distanceValue, costs))
	}

	return responses, total, nil
}

func (s *ticketService) ListWithDetails(ctx context.Context, limit, offset int) ([]dto.TicketResponse, int, error) {
	ticketsWithDetails, total, err := s.ticketRepo.ListWithDetails(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets with details: %w", err)
	}

	var responses []dto.TicketResponse
	for _, ticketDetail := range ticketsWithDetails {
		// Buscar custos para cada ticket (única query adicional necessária)
		costs, err := s.ticketRepo.GetTicketCosts(ctx, ticketDetail.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get ticket costs for ticket %d: %w", ticketDetail.ID, err)
		}

		// Converter para TicketResponse usando a nova função
		response := dto.ToTicketResponseFromDetails(&ticketDetail, costs)
		responses = append(responses, *response)
	}

	return responses, total, nil
}

func (s *ticketService) FindByID(ctx context.Context, id int) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Buscar informações do branch
	branch, err := s.branchRepo.FindByID(ctx, ticket.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to find branch: %w", err)
	}

	// Buscar informações do provider (se existir)
	var provider *domain.Provider
	if ticket.ProviderID != nil {
		provider, err = s.providerRepo.FindByID(ctx, *ticket.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("failed to find provider: %w", err)
		}
	}

	// Buscar custos do ticket
	costs, err := s.ticketRepo.GetTicketCosts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket costs: %w", err)
	}

	// Buscar distance do ticket pelo número
	var distanceValue *float64
	distance, err := s.distanceService.FindByNumber(ctx, ticket.Number)
	if err == nil && distance != nil {
		distanceValue = &distance.Distance
	}

	return dto.ToTicketResponseWithBranchProviderDistanceAndCosts(ticket, branch, provider, distanceValue, costs), nil
}

func (s *ticketService) Update(ctx context.Context, id int, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	// Verificar se ticket existe
	existingTicket, err := s.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Validar se branch existe
	_, err = s.branchRepo.FindByID(ctx, req.BranchID)
	if err != nil {
		return nil, fmt.Errorf("branch not found: %w", err)
	}

	// Validar se provider existe (se fornecido)
	if req.ProviderID != 0 {
		_, err := s.providerRepo.FindByID(ctx, req.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("provider not found: %w", err)
		}
	}

	// Parsear data de abertura
	openDate, err := time.Parse(time.RFC3339, req.OpenDate)
	if err != nil {
		return nil, fmt.Errorf("invalid open_date format: %w", err)
	}

	// Parsear data de fechamento (se fornecida)
	var closeDate *time.Time
	if req.CloseDate != nil {
		parsed, err := time.Parse(time.RFC3339, *req.CloseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid close_date format: %w", err)
		}
		closeDate = &parsed
	}

	// Atualizar campos do ticket
	existingTicket.Number = req.Number
	existingTicket.Status = int(req.Status)
	existingTicket.Priority = req.Priority
	existingTicket.Description = req.Description
	existingTicket.OpenDate = openDate
	existingTicket.CloseDate = closeDate
	existingTicket.BranchID = req.BranchID
	existingTicket.ProviderID = &req.ProviderID

	// Atualizar no repositório
	err = s.ticketRepo.Update(ctx, existingTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	// Processar solution_items e atualizar custos do ticket
	if len(req.SolutionItems) > 0 {
		var costs []domain.TicketCost
		for _, item := range req.SolutionItems {
			cost := domain.TicketCost{
				TicketID:   id,
				SolutionID: nil, // NULL - item customizado, não do catálogo de soluções
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				Subtotal:   float64(item.Quantity) * item.UnitPrice,
			}
			costs = append(costs, cost)
		}

		err = s.ticketRepo.UpdateTicketCosts(ctx, id, costs)
		if err != nil {
			return nil, fmt.Errorf("failed to update ticket costs: %w", err)
		}
	} else {
		// Se não há solution_items, remove todos os custos existentes
		err = s.ticketRepo.DeleteTicketCosts(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ticket costs: %w", err)
		}
	}

	// Buscar ticket atualizado
	updatedTicket, err := s.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated ticket: %w", err)
	}

	// Buscar custos do ticket atualizado
	costs, err := s.ticketRepo.GetTicketCosts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket costs: %w", err)
	}

	return dto.ToTicketResponseWithCosts(updatedTicket, costs), nil
}

func (s *ticketService) Delete(ctx context.Context, id int) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Remover custos do ticket antes de deletar o ticket
	err = s.ticketRepo.DeleteTicketCosts(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket costs: %w", err)
	}

	err = s.ticketRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	return nil
}

func (s *ticketService) GetTicketNumber(ctx context.Context) (int, error) {
	number, err := s.ticketRepo.GetTicketNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get ticket number: %w", err)
	}
	return number, nil
}

func (s *ticketService) AddProvider(ctx context.Context, ticketID int, req *dto.AddProviderRequest) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Verificar se provider existe
	_, err = s.providerRepo.FindByID(ctx, req.ProviderID)
	if err != nil {
		return fmt.Errorf("provider not found: %w", err)
	}

	err = s.ticketRepo.AddProvider(ctx, ticketID, req.ProviderID)
	if err != nil {
		return fmt.Errorf("failed to add provider to ticket: %w", err)
	}

	return nil
}

func (s *ticketService) RemoveProvider(ctx context.Context, ticketID int) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	err = s.ticketRepo.RemoveProvider(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to remove provider from ticket: %w", err)
	}

	return nil
}

func (s *ticketService) GetProviderOnTicket(ctx context.Context, ticketID int) (*dto.ProviderSummaryResponse, error) {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	provider, err := s.ticketRepo.GetProviderOnTicket(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	if provider == nil {
		return nil, fmt.Errorf("no provider associated with ticket")
	}

	return &dto.ProviderSummaryResponse{
		ID:      strconv.Itoa(provider.ID),
		Name:    provider.Name,
		Zipcode: provider.Zipcode,
	}, nil
}

// AddProblemToTicket associa um problema a um ticket
func (s *ticketService) AddProblemToTicket(ctx context.Context, ticketID int, req *dto.TicketProblemRequest) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Verificar se problema existe
	_, err = s.problemRepo.FindByID(ctx, req.ProblemID)
	if err != nil {
		return fmt.Errorf("problem not found: %w", err)
	}

	// Associar problema ao ticket
	err = s.ticketRepo.AddProblemToTicket(ctx, ticketID, req.ProblemID)
	if err != nil {
		return fmt.Errorf("failed to add problem to ticket: %w", err)
	}

	return nil
}

// GetTicketProblems retorna todos os problemas associados a um ticket
func (s *ticketService) GetTicketProblems(ctx context.Context, ticketID int) ([]dto.TicketProblemResponse, error) {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	ticketProblems, err := s.ticketRepo.GetTicketProblems(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket problems: %w", err)
	}

	// Converter para DTO com detalhes do problema
	var responses []dto.TicketProblemResponse
	for _, tp := range ticketProblems {
		// Buscar detalhes do problema
		problem, err := s.problemRepo.FindByID(ctx, tp.ProblemID)
		if err != nil {
			// Se não encontrar o problema, continua sem os detalhes
			responses = append(responses, dto.TicketProblemResponse{
				ID:        tp.ID,
				TicketID:  tp.TicketID,
				ProblemID: tp.ProblemID,
				CreatedAt: tp.CreatedAt,
			})
			continue
		}

		responses = append(responses, dto.TicketProblemResponse{
			ID:        tp.ID,
			TicketID:  tp.TicketID,
			ProblemID: tp.ProblemID,
			Problem: &dto.ProblemResponse{
				ID:          problem.ID,
				Name:        problem.Name,
				Description: problem.Description,
			},
			CreatedAt: tp.CreatedAt,
		})
	}

	return responses, nil
}

// RemoveProblemFromTicket remove a associação entre um problema e um ticket
func (s *ticketService) RemoveProblemFromTicket(ctx context.Context, ticketID int, problemID int) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Remover associação
	err = s.ticketRepo.RemoveProblemFromTicket(ctx, ticketID, problemID)
	if err != nil {
		return fmt.Errorf("failed to remove problem from ticket: %w", err)
	}

	return nil
}

// AddKilometersValueToTicket adiciona um custo de kilometragem a um ticket
func (s *ticketService) AddKilometersValueToTicket(ctx context.Context, ticketID int, kilometers float64) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Adicionar custo de kilometragem
	err = s.ticketRepo.AddKilometersValueToTicket(ctx, ticketID, kilometers)
	if err != nil {
		return fmt.Errorf("failed to add kilometers value to ticket: %w", err)
	}

	return nil
}

// RemoveKilometersValueFromTicket remove o custo de kilometragem de um ticket
func (s *ticketService) RemoveKilometersValueFromTicket(ctx context.Context, ticketID int) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Remover custo de kilometragem
	err = s.ticketRepo.RemoveKilometersValueFromTicket(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to remove kilometers value from ticket: %w", err)
	}

	return nil
}
