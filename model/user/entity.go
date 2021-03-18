package user

type User struct {
	Name string `json:"name"`
}

type RelUserCharacter struct {
	Id            int    `json:"userCharacterID,string"`
	CharacterId   int    `json:"characterID,string"`
	CharacterName string `json:"name"`
}