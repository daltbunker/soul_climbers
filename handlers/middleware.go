package handlers

import (
	"net/http"
	"os"
	"regexp"

	"github.com/daltbunker/soul_climbers/utils"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("MODE") == "dev" {
			NewDevSession(w, r)
		}

		session, err := GetSession(r)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			HandleUnautorized(w, r, false)
			return
		}

		match, err := regexp.MatchString("/admin",  r.URL.Path)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}

		if match && session.Values["role"].(string) != string(utils.Admin) {
			HandleUnautorized(w, r, true)
			return
		}

		next.ServeHTTP(w, r)
	})
}