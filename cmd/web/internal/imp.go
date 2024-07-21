package internal

import (
	"fmt"
	"html/template"
	"net/http"
)

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
