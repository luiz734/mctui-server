package env

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/ini.v1"
)

func CheckRequiredFiles() error {
	if err := checkServerJarFile(); err != nil {
		return fmt.Errorf("check server.jar: %w", err)
	}

	if err := checkEula(); err != nil {
		return fmt.Errorf("check eula: %w", err)
	}

	if err := checkServiceFile(); err != nil {
		return fmt.Errorf("check services: %w", err)
	}

	if err := checkBackupsDir(); err != nil {
		return fmt.Errorf("check backups: %w", err)
	}

	if err := checkServerProperties(); err != nil {
		return fmt.Errorf("check server properties: %w", err)
	}

	return nil
}

func checkServerJarFile() error {
	worldDir := GetWorldDir()
	serverDir := path.Dir(worldDir)
	serverJar := path.Join(serverDir, "server.jar")

	if _, err := os.Stat(serverJar); err != nil {
		return fmt.Errorf("missing server executable: %s", serverJar)
	}
	return nil
}

func checkEula() error {
	worldDir := GetWorldDir()
	serverDir := path.Dir(worldDir)
	eulaFile := path.Join(serverDir, "eula.txt")

	// Check eula.txt exists
	var err error
	if _, err = os.Stat(eulaFile); err != nil {
		return fmt.Errorf("missing eula: %s\nDid you run the server once?", eulaFile)
	}

	// Check eula=true
	var f *os.File
	if f, err = os.Open(eulaFile); err != nil {
		return fmt.Errorf("can't read file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if lastLine != "eula=true" {
		return fmt.Errorf("you must agree with the eula")
	}
	return nil
}

func checkServiceFile() error {
	user := os.Getenv("USER")
	servicesPath := fmt.Sprintf("/home/%s/.config/systemd/user/minecraft.service", user)

	// Check minecraft.service exists
	var err error
	if _, err = os.Stat(servicesPath); err != nil {
		return fmt.Errorf("missing %s\nFile must be created manually.", servicesPath)
	}

	return nil
}

func checkBackupsDir() error {
	backupsDir := GetBackupDir()

	// Check backups dir exists
	var err error
	var f fs.FileInfo
	if f, err = os.Stat(backupsDir); err != nil {
		return fmt.Errorf("missing %s\nDir must be created manually.", backupsDir)
	}
	// Check if it's a dir
	if !f.IsDir() {
		return fmt.Errorf("%s exists but it's not a dir", backupsDir)
	}
	return nil
}

func checkServerProperties() error {
	worldDir := GetWorldDir()
	serverDir := path.Dir(worldDir)
	serverProps := path.Join(serverDir, "server.properties")

	cfg, err := ini.Load(serverProps)
	if err != nil {
		return fmt.Errorf("can't read server config file: %w", err)
	}

	// Check rcon enabled
	enableRcon := cfg.Section("").Key("enable-rcon").MustBool(true)
	if !enableRcon {
		return fmt.Errorf("enable-rcon=false (must be true)")
	}

	// Check port
	_, port, err := splitRconAddr()
	if err != nil {
		return fmt.Errorf("can't parse config: %w", err)
	}
	cfgPort := cfg.Section("").Key("rcon.port").String()
	if port != cfgPort {
		return fmt.Errorf("port in server.properties must match port in env RCON_ADDRESS\nserver.properties: %s\nRCON_ADDRESS: %s", cfgPort, GetRconAddress())
	}

	// Check password
	cfgPassword := cfg.Section("").Key("rcon.password").String()
	if GetRconPassword() != cfgPassword {
		return fmt.Errorf("password in server.properties must match env RCON_PASSWORD")
	}

	// Check server ip
	cfgServerIp := cfg.Section("").Key("server-ip").String()
	if cfgServerIp == "" {
		log.Printf("WARNING: server-ip not set")
        log.Printf("The server will only be visible in you local network")
        log.Printf("To make it avaliable set it to 0.0.0.0")
	}

	return nil
}

func splitRconAddr() (addr string, port string, err error) {
	s := strings.Split(GetRconAddress(), ":")
	if len(s) != 2 {
		return "", "", fmt.Errorf("invalid address: %s", GetRconAddress())
	}
	return s[0], s[1], nil
}
