package app

import (
	"fmt"
	"log"

	"github.com/gorcon/rcon"
)

func AskRconServer(command string) string {
	conn, err := rcon.Dial("127.0.0.1:25575", "minecraft")
	if err != nil {
		return fmt.Sprintf("%s\n%s", err.Error(), "Is the server down?")
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	log.Printf("Send command to rcon: %s", command)
	if err != nil {
		return fmt.Sprintf("%s", err.Error())
	}

	return response
}
