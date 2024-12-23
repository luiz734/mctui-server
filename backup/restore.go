package backup

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"mctui-server/app"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

var ServerLogFile = "/home/tohru/tmp/mcserver.log"

// User wants to restore a backup
func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	// User provides the backup filename
	var err error
	var rr RestoreRequest
	err = json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if server running
	if !serverRunning() {
		errMsg := fmt.Sprintf("Minecraft server not running")
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Check if server is empty
	if !serverEmpty() {
		errMsg := fmt.Sprintf("Minecraft server not empty")
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	// Server empty. We can save now
	// Save the world and close the server
	systemdStop()
	// Restore the backup
	if err = restoreBackup(rr.Filename); err != nil {
		log.Error("Error restoring backup: %v", err)
		log.Error("User IP: %s", r.RemoteAddr)
		http.Error(w, "Forbidden action. IP logged.", http.StatusForbidden)
		defer systemdStart()
		return
	}
	// Starts the server again
	// startServer()
	systemdStart()

	w.Write([]byte("Backup restored!"))
}

// Returns true if server is empty
func serverEmpty() bool {
	output := app.AskRconServer("list")
	return strings.Contains(output, "There are 0 of a max of")
}

// Returns true if server is running
func serverRunning() bool {
	if status, _ := SystemdStatus(); status != Active {
		return false
	}
	return true
}

func checkDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// panic(fmt.Sprintf("dir %s not exists", path))
		return false
	}
	return true
}

type RestoreRequest struct {
	Filename string `json:filename`
}

type badFilenameError struct {
	message string
}

func (e *badFilenameError) Error() string {
	return fmt.Sprintf("invalid filename: %s", e.message)
}

func restoreBackup(backupName string) error {
	/*
		        These variables can be confuse
		        Here are some examples to help

		        server is the server root
		        backups is where you store the backups

		        fileBackupZip       back123.zip
		        fileBackup          back123

				serverRootPath      /server
				currentSavePath     /server/world
				oldSavePath         /server/old
				backupBeforePath    /server/back123
				backupAfterPath     /backups/back123.zip
	*/

	var err error
	// File names
	fileBackupZip := backupName
	fileBackup := strings.TrimSuffix(fileBackupZip, filepath.Ext(fileBackupZip))

	// Directories
	serverRootPath := path.Dir(Dirs.saves)
	currentSavePath := Dirs.saves
	oldSavePath := path.Join(serverRootPath, "old")
	backupAfterPath := path.Join(serverRootPath, fileBackup)
	backupBeforePath := path.Join(Dirs.manual, fileBackupZip)

	// Check for directory traversal
	// E.g filename like ../../../something
	absPath, err := filepath.Abs(backupBeforePath)
	// log.Printf("%s %s", absPath, Dirs.manual)
	if !strings.HasPrefix(absPath, Dirs.manual) {
		return &badFilenameError{"directory traversal attempt detected!"}
	}

	// Remove a dir named the same as the backup
	// Unlikelly to happen. Happens a lot while debbuging tho
	if checkDirExists(backupAfterPath) {
		err = os.RemoveAll(backupAfterPath)
		if err != nil {
			panic(err)
		}
		log.Debugf("Removed file %s", backupAfterPath)
	}

	// Unarchive the backup in the saves dir
	err = archiver.Unarchive(backupBeforePath, serverRootPath)
	if err != nil {
		panic(err)
	}
	log.Debugf("Unarchived file %s to %s", backupBeforePath, serverRootPath)

	// Remove any dir called "old" in saves dir
	// Workaround until debug (see bellow)
	if checkDirExists(oldSavePath) {
		err = os.RemoveAll(oldSavePath)
		if err != nil {
			panic(err)
		}
		log.Debugf("Remove file %s", oldSavePath)
	}

	// Rename "world" to "old"
	err = os.Rename(currentSavePath, oldSavePath)
	if err != nil {
		panic(err)
	}
	log.Debugf("Renamed %s to %s", currentSavePath, oldSavePath)

	// Rename brand new backup to "world"
	// err = cp.Copy(backupBeforePath, currentSavePath, opts)
	err = os.Rename(backupAfterPath, currentSavePath)
	if err != nil {
		panic(err)
	}
	log.Debugf("Renamed %s to %s", backupAfterPath, currentSavePath)

	// We don't remove "old"
	// Can be usefull to undo the last backup
	log.Info("Restore completed")

	return nil
}
