// TODO: Manage sessions when logging in

package handler

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
	"html/template"
	"log"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

func generateSessionID() (string, error) {
	byteSize := 32 // Create a byte slice of size 32
	sessionId := make([]byte, byteSize)

	_, err := rand.Read(sessionId)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sessionId), nil
}

func setSessionCookie(w http.ResponseWriter, username string) error {
	id, err := internal.GetUserIdByName(username)
	if err != nil {
		return err
	}

	exists, err := internal.CheckIfSessionExists(id)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "SessionID",
		Value:    sessionID,
		Secure:   true,
		HttpOnly: true,
		Path:     "/tasks",
	}

	err = internal.SetSessionID(id, sessionID)
	if err != nil {
		if err.Error() == "sessionID already exists" {
			for err.Error() == "sessionID already exists" {
				sessionID, err = generateSessionID()
				if err != nil {
					return err
				}

				cookie.Value = sessionID
				err = internal.SetSessionID(id, sessionID)
				if err == nil {
					break
				}
			}
		} else {
			return err
		}
	}

	http.SetCookie(w, cookie)
	return nil
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
			goto End
		}

		// TODO: Fix redirect
		err = setSessionCookie(w, creds.Username)
		if err != nil {
			log.Println(err)
			goto End
		}

		log.Printf("Login success for %s\n", creds.Username)
		http.Redirect(w, r, "/tasks", http.StatusTemporaryRedirect)
	}

End:
	// If not a POST request
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
