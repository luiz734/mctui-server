package db

import (
	"database/sql"
	"fmt"
	"log"
	env "mctui-server/environment"

	_ "github.com/mattn/go-sqlite3"
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
	log.Printf("Database initialized")

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
			log.Printf("Username %s not found", username)
			return false, nil
		}
		log.Printf("Unknown error: %v", err)
		return false, err
	}

	if password == hashedPassword {
		log.Printf("Credentials of user %s match", username)
		return true, nil
	}
	log.Printf("Invalid password for user %s", username)

	return false, nil
}

func initDB() error {
	db, err := sql.Open("sqlite3", env.GetDatabaseFile())
	if err != nil {
		log.Printf("Missing database file. It will be created")
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
