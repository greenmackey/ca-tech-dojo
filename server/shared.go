package server

import (
	"net/http"
	"os"
	"strings"
)

const InvalidTokenMsg = "Token is invalid."
const InvalidBodyMsg = "Request body is invalid."
const InternalErrMsg = "Internal Server Error."
const EncodingErrMsg = "cannot encode response in %s"

// トークンの取得
func GetToken(r *http.Request) string {
	token := r.Header.Get("X-Token")
	return token
}

// CORSに対応するようにレスポンスヘッダーに書き込み
func CORSHeader(w http.ResponseWriter) {
	// w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token")
}

// CORSに対応するようにAccess-Control-Allow-Originに書き込み
// リストに載っているオリジンだけ許可
func CORSOrigin(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		allowedOrigins := strings.FieldsFunc(os.Getenv("ALLOWED_ORIGINS"), func(r rune) bool { return r == 44 || r == 32 })
		for _, a := range allowedOrigins {
			if origin == a {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				return
			}
		}
	}
}
