package log

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitZap() error {
	// jsonファイルの読み込み
	cfgJson, err := os.ReadFile("zap_config.json")
	if err != nil {
		return errors.Wrap(err, "cannot read")
	}

	// configを作成
	var cfg zap.Config
	if err := json.Unmarshal(cfgJson, &cfg); err != nil {
		return errors.Wrap(err, "cannot unmarshal")
	}
	cfg.OutputPaths = append(cfg.OutputPaths, os.Getenv("LOG_FILE_NAME"))
	cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, os.Getenv("LOG_FILE_NAME"))

	// loggerを生成
	orgLogger, err := cfg.Build()
	if err != nil {
		return errors.Wrap(err, "cannot build")
	}
	Logger = orgLogger.Sugar()
	return nil
}
