package backup

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	// "mctui-server/app"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
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
	now := time.Now().UTC().Format("2006-01-02-15-04-05")
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
	log.Debug("Copied %s to %s", worldPath, backupPath)

	// Compress
	err = archiver.Archive([]string{backupPath}, zipPath)
	if err != nil {
		panic(err)
	}
	log.Debug("Compressed file %s", zipPath)

	// Remove the uncompressed dir
	err = os.RemoveAll(backupPath)
	if err != nil {
		panic(err)
	}
	log.Debug("Removing dir %s", backupPath)
	log.Debug("Backup complete")

	return nil
}

func MakeBackupHandler(w http.ResponseWriter, r *http.Request) {
	Dirs.MakeBackup()
	fmt.Fprintf(w, "Done")
}

func BackupHandler(w http.ResponseWriter, r *http.Request) {
	backups, err := Dirs.LoadBackups()
	if err != nil {
		log.Errorf("Error listing backups: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var backupNames []string
	for _, b := range backups {
		backupNames = append(backupNames, b.Name)
	}
	json.NewEncoder(w).Encode(backupNames)
}

type Backup struct {
	Time time.Time
	Name string
}

func (d directories) LoadBackups() ([]Backup, error) {
	var backups []Backup

	files, err := os.ReadDir(d.manual)
	if err != nil {
        return nil, fmt.Errorf("error reading dir content: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			log.Error("Found directory in backups folder", "dirname", file.Name())
		}

		name := file.Name()
		if strings.HasPrefix(name, "backup-") && strings.HasSuffix(name, ".zip") {
			timestamp := name[len("backup-") : len(name)-len(".zip")] // extract time part
			t, err := time.Parse("2006-01-02-15-04-05", timestamp)
			if err != nil {
                log.Errorf("Error parsing time: %v", err)
			}
			backups = append(backups, Backup{Time: t, Name: name})
		} else {
			log.Error("File in backups dir with bad naming", "filename", file.Name())
		}
	}

	// Sort backups by time
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Time.After(backups[j].Time)
	})

	return backups, nil
}
