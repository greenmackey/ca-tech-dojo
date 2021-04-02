package gacha

import (
	"ca-tech-dojo/controller/gacha"
	"ca-tech-dojo/controller/reluc"
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/character"
	Reluc "ca-tech-dojo/model/reluc"
	"ca-tech-dojo/model/user"
	"ca-tech-dojo/server"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func DrawGacha(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	server.CORSOrigin(w, r)

	// CORS preflight requestをさばく
	if r.Method == "OPTIONS" {
		server.CORSHeader(w)
		return
	}

	// トークンの取得
	token := server.GetToken(r)

	// 該当するユーザの存在確認
	if err := user.Verify(token); err != nil {
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// リクエストbodyの内容取得
	// ガチャ回数を受け取る
	b := struct{ Times int }{Times: -1}
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&b)
	if err != nil || b.Times < 0 {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	// ガチャを作成
	g, err := gacha.New()
	if err != nil {
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		log.Logger.Error(errors.Wrap(err, "gacha.New failed"))
		return
	}
	// ガチャ結果を取得，保存
	characters := g.Draw(uint(b.Times))
	rels := reluc.ToRelationship(token, characters)
	if err := Reluc.BulkCreate(rels); err != nil {
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		log.Logger.Error(errors.Wrap(err, "Reluc.BulkCreate failed"))
		return
	}

	// レスボンスbodyの作成
	// ガチャ結果を返す
	var resp struct {
		Characters []*character.Character `json:"results"`
	}
	resp.Characters = characters
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}
}
