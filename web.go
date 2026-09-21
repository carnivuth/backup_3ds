package console_backupper

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type HomeData struct {
	ConsoleNumber int
	Consoles []os.DirEntry
}

type BackupData struct {
	BackupNumber int
	Backups []os.DirEntry
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get consoles
	datadir := getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")
	consoles, err := os.ReadDir(datadir)
	if err != nil {
		log.Fatal(err)
	}
	ConsoleNumber := len(consoles)
	data := HomeData{ConsoleNumber: ConsoleNumber, Consoles: consoles }
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func backupHandler(w http.ResponseWriter, r *http.Request) {

	console := r.PathValue("console")
	tmpl, err := template.ParseFiles("templates/backup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datadir := getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")
	backups, err := os.ReadDir(filepath.Join(datadir, console))
	if err != nil {
		log.Fatal(err)
	}
	backupNumber := len(backups)
	data := BackupData{ BackupNumber: backupNumber, Backups: backups}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
