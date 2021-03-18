package main

import (
	"github.com/joho/godotenv"
  "log"
	"ca-tech-dojo/db"
	"ca-tech-dojo/server"
	"net/http"
)

// var db *sql.DB
// var err error

func main() {
	// .envファイルから環境変数を設定
	godotenv.Load()

	// // DBに接続
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// ログの設定
	initLog()

	// ルーティングとサーバの起動
	http.HandleFunc("/user/create", server.CreateUser)
	http.HandleFunc("/user/get", server.GetUser)
	http.HandleFunc("/user/update", server.UpdateUser)
	http.HandleFunc("/gacha/draw", server.DrawGacha)
	http.HandleFunc("/character/list", server.ListCharacters)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
