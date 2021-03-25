package server

import (
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/user"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func DrawGacha(w http.ResponseWriter, r *http.Request) {
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
	// ガチャ回数を受け取る
	var b struct{ Times uint }
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&b)
	if err != nil {
		http.Error(w, invalidBodyMsg, http.StatusBadRequest)
		return
	}

	// ガチャを引く
	chars, err := user.DrawGacha(token, b.Times)
	if err != nil {
		log.Print(errors.Wrapf(err, "cannot draw gacha in %s", "DrawGacha"))
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
		log.Print(errors.Wrapf(err, encodingErrMsg, "DrawGacha"))
		http.Error(w, internalErrMsg, http.StatusInternalServerError)
		return
	}
}
