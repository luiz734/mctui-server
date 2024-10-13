package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"mctui-server/app"
	"net/http"
	"strings"
)

type Command struct {
	Command string `json:"command"`
}

var (
	CertFilePath = "cert/cert.pem"
	KeyFilePath  = "cert/key.pem"
)

func commandHandler(w http.ResponseWriter, r *http.Request) {
	var cmd Command
	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&cmd)
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
	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	port := 8090
	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/command", commandHandler)

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
