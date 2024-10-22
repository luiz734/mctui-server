package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mctui-server/app"
	"mctui-server/backup"
	"mctui-server/db"
	env "mctui-server/environment"
	"net/http"
	"os"
	"strings"
)

type Command struct {
	Command string `json:"command"`
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}

		tokenString = tokenString[len("Bearer "):]

		secretKey := env.GetJwtSecret()
		err := app.VerifyToken(secretKey, tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token")
			return
		}

		// Call the next handler if authorization is successful
		next(w, r)
	}
}

// protected
func commandHandler(w http.ResponseWriter, r *http.Request) {
	var cmd Command
	var err error
	// Decode the JSON body
	err = json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	withoutPrefix, found := strings.CutPrefix(cmd.Command, "!")
	if found {
		fmt.Fprintf(w, fmt.Sprintf("%s not implemented yet...", withoutPrefix))
		return
	}
	fmt.Fprintf(w, app.AskRconServer(cmd.Command))
}

func setupLogs() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	} else {
		_ = io.Discard
		// log.SetOutput(io.Discard)
	}
}

func main() {
	// Setup logs
	setupLogs()

	// Check if there are missing environment variables
	env.CheckRequiredVariables()

	// Check for server.jar and other directories/files
	if err := env.CheckRequiredFiles(); err != nil {
		log.Fatalf("Can't setup: %v", err)
	}

	// Setup database
	db.SetupDatabase()

	if status, _ := backup.SystemdStatus(); status != backup.Active {
		log.Fatalf("minecraft service not running")
	}

	serverTLSCert, err := tls.LoadX509KeyPair(env.GetTlsCert(), env.GetTlsKey())
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	port := 8090
	addr := fmt.Sprintf(":%d", port)

	worldDir := env.GetWorldDir()
	backupDir := env.GetBackupDir()
	backup.Dirs = backup.NewDirectories(worldDir, backupDir)

	mux := http.NewServeMux()

	// Protected routes
	mux.HandleFunc("/command", authMiddleware(commandHandler))
	mux.HandleFunc("/backup", authMiddleware(backup.MakeBackupHandler))
	mux.HandleFunc("/backups", authMiddleware(backup.BackupHandler))
	mux.HandleFunc("/restore", authMiddleware(backup.RestoreHandler))
	// mux.HandleFunc("/restore", backup.RestoreHandler)

	// Public routes
	mux.HandleFunc("/login", app.LoginHandler)

	server := http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
	defer server.Close()

	log.Printf("Listen on port %d", port)
	server.ListenAndServeTLS("", "")
	http.ListenAndServe(addr, nil)
}
