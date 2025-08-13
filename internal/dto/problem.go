package dto

import "github.com/ericolvr/maintenance-v2/internal/domain"

// representa o request para criar/atualizar problema
type ProblemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// representa a resposta do problema
type ProblemResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// converte DTO para domain
func (r *ProblemRequest) ToProblemDomain() *domain.Problem {
	return &domain.Problem{
		Name:        r.Name,
		Description: r.Description,
	}
}

// converte domain para DTO
func ToProblemResponse(problem *domain.Problem) *ProblemResponse {
	if problem == nil {
		return nil
	}

	return &ProblemResponse{
		ID:          problem.ID,
		Name:        problem.Name,
		Description: problem.Description,
	}
}

// converte lista de domain para DTO
func ToProblemResponseList(problems []domain.Problem) []ProblemResponse {
	responses := make([]ProblemResponse, len(problems))
	for i, problem := range problems {
		responses[i] = *ToProblemResponse(&problem)
	}
	return responses
}
