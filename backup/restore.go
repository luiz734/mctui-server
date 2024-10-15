package backup

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	cp "github.com/otiai10/copy"
)

func checkDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// panic(fmt.Sprintf("dir %s not exists", path))
		return false
	}
	return true
}

func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	serverDir := path.Dir(Dirs.saves)
	backupBaseName := "backup-2024-10-13-21-54-21.zip"
	backupName := strings.TrimSuffix(backupBaseName, filepath.Ext(backupBaseName))
	backupPath := path.Join(Dirs.manual, backupBaseName)
	// Filename is not necessary. Only dir
	restoreDir := path.Join(serverDir, "")
	restorePath := path.Join(restoreDir, backupName)
	currentSavePath := Dirs.saves
	oldSavePath := path.Join(serverDir, "old")

	// oldWorldPath := "old"
	// newWorldPath := "world"

	// check backups exists

	var err error
    // Remove a dir named the same as the backup
    // Unlikelly to happen. Happens a lot while debbuging tho
	if checkDirExists(restorePath) {
		err = os.RemoveAll(restorePath)
		if err != nil {
			panic(err)
		}
		log.Printf("remove file %s", restorePath)
	}

    // Unarchive the backup in the saves dir
	err = archiver.Unarchive(backupPath, restoreDir)
	if err != nil {
		panic(err)
	}
	log.Printf("unarchived file %s to %s", backupPath, restoreDir)

    // Remove any dir called "old" in saves dir
    // Workaround until debug (see bellow)
	if checkDirExists(oldSavePath) {
		err = os.RemoveAll(oldSavePath)
		if err != nil {
			panic(err)
		}
		log.Printf("remove file %s", restorePath)
	}
    
    // Good. Now we can recreate an empty dir there
	err = os.Mkdir(oldSavePath, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("created dir %d", oldSavePath)

    // Option to replace destination
    // This is not working as expected for some reason
    // Manually removing the dir in the code above
	var opts = cp.Options{
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			return cp.Replace
		},
	}
	err = cp.Copy(currentSavePath, oldSavePath, opts)
	if err != nil {
		panic(err)
	}
	log.Printf("copied %s to %s", currentSavePath, oldSavePath)

	// Rename brand new backup to "world"
	err = cp.Copy(restorePath, currentSavePath)
	if err != nil {
		panic(err)
	}
	log.Printf("renamed %s to %s", restorePath, currentSavePath)

    // We don't remove "old"
    // Can be usefull to undo the last backup
	log.Printf("restore completed")

}
