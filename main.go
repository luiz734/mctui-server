package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mctui-server/app"
	"net/http"
)

type Command struct {
    Command string `json:"command"`
}
func commandHandler(w http.ResponseWriter, r *http.Request) {
    var cmd Command
    // Decode the JSON body
    err := json.NewDecoder(r.Body).Decode(&cmd)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Log or process the command
    fmt.Fprintf(w, app.AskRconServer(cmd.Command))
}

func main() {
    http.HandleFunc("/command", commandHandler)

    port := 8090
    addr := fmt.Sprintf(":%d", port)
    log.Printf("Listen on port %d", port)
    http.ListenAndServe(addr, nil)
}
