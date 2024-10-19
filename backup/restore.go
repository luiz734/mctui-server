package backup

import (
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"mctui-server/app"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	// cp "github.com/otiai10/copy"
)

var ServerLogFile = "/home/tohru/tmp/mcserver.log"

func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	// todo: filter out player chat

	// User provides the backup filename
	var err error
	var rr RestoreRequest
	err = json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if server is empty
	if !serverEmpty() {
		errMsg := fmt.Sprintf("server not empty")
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	// Server empty. We can save now
	// saveAndWait()
	// Close the server
	// stopAndWait()
    // Save the world and close the server
	systemdStop()
	// Restore the backup
	restoreBackup(rr.Filename)
	// Starts the server again
	// startServer()

	w.Write([]byte("Backup restored!"))
}

func startServer() {
	// todo: check server is not running
	serverDir := "/home/tohru/tmp/minecraft-server"
	originalDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Change to new directory
	if err := os.Chdir(serverDir); err != nil {
		panic(err)
	}
	// Revert back to the original directory
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			panic(err)
		}
	}()

	startScript := "./run.sh"
	cmd := exec.Command("/bin/bash", "-c", startScript)

	// Detach the process so it won't terminate when the Go program exits
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Optional: Redirect output to avoid hanging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start() // Start the command but don't wait for it to finish
	if err != nil {
		panic(err)
	}

}

func stopAndWait() {
	command := "stop"
	patternInLog := "All dimensions are saved"
	go func() { _ = app.AskRconServer(command) }()
	_ = <-waitForPattern(ServerLogFile, patternInLog, true)
	return

}
func saveAndWait() {
	command := "save-all"
	patternInLog := "Saved the game"
	// Output not usefull in this case
	// Save may take some time
	// We watch the log file intead
	go func() { _ = app.AskRconServer(command) }()
	// Now we wait. When the pattern appears, the saving is done
	_ = <-waitForPattern(ServerLogFile, patternInLog, true)
	return
}

// Returns true if server is empty
func serverEmpty() bool {
	output := app.AskRconServer("list")
	return strings.Contains(output, "There are 0 of a max of")
}

func checkDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// panic(fmt.Sprintf("dir %s not exists", path))
		return false
	}
	return true
}

var patternServerStopped = "All dimensions are saved"

func waitForPattern(logfile, pattern string, seekAtEnd bool) chan bool {
	seek := io.SeekStart
	if seekAtEnd {
		seek = io.SeekEnd
	}
	_ = seek
	config := tail.Config{
		Follow:    true,
		MustExist: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: seek,
		},
	}

	t, err := tail.TailFile(logfile, config)
	if err != nil {
		panic(err)
	}

	finished := make(chan bool)
	go func() {
		for line := range t.Lines {
			log.Printf("new-line: %s", line.Text)
			if strings.Contains(line.Text, pattern) {
				finished <- true
			}
		}
	}()
	return finished
}

type RestoreRequest struct {
	Filename string `json:filename`
}

func restoreBackup(backupName string) {
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

	// Remove a dir named the same as the backup
	// Unlikelly to happen. Happens a lot while debbuging tho
	if checkDirExists(backupAfterPath) {
		err = os.RemoveAll(backupAfterPath)
		if err != nil {
			panic(err)
		}
		log.Printf("removed file %s", backupAfterPath)
	}

	// Unarchive the backup in the saves dir
	err = archiver.Unarchive(backupBeforePath, serverRootPath)
	if err != nil {
		panic(err)
	}
	log.Printf("unarchived file %s to %s", backupBeforePath, serverRootPath)

	// Remove any dir called "old" in saves dir
	// Workaround until debug (see bellow)
	if checkDirExists(oldSavePath) {
		err = os.RemoveAll(oldSavePath)
		if err != nil {
			panic(err)
		}
		log.Printf("remove file %s", oldSavePath)
	}

	// Good. Now we can recreate an empty dir there
	// err = os.Mkdir(oldSavePath, os.ModePerm)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// log.Printf("created dir %s", oldSavePath)

	// Option to replace destination
	// This is not working as expected for some reason
	// Manually removing the dir in the code above
	// var opts = cp.Options{
	// 	OnDirExists: func(src, dest string) cp.DirExistsAction {
	// 		return cp.Replace
	// 	},
	// }
	// err = cp.Copy(currentSavePath, oldSavePath, opts)
	// Rename "world" to "old"
	err = os.Rename(currentSavePath, oldSavePath)
	if err != nil {
		panic(err)
	}
	log.Printf("renamed %s to %s", currentSavePath, oldSavePath)

	// Rename brand new backup to "world"
	// err = cp.Copy(backupBeforePath, currentSavePath, opts)
	err = os.Rename(backupAfterPath, currentSavePath)
	if err != nil {
		panic(err)
	}
	log.Printf("renamed %s to %s", backupAfterPath, currentSavePath)

	// We don't remove "old"
	// Can be usefull to undo the last backup
	log.Printf("restore completed")

}
