package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strings"
	"os"
)

type User struct {
	Id    int    `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name",omitempty`
}

var db *sql.DB
var err error

func getToken(r *http.Request) string {
	token := r.Header["X-Token"][0]
	return token
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		token := strings.Replace(uuid.New().String(), "-", "", -1)
		fmt.Println(user.Name, token)

		_, err = db.Query("INSERT INTO user (token, name) VALUES ( ?, ? )", token, user.Name)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := getToken(r)
		var user User
		err := db.QueryRow("select name from user where token = ?", token).Scan(&user.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Name)

		ec := json.NewEncoder(w)
		err = ec.Encode(user)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)

	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		token := getToken(r)

		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("update user set name=? where token=?", user.Name, token)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(user.Name, token)

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	godotenv.Load()
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/get", getUser)
	http.HandleFunc("/user/update", updateUser)
	http.ListenAndServe(":8080", nil)
}
