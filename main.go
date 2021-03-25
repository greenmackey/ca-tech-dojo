package main

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルから環境変数を設定
	godotenv.Load()

	// ログの設定
	if err := initLog(); err != nil {
		log.Fatal(err)
	}

	// // DBに接続
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	// ルーティングとサーバの起動
	r := mux.NewRouter()
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/create", server.CreateUser).Methods("POST", "OPTIONS")
	u.HandleFunc("/get", server.GetUser).Methods("GET", "OPTIONS")
	u.HandleFunc("/update", server.UpdateUser).Methods("PUT", "OPTIONS")
	r.HandleFunc("/gacha/draw", server.DrawGacha).Methods("POST", "OPTIONS")
	r.HandleFunc("/character/list", server.ListCharacters).Methods("GET", "OPTIONS")
	// CORS Access-Control-Allowed-Methodsを自動的に付加
	u.Use(mux.CORSMethodMiddleware(u))
	log.Fatal(http.ListenAndServe(":8080", r))
}
