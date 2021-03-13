package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type User struct {
	// Id    int    `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Name  string `json:"name",omitempty`
}

type Character struct {
	Id         int     `json:"characterID"`
	Name       string  `json:"name"`
	likelihood float64 `json:",omitempty"`
}

type Gacha struct {
	characters []*Character
	region     []float64
}

var db *sql.DB
var err error

// var allCharacters Characters
var gacha Gacha

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

		_, err = db.Query("INSERT INTO users (token, name) VALUES ( ?, ? )", token, user.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := getToken(r)
		var user User
		err := db.QueryRow("select name from users where token = ?", token).Scan(&user.Name)
		if err != nil {
			log.Fatal(err)
		}

		ec := json.NewEncoder(w)
		err = ec.Encode(user)
		if err != nil {
			log.Fatal(err)
		}

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

		_, err = db.Exec("update users set name=? where token=?", user.Name, token)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func draw(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// token := getToken(r)

		var b struct{ Times int }
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&b)
		if err != nil {
			log.Fatal(err)
		}

		var chars = gacha.Draw(b.Times)
		var resp struct{
			Results []*Character `json:"results"`
		}
		resp.Results = chars
		ec := json.NewEncoder(w)
		err = ec.Encode(resp)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func setGacha() {
	rows, err := db.Query("select id, name, likelihood from characters order by id asc")
	if err != nil {
		log.Fatal(rows)
	}

	var total float64

	for rows.Next() {
		var c Character
		err := rows.Scan(&c.Id, &c.Name, &c.likelihood)
		if err != nil {
			log.Fatal(err)
		}
		gacha.characters = append(gacha.characters, &c)
		total += c.likelihood
	}

	gacha.region = []float64{0}
	var sum float64
	for _, c := range gacha.characters {
		c.likelihood /= total
		sum += c.likelihood
		gacha.region = append(gacha.region, sum)
	}
}

func (g Gacha) Draw(n int) []*Character {
	var chars []*Character
	rand.Seed(time.Now().UnixNano())
	for ; n > 0; n-- {
		p := rand.Float64()
		for i := 0; i < len(g.characters); i++ {
			if g.region[i] <= p && p < g.region[i+1] {
				chars = append(chars, g.characters[i])
				break
			}
		}
	}
	return chars
}

func main() {
	godotenv.Load()
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	setGacha()
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/get", getUser)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/gacha/draw", draw)
	http.ListenAndServe(":8080", nil)
}
