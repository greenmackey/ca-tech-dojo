package user

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
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
	q := "SELECT name, point FROM users WHERE token = ?"
	if err := db.DB.QueryRow(q, token).Scan(&user.Name, &user.Point); err != nil {
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

func Verify(token string) error {
	q := "SELECT id FROM users WHERE token = ?"
	if err := db.DB.QueryRow(q, token).Scan(0); err == sql.ErrNoRows {
		return errors.Wrap(err, "Select query failed")
	}
	log.Logger.Info("Verify token")
	return nil
}
