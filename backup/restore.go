package backup

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	// cp "github.com/otiai10/copy"
)

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

func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	// User provides the backup filename
	var err error
	var rr RestoreRequest
	err = json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

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

	// File names
	fileBackupZip := rr.Filename
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
