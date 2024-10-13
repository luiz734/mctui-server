package app;

import (
	"log"
	"github.com/gorcon/rcon"
)

func AskRconServer(command string) string {
	conn, err := rcon.Dial("127.0.0.1:25575", "minecraft")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	if err != nil {
		log.Fatal(err)
	}
	
    return response
}
