package internal

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Todo struct {
	Name        string
	Description string
	Deadline    time.Time
	Done        bool
}

type Event struct {
	Name     string
	Duration time.Duration
	Deadline time.Time
	Start    time.Time
	End      time.Time
	TodoList []Todo
}

type Day struct {
	Date   time.Time
	Events []Event
}

type EventPage struct {
	Day    Day
	Events []Event
}

// TemplateData holds data passed to templates
type TemplateData struct {
	Title       string
	LinkToTasks string
}

// RenderTemplate renders HTML templates
func RenderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	tmplPath := fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		return
	}
}

func GenerateSessionID() (string, error) {
	byteSize := 32 // Create a byte slice of size 32
	sessionId := make([]byte, byteSize)

	_, err := rand.Read(sessionId)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sessionId), nil
}

func SetSessionCookie(w http.ResponseWriter, username string) error {
	id, err := GetUserIdByName(username)
	if err != nil {
		return err
	}

	exists, err := CheckIfSessionExists(id)
	if err != nil {
		return err
	}

	var sessionID string
	if exists {
		sessionID, err = GetSessionIDbyUserID(id)
	} else {
		sessionID, err = GenerateSessionID()
		if err != nil {
			return err
		}

		err = SetSessionID(id, sessionID)
		if err != nil {
			if err.Error() == "sessionID already exists" {
				for err.Error() == "sessionID already exists" {
					sessionID, err = GenerateSessionID()
					if err != nil {
						return err
					}

					err = SetSessionID(id, sessionID)
					if err == nil {
						break
					}
				}
			} else {
				return err
			}
		}
	}

	cookie := &http.Cookie{
		Name:     "SessionID",
		Value:    sessionID,
		Secure:   true,
		HttpOnly: true,
		Path:     "/tasks",
	}

	http.SetCookie(w, cookie)
	return nil
}

func CheckIfSessionIsCorrect(username string, r *http.Request) (bool, error) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		return false, nil
	}

	userId, err := GetUserIdByName(username)
	if err != nil {
		return false, err
	}

	storedSessionID, err := GetSessionIDbyUserID(userId)
	if err != nil {
		return false, err
	}

	if storedSessionID != cookie.Value {
		return false, nil
	}

	return true, nil
}
