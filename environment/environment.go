package env

import (
	"github.com/charmbracelet/log"
	"os"
)

// Checks for enviroment variables

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
func GetWorldDir() string {
	return os.Getenv("WORLD_DIR")
}
func GetBackupDir() string {
	return os.Getenv("BACKUP_DIR")
}
func GetRconPassword() string {
	return os.Getenv("RCON_PASSWORD")
}
func GetRconAddress() string {
	return os.Getenv("RCON_ADDRESS")
}
func GetPort() string {
	return os.Getenv("PORT")
}

func CheckRequiredVariables() {
	vars := []string{
        "PORT",
		"JWT_SECRET",
		"DB_FILE",
		"TLS_CERT_FILE",
		"TLS_KEY_FILE",
		"WORLD_DIR",
		"BACKUP_DIR",
		"RCON_PASSWORD",
	}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required env var %s", v)
		}
	}
}
