package main

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

func initLog() error {
	// ログのプリフィックスを設定
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
	log.SetPrefix("😮 ")

	// ログの出力先を標準出力とログファイルに設定
	f, err := os.OpenFile(os.Getenv("LOG_FILE_NAME"), os.O_APPEND|os.O_WRONLY, 0400)
	if err != nil {
		return errors.Wrap(err, "cannot open")
	}
	multi := io.MultiWriter(f, os.Stdout)
	log.SetOutput(multi)
	return nil
}
