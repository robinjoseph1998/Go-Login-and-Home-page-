package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore([]byte("secret-password"))
)

func login(w http.ResponseWriter, r *http.Request) {
	var fileName = "login.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		fmt.Println("Error getting session", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	var errMsg string
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		switch errParam {
		case "incorrect_credentials":
			errMsg = "Invalid username or password."
		default:
			errMsg = "An error occurred."
		}
	}

	type Detail struct {
		Name   string
		ErrMsg string
	}

	detail := Detail{Name: "", ErrMsg: errMsg}
	t.ExecuteTemplate(w, fileName, detail)

	if err != nil {
		fmt.Println("Error when executing Template", err)
		return
	}
}

var userDB = map[string]string{
	"robin": "robin123",
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	if userDB[username] == password {
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println("Error getting session", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Redirecting to home page...")

		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?error=incorrect_credentials", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		fmt.Println("Error getting session", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		login(w, r)
	case "/login_submit":
		loginSubmit(w, r)
	case "/home":
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println("Error getting session", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session.Values["authenticated"] != true {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		fileName := "home.html"
		funcMap := map[string]interface{}{
			"upper": strings.ToUpper,
		}
		t, err := template.New(fileName).Funcs(funcMap).ParseFiles(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		t.ExecuteTemplate(w, fileName, "")
		if err != nil {
			fmt.Println(err)
			return
		}
	case "/logout":
		logout(w, r)
	default:
		fmt.Fprintf(w, "Please Login")

	}
}

func main() {
	http.HandleFunc("/login", handler)
	http.HandleFunc("/login_submit", loginSubmit)
	http.HandleFunc("/home", handler) // add this line
	http.HandleFunc("/logout", logout)
	http.ListenAndServe("", nil)
}
