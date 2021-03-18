package character

type Character struct {
	Id         int    `json:"characterID,string"`
	Name       string `json:"name"`
	Likelihood float64
}