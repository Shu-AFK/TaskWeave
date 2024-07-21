package handler

import (
	"fmt"
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
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
	StartOfDay time.Time
	EndOfDay   time.Time
	Duration   time.Duration
	Events     []Event
}

var parseTemplates = []string{"Jan 2 15:04", "Jan 2 3pm"}

func tryParseTime(input string, layouts []string) (time.Time, string, error) {
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, input)
		if err == nil {
			return parsedTime, layout, nil
		}
	}

	return time.Time{}, "", fmt.Errorf("no layouts matched")
}

func checkIfInDays(timePoint time.Time, days []*Day) bool {
	for _, day := range days {
		if day.StartOfDay.Day() == timePoint.Day() && day.StartOfDay.Month() == timePoint.Month() {
			return true
		}
	}
	return false
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	data := internal.TemplateData{
		Title:       "Welcome",
		LinkToTasks: "/tasks",
	}

	internal.RenderTemplate(w, "tasks", data)
}
