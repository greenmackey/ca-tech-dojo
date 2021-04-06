package usercharacter

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"strings"

	"github.com/pkg/errors"
)

func BulkCreate(rels []Relationship) error {
	if len(rels) == 0 {
		return nil
	}

	partialq := "INSERT INTO rel_user_character (user_token, character_id) VALUES "
	var placeholders []string
	var insert []interface{}
	for _, rel := range rels {
		placeholders = append(placeholders, "(?, ?)")
		insert = append(insert, rel.UserToken, rel.CharacterId)
	}

	q := partialq + strings.Join(placeholders, ", ")
	if _, err := db.DB.Exec(q, insert...); err != nil {
		return errors.Wrap(err, "Insert query failed")
	}

	log.Logger.Info("Save user-character relationships")
	return nil
}

func Get(token string) ([]Relationship, error) {
	var rels []Relationship

	q := "SELECT r.id, chars.id, chars.name FROM rel_user_character AS r INNER JOIN characters AS chars ON r.character_id = chars.id AND r.user_token = ?"
	rows, err := db.DB.Query(q, token)
	if err != nil {
		return nil, errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Get user-character relationships ")

	// ガチャ結果をスライス relsに格納
	for rows.Next() {
		var rel Relationship
		if err := rows.Scan(&rel.Id, &rel.CharacterId, &rel.CharacterName); err != nil {
			return nil, errors.Wrap(err, "rows.Scan failed")
		}
		rels = append(rels, rel)
	}
	return rels, nil
}
