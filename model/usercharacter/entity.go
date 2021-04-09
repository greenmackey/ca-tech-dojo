package usercharacter

type Relationship struct {
	Id            int    `json:"userCharacterID,string"`
	UserToken     string `json:"-"`
	CharacterId   int    `json:"characterID,string"`
	CharacterName string `json:"name"`
}

type queryErr struct {
	errMsg string
}

func newQueryErr(errMsg string) queryErr {
	return queryErr{errMsg: errMsg}
}

func (err queryErr) Error() string {
	return err.errMsg
}

func (queryErr) NotFound() bool {
	return true
}
