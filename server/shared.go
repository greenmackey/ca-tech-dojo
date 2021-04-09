package server

import (
	"net/http"
)

const InvalidTokenMsg = "Token is invalid."
const InvalidBodyMsg = "Request body is invalid."
const InternalErrMsg = "Internal Server Error."

// トークンの取得
func GetToken(r *http.Request) string {
	token := r.Header.Get("X-Token")
	return token
}
