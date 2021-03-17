package main

import (
	"io"
	"log"
	"os"
)


func initLog() {
	// ログの出力先を標準出力とログファイルに設定
	f, err := os.OpenFile(os.Getenv("LOG_FILE_NAME"), os.O_APPEND|os.O_WRONLY, 0400)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(f, os.Stdout)
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("😮 ")
	log.SetOutput(multi)
}
