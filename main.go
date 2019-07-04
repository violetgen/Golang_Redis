package main

import (
	"net/http"

	"./models"
	"./routes"
	"./utils"
	"github.com/go-redis/redis"
)

var client *redis.Client

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)

	//without mux
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":7000", nil)
}

// Instead of checking if a user is authenticated in the "indexGetHandler" alone, let us create a middleware thart will help us handle that:

//note: "AuthRequired" takes in a handler, for instance: "indexGetHandler", which calls the ServeHTTP method of the handler, if the midddleware condition is satisfied

// func testGetHandler(res http.ResponseWriter, req *http.Request) {
// 	session, _ := store.Get(req, "session")
// 	untyped, ok := session.Values["username"]
// 	if !ok {
// 		return
// 	}
// 	username, ok := untyped.(string)
// 	if !ok {
// 		return
// 	}
// 	res.Write([]byte(username))
// }

// func goodbyeHandler(res http.ResponseWriter, req *http.Request) {
// 	fmt.Fprint(res, "Goodbye World!")
// }
