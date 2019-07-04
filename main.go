package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var client *redis.Client
var store = sessions.NewCookieStore([]byte("victor-steven"))

var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //redis port
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	// r.HandleFunc("/test", testGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// r.HandleFunc("/goodbye", goodbyeHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)

	//without mux
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":7000", nil)
}

func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	_, ok := session.Values["username"]
	if !ok {
		http.Redirect(res, req, "/login", 302)

		return
	}
	//get the first ten strings in redis from a string called "comments":
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	// fmt.Fprint(res, "Hello World!")
	templates.ExecuteTemplate(res, "index.html", comments)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
	//parse the form from the request body
	req.ParseForm()

	comment := req.PostForm.Get("comment")

	log.Println(comment)

	err := client.LPush("comments", comment).Err()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	//if no errors, return to the main page
	http.Redirect(res, req, "/", 302)

	// fmt.Fprint(res, "Hello World!")
	// templates.ExecuteTemplate(res, "index.html", comments)
}

func loginGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "login.html", nil)
}

func loginPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(res, "login.html", "unknown user")
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}

	//check if the hash the user entered matched with the one stored:
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		templates.ExecuteTemplate(res, "login.html", "Invalid login details")
		return
	}
	session, _ := store.Get(req, "session")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Save(req, res)
	http.Redirect(res, req, "/", 302)
}

func registerGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "register.html", nil)
}

func registerPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")

	//how much strength the password have
	cost := bcrypt.DefaultCost
	//convert the password from string to byte
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	//Add the user to redis
	//the zero tells the set method that the key should not expire
	err = client.Set("user:"+username, hash, 0).Err()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(res, req, "/login", 302)
}

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
