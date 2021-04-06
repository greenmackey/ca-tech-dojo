package gacha

import "ca-tech-dojo/model/character"

type drawGachaRequest struct {
	Times int `json:"times"`
}

type drawGachaResponse struct {
	Characters []drawGachaResponseCharacter `json:"results"`
}

type drawGachaResponseCharacter struct {
	Id   int    `json:"characterID,string"`
	Name string `json:"name"`
}

func newDrawGachaResponse(characters []*character.Character) drawGachaResponse {
	resp := drawGachaResponse{
		Characters: make([]drawGachaResponseCharacter, 0, len(characters)),
	}
	for _, c := range characters {
		respCharacter := drawGachaResponseCharacter{
			Id:   c.Id,
			Name: c.Name,
		}
		resp.Characters = append(resp.Characters, respCharacter)
	}
	return resp
}
