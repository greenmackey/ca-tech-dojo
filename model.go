package main

type User struct {
	Name string `json:"name"`
}

type Character struct {
	Id         int    `json:"characterID,string"`
	Name       string `json:"name"`
	likelihood float64
}

type RelUserCharacter struct {
	Id            int    `json:"userCharacterID,string"`
	CharacterId   int    `json:"characterID,string"`
	CharacterName string `json:"name"`
}
