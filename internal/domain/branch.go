package domain

type Branch struct {
	ID           int    `json:"id"`
	Client       string `json:"client"`
	Name         string `json:"name"`
	Uniorg       string `json:"uniorg"`
	Zipcode      string `json:"zipcode"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Address      string `json:"address"`
	Complement   string `json:"complement"`
}
