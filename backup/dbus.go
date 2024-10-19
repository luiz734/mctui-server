package backup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	// "github.com/godbus/dbus/v5"
)

func systemdStop() {
	context := context.Background()
	conn, err := dbus.NewUserConnectionContext(context)
	if err != nil {
		log.Fatalf("Failed to connect to systemd: %v", err)
	}
	defer conn.Close()
	log.Printf("connected to dbus")

	// Stop the service
	service := "minecraft.service"
	stopRes, err := conn.StopUnitContext(context, service, "replace", nil)
	if err != nil {
		log.Fatalf("Failed to stop service: %v", err)
	}
	log.Printf("stop result: %s", stopRes)

	// Wait for the service to reach the "active" state (it should be active after a successful restart)
	targetState := "inactive"
	if err := waitForServiceState(conn, service, targetState); err != nil {
		log.Fatalf("Error waiting for service to become active: %v", err)
	}
	fmt.Println("Service has restarted successfully and is now active.")

}

func waitForServiceState(conn *dbus.Conn, service string, targetState string) error {
	for {
		// Get the service's current properties (like ActiveState)
		context := context.Background()
		props, err := conn.GetUnitPropertiesContext(context, service)
		if err != nil {
			return fmt.Errorf("failed to get unit properties: %v", err)
		}

		// Check the current "ActiveState"
		currentState := props["ActiveState"].(string)
		fmt.Printf("Current service state: %s\n", currentState)

		// If it matches the target state (e.g., "active" after a restart), we're done
		if currentState == targetState {
			break
		}

		// Wait for a short while before checking again
		time.Sleep(2 * time.Second)
	}

	return nil
}
