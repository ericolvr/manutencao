package dto

type CostRequest struct {
	ValuePerKm   float64 `json:"value_per_km" binding:"required"`
	InitialValue float64 `json:"initial_value" binding:"required"`
}

type CostResponse struct {
	ID           int     `json:"id"`
	ValuePerKm   float64 `json:"value_per_km"`
	InitialValue float64 `json:"initial_value"`
}
