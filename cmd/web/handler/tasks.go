package handler

import (
	"fmt"
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
	"html/template"
	"net/http"
	"time"
)

func RenderDays(w http.ResponseWriter, tmpl string, days []internal.Day) {
	tmplPath := fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	todos := []internal.Todo{
		{
			Name:        "TODO",
			Description: "hello, I'm a todo",
			Deadline:    time.Now(),
			Done:        false,
		},
	}

	events := []internal.Event{
		{
			Name:     "ABC",
			Duration: 0,
			Deadline: time.Now(),
			Start:    time.Now(),
			End:      time.Now(),
			TodoList: todos,
		},
	}

	days := []internal.Day{
		{
			Date:   time.Now(),
			Events: events,
		},
		{
			Date:   time.Now().Add(24 * time.Hour),
			Events: events,
		},
	}

	RenderDays(w, "tasks", days)
}
