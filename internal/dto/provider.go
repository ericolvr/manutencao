package dto

type ProviderRequest struct {
	Name         string `json:"name" binding:"required"`
	Mobile       string `json:"mobile" binding:"required"`
	Zipcode      string `json:"zipcode" binding:"required"`
	State        string `json:"state" binding:"required"`
	City         string `json:"city" binding:"required"`
	Neighborhood string `json:"neighborhood" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Complement   string `json:"complement"`
}

type ProviderResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Mobile       string `json:"mobile"`
	Zipcode      string `json:"zipcode"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Address      string `json:"address"`
	Complement   string `json:"complement"`
}

type ProviderSummaryResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Zipcode string `json:"zipcode"`
}
