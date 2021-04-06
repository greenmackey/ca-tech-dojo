package gacha

import "ca-tech-dojo/model/character"

type DrawGachaRequest struct {
	Times int `json:"times"`
}

type DrawGachaResponse struct {
	Characters []*character.Character `json:"results"`
}
