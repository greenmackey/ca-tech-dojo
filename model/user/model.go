package user

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/gacha"
	"ca-tech-dojo/model/reluc"
	"database/sql"

	"github.com/pkg/errors"
)

func Create(token, name string) error {
	q := "INSERT INTO users (token, name) VALUES ( ?, ? )"
	if _, err := db.DB.Exec(q, token, name); err != nil {
		return errors.Wrap(err, "Insert query failed")
	}
	log.Logger.Info("Create a user")
	return nil
}

func Get(token string) (User, error) {
	var user User
	q := "SELECT name FROM users WHERE token = ?"
	if err := db.DB.QueryRow(q, token).Scan(&user.Name); err != nil {
		return user, errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Get a user")
	return user, nil
}

func Update(token, name string) error {
	q := "UPDATE users SET name=? WHERE token=?"
	if _, err := db.DB.Exec(q, name, token); err != nil {
		return errors.Wrap(err, "Update query failed")
	}
	log.Logger.Info("Update a user")
	return nil
}

func DrawGacha(token string, times uint) ([]*character.Character, error) {
	if times <= 0 {
		return []*character.Character{}, nil
	}

	// ガチャを引く
	g, err := gacha.NewGacha()
	if err != nil {
		return nil, errors.Wrap(err, "gacha.NewGacha() failed")
	}

	characters := g.Draw(times)

	rels := make([]reluc.Relationship, 0, len(characters))

	for _, c := range characters {
		rel := reluc.Relationship{
			UserToken:   token,
			CharacterId: c.Id,
		}
		rels = append(rels, rel)
	}

	if err := reluc.BulkCreate(rels); err != nil {
		return nil, errors.Wrap(err, "reluc.BulkCreate failed")
	}

	return characters, nil
}

func Verify(token string) error {
	q := "SELECT id FROM users WHERE token = ?"
	if err := db.DB.QueryRow(q, token).Scan(0); err == sql.ErrNoRows {
		return errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Verify token")
	return nil
}
