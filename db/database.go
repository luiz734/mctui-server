package db

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	env "mctui-server/environment"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string
	Username string
	Password string
}

func SetupDatabase() error {
	//    log.Printf("hello")
	// if _, err := os.Stat(env.GetDatabaseFile()); err != nil {
	// 	log.Printf("Missing database file. It will be created")
	// } else {
	// 	log.Printf("Found database file")
	// }
	//
	if err := initDB(); err != nil {
		return fmt.Errorf("can't init database: %w", err)
	}
	log.Info("Database initialized")

	return nil
}

func CheckCredentials(username, password string) (bool, error) {
	// Open the database
	db, err := sql.Open("sqlite3", env.GetDatabaseFile())
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Prepare the query
	query := `SELECT password FROM users WHERE username = ?`

	var hashedPassword string
	err = db.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Username %s not found", username)
			return false, nil
		}
		log.Error("Unknown error: %v", err)
		return false, err
	}

	if checkPassword(password, hashedPassword) {
		log.Debug("Credentials of user %s match", username)
		return true, nil
	}
	log.Debug("Invalid password for user %s", username)

	return false, nil
}

func initDB() error {
	db, err := sql.Open("sqlite3", env.GetDatabaseFile())
	if err != nil {
		log.Info("Missing database file. It will be created")
		// log.Fatalf("Can't open database file")
		// return err
	}
	defer db.Close()

	// Ensure the table is created if it doesn't exist
	err = createUserTable(db)
	if err != nil {
		return err
	}
	return nil
}

func createUserTable(db *sql.DB) error {
	const create string = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	_, err := db.Exec(create)
	return err
}

func AddUser(username, password string) error {
	var err error
	var hashedPassword string
	if hashedPassword, err = hashPassword(password); err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	// Opens database
	db, err := sql.Open("sqlite3", env.GetDatabaseFile())
	if err != nil {
		return fmt.Errorf("can't open database file: %w", err)
	}
	defer db.Close()
	// Insert user
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err = db.Exec(query, username, hashedPassword)
	if err != nil {
		return err
	}
	log.Info("User added successfully")
	return nil
}

func GetAllUsernames() ([]string, error) {
	var usernames []string
	// Opens database
	db, err := sql.Open("sqlite3", env.GetDatabaseFile())
	if err != nil {
		return usernames, fmt.Errorf("can't open database file: %w", err)
	}
	defer db.Close()

	// Query usernames
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return usernames, nil
}

// Hash the password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Check the password using the hash (salt is extracted from the hash)
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
