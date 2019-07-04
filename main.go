package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {
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
	// fmt.Fprint(res, "Hello World!")
	templates.ExecuteTemplate(res, "index.html", nil)
}

// func goodbyeHandler(res http.ResponseWriter, req *http.Request) {
// 	fmt.Fprint(res, "Goodbye World!")
// }
