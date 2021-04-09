package character

import "ca-tech-dojo/model/usercharacter"

type listCharactersResponse struct {
	UserCharacters []listCharactersResponseRelationship `json:"characters"`
}

type listCharactersResponseRelationship struct {
	UserCharacterId int    `json:"userCharacterID,string"`
	CharacterId     int    `json:"characterID,string"`
	CharacterName   string `json:"name"`
}

func NewListCharactersResponse(relationships []usercharacter.Relationship) listCharactersResponse {
	resp := listCharactersResponse{
		UserCharacters: make([]listCharactersResponseRelationship, 0, len(relationships)),
	}
	for _, relationship := range relationships {
		userCharacterEntity := listCharactersResponseRelationship{
			UserCharacterId: relationship.Id,
			CharacterId:     relationship.CharacterId,
			CharacterName:   relationship.CharacterName,
		}
		resp.UserCharacters = append(resp.UserCharacters, userCharacterEntity)
	}
	return resp
}

type SellCharacterRequest struct {
	Id int `json:"characterID,string"`
}
