package main

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"ca-tech-dojo/server/character"
	"ca-tech-dojo/server/gacha"
	"ca-tech-dojo/server/middleware"
	"ca-tech-dojo/server/user"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	// .envファイルから環境変数を設定
	if err := godotenv.Load(); err != nil {
		fmt.Println(errors.Wrap(err, "godotenv.Load failed"))
		os.Exit(1)
	}

	// ログの設定
	if err := log.InitZap(); err != nil {
		fmt.Println(errors.Wrap(err, "initZap failed"))
		os.Exit(1)
	}

	// DBに接続
	if err := db.InitDB(); err != nil {
		log.Logger.Fatal(err)
	}

	// ルーティングとサーバの起動
	r := mux.NewRouter()
	// CORS Access-Control-Allowed-Methodsを自動的に付加
	r.Use(mux.CORSMethodMiddleware(r))
	// CORS その他のヘッダーを付加 & preflight requestをさばく
	r.Use(middleware.CORSMiddleware(r))

	r.HandleFunc("/user/create", user.CreateUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/user/get", user.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/user/update", user.UpdateUser).Methods("PUT", "OPTIONS")
	r.HandleFunc("/gacha/draw", gacha.DrawGacha).Methods("POST", "OPTIONS")
	r.HandleFunc("/character/list", character.ListCharacters).Methods("GET", "OPTIONS")
	r.HandleFunc("/character/sell", character.SellCharacter).Methods("POST", "OPTIONS")
	log.Logger.Fatal(http.ListenAndServe(":8080", r))
}
