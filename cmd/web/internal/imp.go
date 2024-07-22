package internal

import (
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
