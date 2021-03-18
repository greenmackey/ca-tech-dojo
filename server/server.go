package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"ca-tech-dojo/model/user"
	"ca-tech-dojo/model/character"
)

const invalidTokenMsg = "Token is invalid."
const invalidBodyMsg = "Request body is invalid."
const invalidMethodMsg = "Method is not allowed."
const internalErrMsg = "Internal Server Error."
const secret = "***"

// トークンの取得
func getToken(r *http.Request) string {
	token := r.Header.Get("X-Token")
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
	if origin := r.Header.Get("Origin"); origin == "" {
		allowedOrigins := strings.FieldsFunc(os.Getenv("ALLOWED_ORIGINS"), func(r rune) bool { return r == 44 || r == 32 })
		for _, a := range allowedOrigins {
			if origin == a {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				return
			}
		}
	}
}


func CreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// CORS対応
		allowOrigins(w, r)

		// リクエストbodyの内容を取得
		// 新規ユーザの名前を受け取る
		var u user.User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&u)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		// トークンの生成
		token := strings.Replace(uuid.New().String(), "-", "", -1)

		// DBに追加
		err = user.Create(token, u.Name)
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// CORS対応
		allowOrigins(w, r)

		// トークン取得
		token := getToken(r)

		// DBからユーザ情報取得
		u, err := user.Get(token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// レスポンスbodyの作成
		// ユーザの名前を返す
		ec := json.NewEncoder(w)
		err = ec.Encode(u)
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// 該当するユーザの存在確認
		if err := user.VerifyToken(token); err != nil {
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// リクエストbodyの内容取得
		// 新しいユーザの名前を取得
		var u user.User
		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&u)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidBodyMsg, http.StatusBadRequest)
			return
		}

		// DB更新
		// ユーザの名前を更新
		if err := user.Update(token, u.Name); err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
		}

	case "OPTIONS":
		dealCORS(w, r)
	default:
		http.Error(w, invalidMethodMsg, http.StatusMethodNotAllowed)
		return
	}
}

func DrawGacha(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// 該当するユーザの存在確認
		if err := user.VerifyToken(token); err != nil {
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
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
		chars, err := user.DrawGacha(token, b.Times)
		if err != nil {
			log.Print(err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
		}

		// レスボンスbodyの作成
		// ガチャ結果を返す
		var resp struct {
			Characters []*character.Character `json:"results"`
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

func ListCharacters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// CORS対応
		allowOrigins(w, r)

		// トークンの取得
		token := getToken(r)

		// 該当するユーザの存在確認
		if err := user.VerifyToken(token); err != nil {
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// DBからユーザのガチャ結果を取得
		rels, err := user.RelCharacters(token)
		if err != nil {
			log.Print(err)
			http.Error(w, invalidTokenMsg, http.StatusBadRequest)
			return
		}

		// レスポンスbodyの作成
		// 該当ユーザのガチャ結果を返す
		var resp struct {
			RelsUserCharacter []user.RelUserCharacter `json:"characters"`
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
