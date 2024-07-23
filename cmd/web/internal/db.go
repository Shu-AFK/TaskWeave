package internal

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/mail"
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
			userId INTEGER,
			date TEXT,
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

func emailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func valueExistsUserDB(db *sql.DB, isUsername bool, toCheck string) (bool, error) {
	var exists bool
	var query string
	if isUsername {
		query = fmt.Sprintf(`SELECT EXISTS(SELECT * FROM Users WHERE %s=?)`, "username")
	} else {
		query = fmt.Sprintf(`SELECT EXISTS(SELECT * FROM Users WHERE %s=?)`, "email")
	}

	err := db.QueryRow(query, toCheck).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func AddUser(username string, email string, password string) error {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if email is valid
	if username == "" || password == "" || email == "" {
		return errors.New("username, password or email is empty")
	}
	if !emailValid(email) {
		return errors.New("invalid email")
	}

	exists, err := valueExistsUserDB(db, true, username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	exists, err = valueExistsUserDB(db, false, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	_, err = db.Exec(`
		INSERT INTO Users (username, email, password)
		VALUES (?, ?, ?)
	`, username, email, password)
	if err != nil {
		return err
	}

	return nil
}
