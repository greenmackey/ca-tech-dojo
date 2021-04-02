package user

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/gacha"
	"database/sql"
	"strings"

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

	var characters []*character.Character
	// ガチャを引く
	g, err := gacha.NewGacha()
	if err != nil {
		return characters, errors.Wrap(err, "gacha.NewGacha() failed")
	}

	characters = g.Draw(times)

	// DBにガチャの結果を反映
	partialq := "INSERT INTO rel_user_character (user_token, character_id) VALUES "
	var placeholders []string
	var insert []interface{}
	for i := uint(0); i < times; i++ {
		placeholders = append(placeholders, "(?, ?)")
		insert = append(insert, token, characters[i].Id)
	}
	q := partialq + strings.Join(placeholders, ", ")
	if _, err := db.DB.Exec(q, insert...); err != nil {
		return []*character.Character{}, errors.Wrap(err, "Insert query failed")
	}
	log.Logger.Info("Save gacha results")
	return characters, nil
}

func RelCharacters(token string) ([]RelUserCharacter, error) {
	q := "SELECT r.id, chars.id, chars.name FROM rel_user_character AS r INNER JOIN characters AS chars ON r.character_id = chars.id AND r.user_token = ?"
	rows, err := db.DB.Query(q, token)
	if err != nil {
		return nil, errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Get gacha results")

	var rels []RelUserCharacter

	// ガチャ結果をスライス relsに格納
	for rows.Next() {
		var rel RelUserCharacter
		if err := rows.Scan(&rel.Id, &rel.CharacterId, &rel.CharacterName); err != nil {
			return nil, errors.Wrap(err, "rows.Scan failed")
		}
		rels = append(rels, rel)
	}
	return rels, nil
}

func VerifyToken(token string) error {
	q := "SELECT id FROM users WHERE token = ?"
	if err := db.DB.QueryRow(q, token).Scan(0); err == sql.ErrNoRows {
		return errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Verify token")
	return nil
}
