package dto

type ClientRequest struct {
	Name string `json:"name" binding:"required"`
}

type ClientResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
