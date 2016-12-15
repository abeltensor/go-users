package main

import (
	"crypto/rand"
	"html/template"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

type User struct {
	Uuid     string
	FName    string
	LName    string
	UserName string
	Email    string
	Password string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "header.html", "base.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "main", u)
	case "POST":

		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("email")
		un := r.FormValue("userName")
		p := r.FormValue("password")
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			fmt.Printf("can't initialize UUID")
			return
		}
		uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		u := &User{Uuid: uuid, FName: f, LName: l, UserName: un, Email: em, Password: p}
		setSession(u, w)
		err = saveData(u)
		if err != nil {
			fmt.Printf("unable to save data in DB")
		}
		redirectTarget := "/internal"
		http.Redirect(w, r, redirectTarget, 302)
	}

}
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	u := &User{}
	tmpl, _ := template.ParseFiles("index.html", "header.html", "base.html")
	err := tmpl.ExecuteTemplate(w, "main", u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "internal.html", "header.html")
	username := getUserName(r)

	if username != "" {
		err := tmpl.ExecuteTemplate(w, "main", &User{UserName: username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pass := r.FormValue("password")

	redirectTarget := "/"
	if name != "" && pass != "" {
		u, err := loadUser(name, pass)
		if err != nil {
			http.Redirect(w, r, "/", 302)
		}
		setSession(u, w)
		redirectTarget = "/internal"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/signup", signUpHandler).Methods("GET", "POST")
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/internal", internalPageHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
