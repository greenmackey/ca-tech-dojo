package main

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
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
