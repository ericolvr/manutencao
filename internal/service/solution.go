package service

import (
	"context"
	"fmt"
	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type SolutionService interface {
	Create(ctx context.Context, solution *domain.Solution) (*domain.Solution, error)
	GetByID(ctx context.Context, id int) (*domain.Solution, error)
	GetAll(ctx context.Context) ([]domain.Solution, error)
	GetByProblem(ctx context.Context, problemID int) ([]domain.Solution, error)
	Update(ctx context.Context, id int, solution *domain.Solution) (*domain.Solution, error)
	Delete(ctx context.Context, id int) error
}

type solutionService struct {
	solutionRepo repository.SolutionRepository
	problemRepo  repository.ProblemRepository
}

func NewSolutionService(solutionRepo repository.SolutionRepository, problemRepo repository.ProblemRepository) SolutionService {
	return &solutionService{
		solutionRepo: solutionRepo,
		problemRepo:  problemRepo,
	}
}

func (s *solutionService) Create(ctx context.Context, solution *domain.Solution) (*domain.Solution, error) {
	// Validações básicas
	if solution.Name == "" {
		return nil, fmt.Errorf("solution name is required")
	}

	if solution.UnitPrice < 0 {
		return nil, fmt.Errorf("unit price must be non-negative")
	}

	if solution.ProblemID <= 0 {
		return nil, fmt.Errorf("valid problem ID is required")
	}

	// Verificar se o problema existe
	_, err := s.problemRepo.FindByID(ctx, solution.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	// Criar a solução
	id, err := s.solutionRepo.Create(ctx, solution)
	if err != nil {
		return nil, fmt.Errorf("failed to create solution: %w", err)
	}

	// Buscar a solução criada para retornar com dados completos
	createdSolution, err := s.solutionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created solution: %w", err)
	}

	return createdSolution, nil
}

func (s *solutionService) GetByID(ctx context.Context, id int) (*domain.Solution, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid solution ID")
	}

	solution, err := s.solutionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get solution: %w", err)
	}

	return solution, nil
}

func (s *solutionService) GetAll(ctx context.Context) ([]domain.Solution, error) {
	solutions, err := s.solutionRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get solutions: %w", err)
	}

	return solutions, nil
}

func (s *solutionService) GetByProblem(ctx context.Context, problemID int) ([]domain.Solution, error) {
	if problemID <= 0 {
		return nil, fmt.Errorf("invalid problem ID")
	}

	// Verificar se o problema existe
	_, err := s.problemRepo.FindByID(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	solutions, err := s.solutionRepo.FindByProblem(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get solutions for problem: %w", err)
	}

	return solutions, nil
}

func (s *solutionService) Update(ctx context.Context, id int, solution *domain.Solution) (*domain.Solution, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid solution ID")
	}

	if solution.Name == "" {
		return nil, fmt.Errorf("solution name is required")
	}

	if solution.UnitPrice < 0 {
		return nil, fmt.Errorf("unit price must be non-negative")
	}

	if solution.ProblemID <= 0 {
		return nil, fmt.Errorf("valid problem ID is required")
	}

	// Verificar se a solução existe
	_, err := s.solutionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("solution not found: %w", err)
	}

	// Verificar se o problema existe
	_, err = s.problemRepo.FindByID(ctx, solution.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	// Atualizar a solução
	err = s.solutionRepo.Update(ctx, id, solution)
	if err != nil {
		return nil, fmt.Errorf("failed to update solution: %w", err)
	}

	// Buscar a solução atualizada
	updatedSolution, err := s.solutionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated solution: %w", err)
	}

	return updatedSolution, nil
}

func (s *solutionService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid solution ID")
	}

	// Verificar se a solução existe
	_, err := s.solutionRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("solution not found: %w", err)
	}

	// Deletar a solução
	err = s.solutionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete solution: %w", err)
	}

	return nil
}
