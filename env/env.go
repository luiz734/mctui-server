package env

import (
	"log"
	"os"
)

func GetJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GetDatabaseFile() string {
	return os.Getenv("DB_FILE")
}

func CheckRequiredVariables() {
	vars := []string{
		"JWT_SECRET",
		"DB_FILE",
	}
	for _, v := range vars {
		// log.Printf("%s %s", v, os.Getenv(v))
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required env var %s", v)
		}
	}
}
