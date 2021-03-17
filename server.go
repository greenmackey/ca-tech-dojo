package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
)

const invalidTokenMsg = "Token is invalid."
const invalidBodyMsg = "Request body is invalid."
const invalidMethodMsg = "Method is not allowed."
const internalErrMsg = "Internal Server Error."
const secret = "***"

// トークンの取得
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

// トークンに対応するユーザがいるかどうか確認
// いればtrueを返す
// いなければエラーメッセージをレスポンスボディに書き込み，falseを返す
func verifyToken(w http.ResponseWriter, token string) bool {
	fq := "SELECT id FROM users WHERE token = %v"
	if err := db.QueryRow(strings.Replace(fq, "%v", "?", -1), token).Scan(0); err == sql.ErrNoRows {
		log.Print(err)
		http.Error(w, invalidTokenMsg, http.StatusBadRequest)
		return false
	} else {
		log.Print(fmt.Sprintf(fq, secret))
	}
	return true
}

// ルーティングとサーバの起動
func initServer() {
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/user/get", getUser)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/gacha/draw", drawGacha)
	http.HandleFunc("/character/list", listCharacters)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
		fq := "INSERT INTO users (token, name) VALUES ( %v, %v )"
		_, err = db.Exec(fmt.Sprintf(fq, "?", "?"), token, user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		} else {
			log.Print(fmt.Sprintf(fq, secret, user.Name))
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
		fq := "SELECT name FROM users WHERE token = %v"
		err := db.QueryRow(strings.Replace(fq, "%v", "?", -1), token).Scan(&user.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		} else {
			log.Print(fmt.Sprintf(fq, secret))
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

		// 該当するユーザの存在確認
		if !verifyToken(w, token) {
			return
		}

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
		fq := "UPDATE users SET name=%v WHERE token=%v"
		_, err = db.Exec(strings.Replace(fq, "%v", "?", -1), user.Name, token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
		} else {
			log.Print(fmt.Sprintf(fq, user.Name, secret))
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

		// 該当するユーザの存在確認
		if !verifyToken(w, token) {
			return
		}

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
		gacha, err := NewGacha()
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
		}
		var chars = gacha.Draw(b.Times)

		// DBにガチャの結果を反映
		partialfq := "INSERT INTO rel_user_character (user_token, character_id) VALUES "
		var placeholders []string
		var insert []interface{}
		for i := 0; i < b.Times; i++ {
			placeholders = append(placeholders, "(%v, %v)")
			insert = append(insert, token, chars[i].Id)
		}
		fq := partialfq + strings.Join(placeholders, ", ")
		_, err = db.Exec(strings.Replace(fq, "%v", "?", -1), insert...)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		} else {
			log.Print(strings.Replace(fmt.Sprintf(fq, insert...), token, secret, -1))
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
		if !verifyToken(w, token) {
			return
		}

		// DBからユーザのガチャ結果を取得
		var rels []RelUserCharacter
		fq := `SELECT r.id, chars.id, chars.name FROM rel_user_character as r
		INNER JOIN characters AS chars
		ON r.character_id = chars.id AND r.user_token = %v`
		rows, err := db.Query(strings.Replace(fq, "%v", "?", -1), token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		} else {
			log.Print(fmt.Sprintf(fq, secret))
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
