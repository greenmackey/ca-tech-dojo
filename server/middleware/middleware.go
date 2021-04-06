package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func CORSMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			// Access-Control-Allow-Originを設定
			if origin := req.Header.Get("Origin"); origin != "" {
				allowedOrigins := strings.FieldsFunc(os.Getenv("ALLOWED_ORIGINS"), func(r rune) bool { return r == 44 || r == 32 })
				for _, a := range allowedOrigins {
					if origin == a {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			// CORS preflight requestをさばく
			if req.Method == "OPTIONS" {
				// Access-Control-Allow-Headersを設定
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Token")
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
