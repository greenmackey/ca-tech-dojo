package middleware

import (
	"ca-tech-dojo/model/user"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const invalidTokenMsg = "Token is invalid."

func CORSMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			// Access-Control-Allow-Methodsを設定
			allMethods, err := getAllMethodsForRoute(r, req)
			if err == nil {
				for _, v := range allMethods {
					if v == "OPTIONS" {
						w.Header().Set("Access-Control-Allow-Methods", strings.Join(allMethods, ","))
					}
				}
			}

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

func TokenMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("X-Token")

			// 該当するユーザの存在確認
			if err := user.Verify(token); err != nil {
				http.Error(w, invalidTokenMsg, http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

/* mux.CORSMethodMiddlewareがOPTIONS requestに対してヘッダーAccess-Control-Allow-Methodsを設定しないというバグがあったので，
ソースコードを参考に作成
*/
func getAllMethodsForRoute(r *mux.Router, req *http.Request) ([]string, error) {
	var allMethods []string

	var match mux.RouteMatch
	if r.Match(req, &match) || match.MatchErr == mux.ErrMethodMismatch {
		methods, err := match.Route.GetMethods()
		if err != nil {
			return nil, err
		}

		allMethods = append(allMethods, methods...)
	}

	return allMethods, nil
}
