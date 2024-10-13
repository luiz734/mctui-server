package backup

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mholt/archiver/v3"
	cp "github.com/otiai10/copy"
)

var Dirs directories

type directories struct {
	saves  string
	manual string
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return true
}

func NewDirectories(savesPath, manualPath string) directories {
	if !dirExists(savesPath) {
		panic("missing dir")
	}
	if !dirExists(manualPath) {
		panic("missing dir")
	}

	return directories{
		saves:  savesPath,
		manual: manualPath,
	}
}

func (d directories) MakeBackup() error {
	now := time.Now().Format("2006-01-02-15-04-05")
	dirName := fmt.Sprintf("backup-%s", now)
	zipName := fmt.Sprintf("%s.zip", dirName)

	worldPath := d.saves
	backupPath := path.Join(d.manual, dirName)
	zipPath := path.Join(d.manual, zipName)

	var err error

	// Make a copy of the current world
	err = cp.Copy(worldPath, backupPath)
	if err != nil { /* success */
		panic(err)
	}
	log.Printf("copied %s to %s", worldPath, backupPath)

	// Compress
	err = archiver.Archive([]string{backupPath}, zipPath)
	if err != nil {
		panic(err)
	}
	log.Printf("compressed file %s", zipPath)

	// Remove the uncompressed dir
	err = os.RemoveAll(backupPath)
	if err != nil {
        panic(err)
	}
	log.Printf("removing dir %s", backupPath)
	log.Printf("backup complete")

	return nil
}

type Backup struct {}

func BackupHandler(w http.ResponseWriter, r *http.Request) {
	Dirs.MakeBackup()
	fmt.Fprintf(w, "Done")
}
