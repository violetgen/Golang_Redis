package routes

import (
	"log"
	"net/http"

	"../middleware"
	"../models"
	"../utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("victor-steven"))

func NewRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	// r.HandleFunc("/test", testGetHandler).Methods("GET")
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return r
}

func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	// session, _ := store.Get(req, "session")
	// _, ok := session.Values["username"]
	// if !ok {
	// 	http.Redirect(res, req, "/login", 302)
	// 	return
	// }

	//get the first ten strings in redis from a string called "comments":
	comments, err := models.GetComments()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	// fmt.Fprint(res, "Hello World!")
	utils.ExecuteTemplate(res, "index.html", comments)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
	//parse the form from the request body
	req.ParseForm()

	comment := req.PostForm.Get("comment")

	log.Println(comment)

	err := models.PostComment(comment)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	//if no errors, return to the main page
	http.Redirect(res, req, "/", 302)

	// fmt.Fprint(res, "Hello World!")
	// utils.ExecuteTemplate(res, "index.html", comments)
}

func loginGetHandler(res http.ResponseWriter, req *http.Request) {
	utils.ExecuteTemplate(res, "login.html", nil)
}

func loginPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(res, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(res, "login.html", "Invalid login details")
		default:
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Internal server error"))
		}
		return
	}

	session, _ := Store.Get(req, "session")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Save(req, res)
	http.Redirect(res, req, "/", 302)
}

func registerGetHandler(res http.ResponseWriter, req *http.Request) {
	utils.ExecuteTemplate(res, "register.html", nil)
}

func registerPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")

	err := models.RegisterUser(username, password)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(res, req, "/login", 302)
}
