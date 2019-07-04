package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler).Methods("GET")
	r.HandleFunc("/goodbye", goodbyeHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)

	//without mux
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":7000", nil)
}

func helloHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Hello World!")
}

func goodbyeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Goodbye World!")
}
