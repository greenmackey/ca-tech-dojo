package character

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"

	"github.com/pkg/errors"
)

func All() ([]*Character, error) {
	q := "SELECT id, name, likelihood, point FROM characters ORDER BY id ASC"
	rows, err := db.DB.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Get all characters")

	var characters []*Character

	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.Id, &c.Name, &c.Likelihood, &c.Point); err != nil {
			return nil, errors.Wrap(err, "rows.Scan failed")
		}
		characters = append(characters, &c)
	}
	return characters, nil
}
