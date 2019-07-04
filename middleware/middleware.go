package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("victor-steven"))

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	//return a handler
	return func(res http.ResponseWriter, req *http.Request) {
		session, _ := Store.Get(req, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(res, req, "/login", 302)
			return
		}
		handler.ServeHTTP(res, req)
	}
}
