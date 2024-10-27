package app

import (
	"fmt"
	"github.com/charmbracelet/log"
	env "mctui-server/environment"

	"github.com/gorcon/rcon"
)

func AskRconServer(command string) string {
	conn, err := rcon.Dial("127.0.0.1:25575", env.GetRconPassword())
	if err != nil {
        return fmt.Sprintf("%s: %s", err.Error(), "is the server down?")
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	log.Infof("Send command to rcon: %s", command)
	if err != nil {
		return fmt.Sprintf("%s", err.Error())
	}

	return response
}
