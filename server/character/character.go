package character

import (
	"ca-tech-dojo/log"
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/usercharacter"
	"ca-tech-dojo/server"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func ListCharacters(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := r.Header.Get("X-Token")

	// DBからユーザのガチャ結果を取得
	relationships, err := usercharacter.Get(token)
	if err != nil {
		log.Logger.Error(errors.Wrap(err, "usercharacter.Get failed"))
		http.Error(w, server.InvalidTokenMsg, http.StatusBadRequest)
		return
	}

	// レスポンスbodyの作成
	// 該当ユーザのガチャ結果を返す
	resp := NewListCharactersResponse(relationships)
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}
}

func SellCharacter(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := r.Header.Get("X-Token")

	// リクエストbodyの内容取得
	// CharacterIdを受け取る
	reqBody := new(SellCharacterRequest)
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&reqBody)
	if err != nil {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	//	ユーザのポイント変更とuserCharacterの削除
	if err := usercharacter.Sell(token, reqBody.Id); err != nil {
		if _, ok := errors.Cause(err).(interface{ NotFound() bool }); ok {
			http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		} else {
			log.Logger.Error(errors.Wrap(err, "usercharacter.Sell failed"))
			http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		}
	}
}

func BuyCharacter(w http.ResponseWriter, r *http.Request) {
	// トークンの取得
	token := r.Header.Get("X-Token")

	// リクエストbodyの内容取得
	// CharacterIdを受け取る
	reqBody := new(BuyCharacterRequest)
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&reqBody)
	if err != nil {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	// ポイントが十分にあるかチェック
	if ok, err := usercharacter.CheckBuyFeasibility(token, reqBody.Id); err != nil {
		log.Logger.Error(errors.Wrap(err, "usercharacter.CheckBuyFeasibility failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, server.InvalidBodyMsg, http.StatusBadRequest)
		return
	}

	//	ユーザのポイント変更とuserCharacterの追加
	if err := usercharacter.Buy(token, reqBody.Id); err != nil {
		log.Logger.Error(errors.Wrap(err, "usercharacter.Buy failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
	}
}

func GetAllCharacters(w http.ResponseWriter, r *http.Request) {

	characters, err := character.All()
	if err != nil {
		log.Logger.Error(errors.Wrap(err, "character.All failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}

	// レスポンスbodyの作成
	// 該当ユーザのガチャ結果を返す
	resp := NewGetAllCharactersResponse(characters)
	ec := json.NewEncoder(w)
	if err := ec.Encode(resp); err != nil {
		log.Logger.Error(errors.Wrap(err, "ec.Encode failed"))
		http.Error(w, server.InternalErrMsg, http.StatusInternalServerError)
		return
	}

}
