package character

import (
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/usercharacter"
)

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

type BuyCharacterRequest struct {
	Id int `json:"characterID,string"`
}

type GetAllCharactersResponse struct {
	Characters []GetAllCharactersResponseCharacter `json:"characters"`
}

type GetAllCharactersResponseCharacter struct {
	Id    int    `json:"characterId,string"`
	Name  string `json:"name"`
	Point uint   `json:"point"`
}

func NewGetAllCharactersResponse(characters []*character.Character) GetAllCharactersResponse {
	resp := GetAllCharactersResponse{
		Characters: make([]GetAllCharactersResponseCharacter, 0, len(characters)),
	}
	for _, c := range characters {
		characterEntity := GetAllCharactersResponseCharacter{
			Id:    c.Id,
			Name:  c.Name,
			Point: c.Point,
		}
		resp.Characters = append(resp.Characters, characterEntity)
	}
	return resp
}
