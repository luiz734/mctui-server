package env

import (
	"log"
	"os"
)

func GetJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func CheckRequiredVariables() {
	vars := []string{
		"JWT_SECRET",
	}
	for _, v := range vars {
		// log.Printf("%s %s", v, os.Getenv(v))
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required env var %s", v)
		}
	}
}
