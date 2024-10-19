package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"mctui-server/app"
	"mctui-server/backup"
	"net/http"
	"os"
	"strings"
)

type Command struct {
	Command string `json:"command"`
}

var (
	CertFilePath = "cert/cert.pem"
	KeyFilePath  = "cert/key.pem"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}

		tokenString = tokenString[len("Bearer "):]

		// Assuming `app.VerifyToken` is your token verification function
		err := app.VerifyToken(tokenString)
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
	// w.Header().Set("Content-Type", "application/json")
	//
	// tokenString := r.Header.Get("Authorization")
	// if tokenString == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprint(w, "Missing authorization header")
	// 	return
	// }
	// tokenString = tokenString[len("Bearer "):]
	//
	// err := app.VerifyToken(tokenString)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprint(w, "Invalid token")
	// 	return
	// }

	// User authenticated

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

func main() {
	if status, _ := backup.SystemdStatus(); status != backup.Active {
		log.Fatalf("minecraft service not running")
	}

	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	port := 8090
	addr := fmt.Sprintf(":%d", port)

	envUser := os.Getenv("USER")
	worldDir := fmt.Sprintf("/home/%s/tmp/minecraft-server/world", envUser)
	backupDir := fmt.Sprintf("/home/%s/tmp/backups", envUser)
	backup.Dirs = backup.NewDirectories(worldDir, backupDir)

	mux := http.NewServeMux()

	// Protected routes
	mux.HandleFunc("/command", authMiddleware(commandHandler))
	mux.HandleFunc("/backup", authMiddleware(backup.MakeBackupHandler))
	mux.HandleFunc("/backups", authMiddleware(backup.BackupHandler))
	// mux.HandleFunc("/restore", authMiddleware(backup.RestoreHandler))
	mux.HandleFunc("/restore", backup.RestoreHandler)

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
