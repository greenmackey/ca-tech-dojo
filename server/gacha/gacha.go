package gacha

import (
	"ca-tech-dojo/controller/gacha"
	"ca-tech-dojo/controller/usercharacter"
	"ca-tech-dojo/log"
	Usercharacter "ca-tech-dojo/model/usercharacter"
	"ca-tech-dojo/server"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func DrawGacha(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := r.Header.Get("X-Token")

	// リクエストbodyの内容取得
	// ガチャ回数を受け取る
	reqBody := drawGachaRequest{Times: -1}
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&reqBody)
	if err != nil || reqBody.Times < 0 {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	// ガチャを作成
	gachaEntity, err := gacha.New()
	if err != nil {
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		log.Logger.Error(errors.Wrap(err, "gachaEntity.New failed"))
		return
	}
	// ガチャ結果を取得，保存
	characters := gachaEntity.Draw(uint(reqBody.Times))
	relationships := usercharacter.ToRelationship(token, characters)
	if err := Usercharacter.BulkCreate(relationships); err != nil {
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		log.Logger.Error(errors.Wrap(err, "usercharacter.BulkCreate failed"))
		return
	}

	// レスボンスbodyの作成
	// ガチャ結果を返す
	resp := newDrawGachaResponse(characters)
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}
}
