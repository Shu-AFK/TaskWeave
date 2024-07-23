package handler

import (
	"html/template"
	"log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	if r.Method == http.MethodPost {
		creds.Username = r.PostFormValue("username")
		creds.Password = r.PostFormValue("password")

		// Do authentication
		log.Printf("Username: %s, Password: %s", creds.Username, creds.Password)
	}

	// If not a POST request
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
