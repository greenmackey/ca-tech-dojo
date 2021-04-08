package character

import "ca-tech-dojo/model/usercharacter"

type ListCharactersResponse struct {
	UserCharacters []usercharacter.Relationship `json:"characters"`
}

type SellCharacterRequest struct {
	Id int `json:"userCharacterID,string"`
}
