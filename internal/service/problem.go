package service

import (
	"context"
	"fmt"
	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type ProblemService interface {
	Create(ctx context.Context, problem *domain.Problem) (*domain.Problem, error)
	GetByID(ctx context.Context, id int) (*domain.Problem, error)
	GetAll(ctx context.Context) ([]domain.Problem, error)
	Update(ctx context.Context, id int, problem *domain.Problem) (*domain.Problem, error)
	Delete(ctx context.Context, id int) error
}

type problemService struct {
	problemRepo repository.ProblemRepository
}

func NewProblemService(problemRepo repository.ProblemRepository) ProblemService {
	return &problemService{
		problemRepo: problemRepo,
	}
}

func (s *problemService) Create(ctx context.Context, problem *domain.Problem) (*domain.Problem, error) {
	// Validações básicas
	if problem.Name == "" {
		return nil, fmt.Errorf("problem name is required")
	}

	// Criar o problema
	id, err := s.problemRepo.Create(ctx, problem)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem: %w", err)
	}

	// Buscar o problema criado para retornar com dados completos
	createdProblem, err := s.problemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created problem: %w", err)
	}

	return createdProblem, nil
}

func (s *problemService) GetByID(ctx context.Context, id int) (*domain.Problem, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid problem ID")
	}

	problem, err := s.problemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (s *problemService) GetAll(ctx context.Context) ([]domain.Problem, error) {
	problems, err := s.problemRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get problems: %w", err)
	}

	return problems, nil
}

func (s *problemService) Update(ctx context.Context, id int, problem *domain.Problem) (*domain.Problem, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid problem ID")
	}

	if problem.Name == "" {
		return nil, fmt.Errorf("problem name is required")
	}

	// Verificar se o problema existe
	_, err := s.problemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	// Atualizar o problema
	err = s.problemRepo.Update(ctx, id, problem)
	if err != nil {
		return nil, fmt.Errorf("failed to update problem: %w", err)
	}

	// Buscar o problema atualizado
	updatedProblem, err := s.problemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated problem: %w", err)
	}

	return updatedProblem, nil
}

func (s *problemService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid problem ID")
	}

	// Verificar se o problema existe
	_, err := s.problemRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("problem not found: %w", err)
	}

	// Deletar o problema
	err = s.problemRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete problem: %w", err)
	}

	return nil
}
