package handler

import (
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
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
		err := internal.ValidateUser(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Println(err.Error())
			return
		}

		// TODO: Fix redirect
		log.Printf("Login success for %s\n", creds.Username)
		http.Redirect(w, r, "/tasks", http.StatusTemporaryRedirect)
	}

	// If not a POST request
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
