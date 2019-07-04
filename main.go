package main

import (
	"fmt"
	"net/http"
)

func handler(res http.ResponseWriter, req *http.Request) {
	// print "Hello WOrld" to ur response writer
	fmt.Fprint(res, "Hello World!")
}

func main() {
	//when a request is made, we want to use our handler
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7000", nil)
}
