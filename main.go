package main

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client

var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //redis port
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	// r.HandleFunc("/goodbye", goodbyeHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)

	//without mux
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":7000", nil)
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	//get the first ten strings in redis from a string called "comments":
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	// fmt.Fprint(res, "Hello World!")
	templates.ExecuteTemplate(res, "index.html", comments)
}

// func goodbyeHandler(res http.ResponseWriter, req *http.Request) {
// 	fmt.Fprint(res, "Goodbye World!")
// }
