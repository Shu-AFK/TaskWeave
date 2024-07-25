package internal

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password TEXT
		);

		CREATE TABLE IF NOT EXISTS Sessions (
		    sessionId TEXT PRIMARY KEY,
		    userId INTEGER,
		    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY(userId) REFERENCES Users(userId)
		);

		CREATE TABLE IF NOT EXISTS Todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT, 
			description TEXT, 
			deadline TEXT,
			done BOOLEAN CHECK(done IN (0, 1))
		);

		CREATE TABLE IF NOT EXISTS Events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT, 
			duration TEXT,
			deadline TEXT,
			start TEXT,
			end TEXT
		);

		CREATE TABLE IF NOT EXISTS Days (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO Users (username, email, password)
		VALUES (?, ?, ?)
	`, username, email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUser(username string, password string) error {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var storedHashedPassword []byte

	err = db.QueryRow("SELECT password FROM Users WHERE username = ?", username).Scan(&storedHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedHashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid password")
		}
		return err
	}

	return nil
}

func CheckIfSessionExists(userId int) (bool, error) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT * FROM Sessions WHERE userId=?)", userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == true {
		return true, nil
	}

	return false, nil
}

func GetUserIdByName(username string) (int, error) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id int
	err = db.QueryRow("SELECT id FROM Users WHERE username=?", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, errors.New("user does not exist")
		} else {
			return -1, err
		}
	}

	return id, nil
}

func GetUsernameById(userId int) (string, error) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	var username string
	err = db.QueryRow("SELECT username FROM Users WHERE id=?", userId).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user does not exist")
		} else {
			return "", err
		}
	}

	return username, nil
}

func SetSessionID(userId int, sessionID string) error {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT * FROM Sessions WHERE sessionId=?)", sessionID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == true {
		return errors.New("sessionID already exists")
	}

	_, err = db.Exec("INSERT INTO Sessions (sessionId, userId) VALUES (?, ?)", sessionID, userId)
	if err != nil {
		return err
	}

	return nil
}

func GetSessionIDbyUserID(userid int) (string, error) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		return "", nil
	}
	defer db.Close()

	var sessionID string
	err = db.QueryRow("SELECT sessionId FROM Sessions WHERE userId=?", userid).Scan(&sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
