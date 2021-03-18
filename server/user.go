package server

import (
	"github.com/google/uuid"
	"net/http"
	"encoding/json"
	"log"
	"ca-tech-dojo/model/user"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	CORSOrigin(w, r)

	// CORS preflight requestをさばく
	if r.Method == "OPTIONS" {
		CORSHeader(w)
		return
	}

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
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	CORSOrigin(w, r)

	// CORS preflight requestをさばく
	if r.Method == "OPTIONS" {
		CORSHeader(w)
		return
	}

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
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	CORSOrigin(w, r)

	// CORS preflight requestをさばく
	if r.Method == "OPTIONS" {
		CORSHeader(w)
		return
	}

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
}
