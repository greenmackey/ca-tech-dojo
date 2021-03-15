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
	Name string `json:"name"`
}

type Character struct {
	Id         int    `json:"characterID,string"`
	Name       string `json:"name"`
	likelihood float64
}

type RelUserCharacter struct {
	Id            int    `json:"userCharacterID,string"`
	CharacterId   int    `json:"characterID,string"`
	CharacterName string `json:"name"`
}

type Gacha struct {
	characters []*Character
	region     []float64
}

var db *sql.DB
var err error
var gacha Gacha
const invalidTokenMsg = "Token is invalid."
const invalidBodyMsg = "Request body is invalid."
const invalidMethodMsg = "Method is not allowed."
const internalErrMsg = "Internal Server Error."

func getToken(r *http.Request) string {
	token := r.Header["X-Token"][0]
	return token
}

func dealCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}
		token := strings.Replace(uuid.New().String(), "-", "", -1)

		_, err = db.Exec("INSERT INTO users (token, name) VALUES ( ?, ? )", token, user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
		var resp struct {
			Token string `json:"token"`
		}
		resp.Token = token
		ec := json.NewEncoder(w)
		if err := ec.Encode(resp); err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}
	case "OPTIONS":
		dealCORS(w)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token := getToken(r)
		var user User
		err := db.QueryRow("select name from users where token = ?", token).Scan(&user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
		ec := json.NewEncoder(w)
		err = ec.Encode(user)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}
	case "OPTIONS":
		dealCORS(w)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		token := getToken(r)

		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("update users set name=? where token=?", user.Name, token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		} else if affected, err := result.RowsAffected(); err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		} else if affected == 0 {
			log.Print("sql: no rows in result set")
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
	case "OPTIONS":
		dealCORS(w)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func drawGacha(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		token := getToken(r)

		var b struct{ Times int }
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&b)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		var chars = gacha.Draw(b.Times)

		partialq := "insert into rel_user_character (user_token, character_id) values "
		var placeholders []string
		var insert []interface{}
		for i := 0; i < b.Times; i++ {
			placeholders = append(placeholders, "(?, ?)")
			insert = append(insert, token, chars[i].Id)
		}
		q := partialq + strings.Join(placeholders, ", ")
		_, err = db.Exec(q, insert...)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
		var resp struct {
			Characters []*Character `json:"results"`
		}
		resp.Characters = chars
		ec := json.NewEncoder(w)
		if err := ec.Encode(resp); err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}
	case "OPTIONS":
		dealCORS(w)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func listCharacters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token := getToken(r)

		q := "select id from users where token = ?"

		if err := db.QueryRow(q, token).Scan(0); err == sql.ErrNoRows {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		var rels []RelUserCharacter
		q = `select r.id, chars.id, chars.name from rel_user_character as r
		inner join characters as chars
		on r.character_id = chars.id and r.user_token = ?`
		rows, err := db.Query(q, token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		for rows.Next() {
			var rel RelUserCharacter
			if err := rows.Scan(&rel.Id, &rel.CharacterId, &rel.CharacterName); err != nil {
				log.Print(err)
				http.Error(w, internalErrMsg, http.StatusInternalServerError)
				return
			}
			rels = append(rels, rel)
		}

		var resp struct {
			RelsUserCharacter []RelUserCharacter `json:"characters"`
		}
		resp.RelsUserCharacter = rels
		ec := json.NewEncoder(w)
		if err := ec.Encode(resp); err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
	case "OPTIONS":
		dealCORS(w)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func NewGacha() (Gacha, error) {
	var gacha Gacha

	rows, err := db.Query("select id, name, likelihood from characters order by id asc")
	if err != nil {
		return Gacha{}, err
	}

	var total float64

	for rows.Next() {
		var c Character
		err := rows.Scan(&c.Id, &c.Name, &c.likelihood)
		if err != nil {
			return Gacha{}, err
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
	return gacha, nil
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
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	gacha, err = NewGacha()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/get", getUser)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/gacha/draw", drawGacha)
	http.HandleFunc("/character/list", listCharacters)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
