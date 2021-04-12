package usercharacter

type Relationship struct {
	Id             int
	UserToken      string
	CharacterId    int
	CharacterName  string
	CharacterPoint uint
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
