package handler

import (
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
	"net/http"
)

// Index handles the homepage request
func Index(w http.ResponseWriter, r *http.Request) {
	data := internal.TemplateData{
		Title:       "Welcome",
		LinkToTasks: "/tasks",
	}

	internal.RenderTemplate(w, "index", data)
}
