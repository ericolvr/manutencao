package dto

import "github.com/ericolvr/maintenance-v2/internal/domain"

// SolutionRequest representa uma requisição para criar/atualizar solução
type SolutionRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	ProblemID   int     `json:"problem_id" binding:"required"`
}

// SolutionResponse representa uma solução na resposta
type SolutionResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	ProblemID   int     `json:"problem_id"`
}

// ToSolutionDomain converte DTO para domain
func (r *SolutionRequest) ToSolutionDomain() *domain.Solution {
	return &domain.Solution{
		Name:        r.Name,
		Description: r.Description,
		UnitPrice:   r.UnitPrice,
		ProblemID:   r.ProblemID,
	}
}

// ToSolutionResponse converte domain para DTO
func ToSolutionResponse(solution *domain.Solution) *SolutionResponse {
	if solution == nil {
		return nil
	}
	
	return &SolutionResponse{
		ID:          solution.ID,
		Name:        solution.Name,
		Description: solution.Description,
		UnitPrice:   solution.UnitPrice,
		ProblemID:   solution.ProblemID,
	}
}

// ToSolutionResponseList converte lista de domain para DTO
func ToSolutionResponseList(solutions []domain.Solution) []SolutionResponse {
	responses := make([]SolutionResponse, len(solutions))
	for i, solution := range solutions {
		responses[i] = *ToSolutionResponse(&solution)
	}
	return responses
}


