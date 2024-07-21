package main

import (
	"fmt"
	"github.com/Shu-AFK/TaskWeave/cmd/web/handler"
	"github.com/Shu-AFK/TaskWeave/cmd/web/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Creates the DB if it doesn't exist already
	internal.CreateDB()

	// Use Gorilla Mux for routing
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/tasks", handler.TasksHandler)

	// Serve static files (CSS, JS, images)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Start the server
	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
