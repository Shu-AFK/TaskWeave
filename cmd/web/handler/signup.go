// TODO: Add session to signup

package handler

import (
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
	"html/template"
	"log"
	"net/http"
)

type SignupCreds struct {
	Username        string
	Email           string
	Password        string
	PasswordRetyped string
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	creds := SignupCreds{}

	if r.Method == http.MethodPost {
		creds.Username = r.PostFormValue("username")
		creds.Email = r.PostFormValue("email")
		creds.Password = r.PostFormValue("password")
		creds.PasswordRetyped = r.PostFormValue("password_retyped")

		if creds.Password != creds.PasswordRetyped {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Handle Adding to db
		err := internal.AddUser(creds.Username, creds.Email, creds.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error:", err)
			return
		}

		correct, err := internal.CheckIfSessionIsCorrect(creds.Username, r)
		if err != nil {
			log.Println(err)
			return
		}

		if !correct {
			err = internal.SetSessionCookie(w, creds.Username)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("User has been created:", creds.Username)
		http.Redirect(w, r, "/tasks", http.StatusTemporaryRedirect)
	}

	// Else server the site
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
