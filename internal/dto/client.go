package dto

type ClientRequest struct {
	Name string `json:"name" binding:"required"`
}

type ClientResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
