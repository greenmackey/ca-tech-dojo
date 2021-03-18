package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"fmt"
)

var DB *sql.DB

func InitDB() error {
	// DBに接続
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_NAME"))
	DB, err = sql.Open("mysql", dataSourceName)
	return err
}
