package server

import (
	"ca-tech-dojo/model/user"
	"encoding/json"
	"log"
	"net/http"
)

func ListCharacters(w http.ResponseWriter, r *http.Request) {
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
}
