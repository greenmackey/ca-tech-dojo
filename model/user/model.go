package user

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/gacha"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const secret = "***"

func Create(token, name string) error {
	fq := "INSERT INTO users (token, name) VALUES ( %v, %v )"
	if _, err := db.DB.Exec(fmt.Sprintf(fq, "?", "?"), token, name); err != nil {
		return errors.Wrapf(err, "query failed in %s", "Create")
	}
	log.Logger.Info(fmt.Sprintf(fq, secret, name))
	return nil
}

func Get(token string) (User, error) {
	var user User
	fq := "SELECT name FROM users WHERE token = %v"
	if err := db.DB.QueryRow(strings.Replace(fq, "%v", "?", -1), token).Scan(&user.Name); err != nil {
		return user, errors.Wrapf(err, "query failed in %s", "Get")
	}
	log.Logger.Info(fmt.Sprintf(fq, secret))
	return user, nil
}

func Update(token, name string) error {
	fq := "UPDATE users SET name=%v WHERE token=%v"
	if _, err := db.DB.Exec(strings.Replace(fq, "%v", "?", -1), name, token); err != nil {
		return errors.Wrapf(err, "query failed in %s", "Update")
	}
	log.Logger.Info(fmt.Sprintf(fq, name, secret))
	return nil
}

func DrawGacha(token string, times uint) ([]*character.Character, error) {
	if times <= 0 {
		return []*character.Character{}, nil
	}

	var chars []*character.Character
	// ガチャを引く
	g, err := gacha.NewGacha()
	if err != nil {
		return chars, errors.Wrapf(err, "cannot get new gacha in %s", "DrawGacha")
	}

	chars = g.Draw(times)

	// DBにガチャの結果を反映
	partialfq := "INSERT INTO rel_user_character (user_token, character_id) VALUES "
	var placeholders []string
	var insert []interface{}
	for i := uint(0); i < times; i++ {
		placeholders = append(placeholders, "(%v, %v)")
		insert = append(insert, token, chars[i].Id)
	}
	fq := partialfq + strings.Join(placeholders, ", ")
	if _, err := db.DB.Exec(strings.Replace(fq, "%v", "?", -1), insert...); err != nil {
		return []*character.Character{}, errors.Wrapf(err, "query failed in %s", "DrawGacha")
	}
	log.Logger.Info(strings.Replace(fmt.Sprintf(fq, insert...), token, secret, -1))
	return chars, nil
}

func RelCharacters(token string) ([]RelUserCharacter, error) {
	fq := `SELECT r.id, chars.id, chars.name FROM rel_user_character AS r INNER JOIN characters AS chars ON r.character_id = chars.id AND r.user_token = %v`
	rows, err := db.DB.Query(strings.Replace(fq, "%v", "?", -1), token)
	if err != nil {
		return nil, errors.Wrapf(err, "query failed in %s", "RelCharacters")
	}
	log.Logger.Info(fmt.Sprintf(fq, secret))

	var rels []RelUserCharacter

	// ガチャ結果をスライス relsに格納
	for rows.Next() {
		var rel RelUserCharacter
		if err := rows.Scan(&rel.Id, &rel.CharacterId, &rel.CharacterName); err != nil {
			return nil, errors.Wrapf(err, "query failed in %s", "RelCharacters")
		}
		rels = append(rels, rel)
	}
	return rels, nil
}

func VerifyToken(token string) error {
	fq := "SELECT id FROM users WHERE token = %v"
	if err := db.DB.QueryRow(strings.Replace(fq, "%v", "?", -1), token).Scan(0); err == sql.ErrNoRows {
		return errors.Wrapf(err, "query failed in %s", "VerifyToken")
	}
	log.Logger.Info(fmt.Sprintf(fq, secret))
	return nil
}
