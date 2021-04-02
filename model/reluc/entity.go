package reluc

type Relationship struct {
	Id            int    `json:"userCharacterID,string"`
	UserToken     string `json:"-"`
	CharacterId   int    `json:"characterID,string"`
	CharacterName string `json:"name"`
}
