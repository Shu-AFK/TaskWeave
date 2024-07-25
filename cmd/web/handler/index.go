package handler

import (
	"html/template"
	"net/http"
)

// Index handles the homepage request
func Index(w http.ResponseWriter, r *http.Request) {
	tmplPath := "templates/index.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
