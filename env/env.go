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
func GetTlsCert() string {
	return os.Getenv("TLS_CERT_FILE")
}
func GetTlsKey() string {
	return os.Getenv("TLS_KEY_FILE")
}

func CheckRequiredVariables() {
	vars := []string{
		"JWT_SECRET",
		"DB_FILE",
		"TLS_CERT_FILE",
		"TLS_KEY_FILE",
	}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required env var %s", v)
		}
	}
}
