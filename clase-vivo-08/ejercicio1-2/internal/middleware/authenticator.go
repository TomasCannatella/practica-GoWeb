package middleware

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type Authenticathor struct {
	token string
}

func NewAuthenticator(token string) *Authenticathor {
	return &Authenticathor{
		token: token,
	}
}

func (a *Authenticathor) Auth(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("token")

		if token != a.token {
			response.Error(w, http.StatusUnauthorized, "invalid Token")
			return
		}

		hd.ServeHTTP(w, r)
	})
}
