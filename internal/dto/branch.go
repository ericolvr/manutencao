package dto

type BranchRequest struct {
	Name         string `json:"name" binding:"required"`
	Client       string `json:"client" binding:"required"`
	Uniorg       string `json:"uniorg" binding:"required"`
	Zipcode      string `json:"zipcode" binding:"required"`
	State        string `json:"state" binding:"required"`
	City         string `json:"city" binding:"required"`
	Neighborhood string `json:"neighborhood" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Complement   string `json:"complement"`
}

type BranchResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Client       string `json:"client"`
	Uniorg       string `json:"uniorg"`
	Zipcode      string `json:"zipcode"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Address      string `json:"address"`
	Complement   string `json:"complement"`
}

type BranchSummaryResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Uniorg  string `json:"uniorg"`
	Zipcode string `json:"zipcode"`
}
