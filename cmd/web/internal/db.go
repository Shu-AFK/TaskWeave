package internal

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func CreateDB() {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Users (
			id INTEGER PRIMARY KEY,
			username TEXT,
			email TEXT,
			password TEXT
		);

		CREATE TABLE IF NOT EXISTS Todos (
			id INTEGER PRIMARY KEY, 
			name TEXT, 
			description TEXT, 
			deadline TEXT,
			done BOOLEAN CHECK(done IN (0, 1))
		);

		CREATE TABLE IF NOT EXISTS Events (
			id INTEGER PRIMARY KEY,
			name TEXT, 
			duration TEXT,
			deadline TEXT,
			start TEXT,
			end TEXT
		);

		CREATE TABLE IF NOT EXISTS Days (
			id INTEGER PRIMARY KEY,
			startOfDay TEXT, 
			endOfDay TEXT, 
			duration TEXT,
			userId INTEGER,
			FOREIGN KEY(userId) REFERENCES Users(id)
		);

		CREATE TABLE IF NOT EXISTS EventTodos (
			eventId INTEGER,
			todoId INTEGER,
			FOREIGN KEY(eventId) REFERENCES Events(id),
			FOREIGN KEY(todoId) REFERENCES Todos(id)
		);

		CREATE TABLE IF NOT EXISTS DayEvents (
			dayId INTEGER,
			eventId INTEGER,
			FOREIGN KEY(dayId) REFERENCES Days(id),
			FOREIGN KEY(eventId) REFERENCES Events(id)
		);		
	`)
	if err != nil {
		log.Fatal(err)
	}
}
