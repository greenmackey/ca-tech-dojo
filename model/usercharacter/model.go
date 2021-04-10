package usercharacter

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"context"
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

func Sell(token string, characterId int) error {
	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	userCharacterDeleteQuery := "DELETE FROM rel_user_character WHERE user_token = ? AND character_id = ? LIMIT 1"
	if result, err := tx.Exec(userCharacterDeleteQuery, token, characterId); err != nil {
		err = errors.Wrap(err, "userCharacterDeleteQuery failed")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "Rollback failed: %v", rollbackErr)
		}
		return err
	} else if n, _ := result.RowsAffected(); n == 0 {
		err = errors.Wrap(newQueryErr("The user has no character of that type"), "userCharacterDeleteQuery failed")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "Rollback failed: %v", rollbackErr)
		}
		return err
	}

	userPointQuery := "UPDATE users SET point = point + (SELECT point FROM characters WHERE id = ?) WHERE token = ?"
	if _, err := tx.Exec(userPointQuery, characterId, token); err != nil {
		err = errors.Wrap(err, "userPointQuery failed")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "Rollback failed: %v", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "tx.Commit failed")
	}

	log.Logger.Info("Delete user-character relationships and update point of user")
	return nil
}

func CheckBuyFeasibility(token string, characterId int) (bool, error) {
	var remainingPoint int
	pointQuery := "SELECT CAST(point AS SIGNED) - CAST((SELECT point FROM characters WHERE id = ?) AS SIGNED) FROM users WHERE token = ?"
	if err := db.DB.QueryRow(pointQuery, characterId, token).Scan(&remainingPoint); err != nil {
		return false, errors.Wrap(err, "pointQuery failed")
	}
	log.Logger.Info("Check feasibility of buying a character")

	if remainingPoint >= 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func Buy(token string, characterId int) error {
	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	userCharacterAddQuery := "INSERT INTO rel_user_character (user_token, character_id) VALUES (?, ?)"
	if _, err := tx.Exec(userCharacterAddQuery, token, characterId); err != nil {
		err = errors.Wrap(err, "userCharacterAddQuery failed")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "Rollback failed: %v", rollbackErr)
		}
		return err
	}

	userPointQuery := "UPDATE users SET point = point - (SELECT point FROM characters WHERE id = ?) WHERE token = ?"
	if _, err := tx.Exec(userPointQuery, characterId, token); err != nil {
		err = errors.Wrap(err, "userPointQuery failed")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "Rollback failed: %v", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "tx.Commit failed")
	}

	log.Logger.Info("Add user-character relationship and reduce point of user")

	return nil
}
