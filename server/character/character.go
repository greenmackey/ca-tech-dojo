package character

import (
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/reluc"
	"ca-tech-dojo/model/user"
	"ca-tech-dojo/server"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func ListCharacters(w http.ResponseWriter, r *http.Request) {
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

	// DBからユーザのガチャ結果を取得
	rels, err := reluc.Get(token)
	if err != nil {
		log.Logger.Error(errors.Wrap(err, "reluc.Get failed"))
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// レスポンスbodyの作成
	// 該当ユーザのガチャ結果を返す
	var resp struct {
		RelsUserCharacter []reluc.Relationship `json:"characters"`
	}
	resp.RelsUserCharacter = rels
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}
}
