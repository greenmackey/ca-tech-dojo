package main

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/log"
	"ca-tech-dojo/server"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	// .envファイルから環境変数を設定
	godotenv.Load()

	// ログの設定
	if err := log.InitZap(); err != nil {
		fmt.Println(errors.Wrap(err, "initZap failed"))
		os.Exit(1)
	}

	// // DBに接続
	if err := db.InitDB(); err != nil {
		log.Logger.Fatal(err)
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
	log.Logger.Fatal(http.ListenAndServe(":8080", r))
}
