package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type taskJson struct {
	TaskName string `json:"task"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var t taskJson
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Printf("Can't parse json: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Printf("Got task: %s", t.TaskName)
	task := t.TaskName
	switch task {
	case "start":
		StartHandler(w, r)
		return
	case "stop":
		StopHandler(w, r)
		return
	default:
		errMsg := fmt.Sprintf("Unknown task %s", task)
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

// !start
func StartHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var state ServiceState

	// Check errors on checking status
	if state, err = SystemdStatus(); err != nil {
		log.Printf("Can't get minecraft server status: %v", err)
		http.Error(w, "Something wrong. Check the server logs", http.StatusInternalServerError)
		return
	}
	// Check server is already running
	if state == Active {
		errMsg := fmt.Sprintf("Server already running")
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	log.Printf("Server will start")
	systemdStart()
	w.Write([]byte("Server started"))
}

// !stop
func StopHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var state ServiceState

	// Check errors on checking status
	if state, err = SystemdStatus(); err != nil {
		log.Printf("Can't get minecraft server status: %v", err)
		http.Error(w, "Something wrong. Check the server logs", http.StatusInternalServerError)
		return
	}
	// Check server is already running
	if state == Inactive {
		errMsg := fmt.Sprintf("Server already stopped")
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	log.Printf("Server will stop")
	systemdStop()
	w.Write([]byte("Server stopped"))
}
