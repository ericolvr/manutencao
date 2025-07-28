package domain

type Provider struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Mobile       string `json:"mobile"`
	Zipcode      string `json:"zipcode"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Address      string `json:"address"`
	Complement   string `json:"complement"`
}
