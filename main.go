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

// CORSに対応するようにレスポンスヘッダーに書き込み
func dealCORS(w http.ResponseWriter, r *http.Request) {
	allowOrigins(w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token")
}

// CORSに対応するようにAccess-Control-Allow-Originに書き込み
// リストに載っているオリジンだけ許可
func allowOrigins(w http.ResponseWriter, r *http.Request) {
	if origins, ok := r.Header["Origin"]; ok {
		origin := origins[0]
		allowedOrigins := strings.FieldsFunc(os.Getenv("ALLOWED_ORIGINS"), func(r rune) bool { return r == 44 || r == 32 })
		for _, a := range allowedOrigins {
			if origin == a {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				return
			}
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// CORS対応
		allowOrigins(w, r)

		// リクエストbodyの内容を取得
		// 新規ユーザの名前を受け取る
		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		// トークンの生成
		token := strings.Replace(uuid.New().String(), "-", "", -1)

		// DBに追加
		// ユーザのトークンと名前を追加
		_, err = db.Exec("INSERT INTO users (token, name) VALUES ( ?, ? )", token, user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}

		// レスポンスbodyの作成
		// 生成したトークンを返す
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
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// CORS対応
		allowOrigins(w, r)

		// トークン取得
		token := getToken(r)

		// DBからユーザ情報取得
		var user User
		err := db.QueryRow("select name from users where token = ?", token).Scan(&user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// レスポンスbodyの作成
		// ユーザの名前を返す
		ec := json.NewEncoder(w)
		err = ec.Encode(user)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}
	case "OPTIONS":
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// リクエストbodyの内容取得
		// 新しいユーザの名前を取得
		var user User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&user)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		// DB更新
		// ユーザの名前を更新
		// 該当するユーザが見つからなければエラー
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
	case "OPTIONS":
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func drawGacha(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// リクエストbodyの内容取得
		// ガチャ回数を受け取る
		var b struct{ Times int }
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&b)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		// ガチャを引く
		var chars = gacha.Draw(b.Times)

		// DBにガチャの結果を反映
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

		// レスボンスbodyの作成
		// ガチャ結果を返す
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
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func listCharacters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// 該当するユーザの存在確認
		q := "select id from users where token = ?"
		if err := db.QueryRow(q, token).Scan(0); err == sql.ErrNoRows {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// DBからユーザのガチャ結果を取得
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

		// ガチャ結果をスライス relsに格納
		for rows.Next() {
			var rel RelUserCharacter
			if err := rows.Scan(&rel.Id, &rel.CharacterId, &rel.CharacterName); err != nil {
				log.Print(err)
				http.Error(w, internalErrMsg, http.StatusInternalServerError)
				return
			}
			rels = append(rels, rel)
		}

		// レスポンスbodyの作成
		// 該当ユーザのガチャ結果を返す
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
	case "OPTIONS":
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

// Gachaを生成
// キャラクターのリストとその出現確率の累積値を管理するregionを格納
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

// ガチャを引く回数に対して，ガチャの結果を返す
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
	// .envファイルから環境変数を設定
	godotenv.Load()

	// DBに接続
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// ガチャを設定
	gacha, err = NewGacha()
	if err != nil {
		log.Fatal(err)
	}

	// ルーティングとサーバの起動
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/get", getUser)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/gacha/draw", drawGacha)
	http.HandleFunc("/character/list", listCharacters)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
