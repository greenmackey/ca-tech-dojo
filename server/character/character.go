package character

import (
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/user"
	"ca-tech-dojo/model/usercharacter"
	"ca-tech-dojo/server"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func ListCharacters(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := server.GetToken(r)

	// 該当するユーザの存在確認
	if err := user.Verify(token); err != nil {
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// DBからユーザのガチャ結果を取得
	relationships, err := usercharacter.Get(token)
	if err != nil {
		log.Logger.Error(errors.Wrap(err, "usercharacter.Get failed"))
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// レスポンスbodyの作成
	// 該当ユーザのガチャ結果を返す
	resp := ListCharactersResponse{UserCharacters: relationships}
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}
}

func SellCharacter(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := server.GetToken(r)

	// 該当するユーザの存在確認
	if err := user.Verify(token); err != nil {
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// リクエストbodyの内容取得
	// userCharacterIdを受け取る
	reqBody := new(SellCharacterRequest)
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&reqBody)
	if err != nil {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	//	ユーザのポイント変更とuserCharacterの削除
	if err := usercharacter.Sell(token, reqBody.Id); err != nil {
		log.Logger.Error(errors.Wrap(err, "usercharacter.Sell failed"))
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
	}
}
