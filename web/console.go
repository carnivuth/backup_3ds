package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"console_backupper/utils"
)


type ConsoleData struct {
	ConsoleNumber int
	Consoles []os.DirEntry
}


func ConsoleHandler(w http.ResponseWriter, r *http.Request) {

	console := r.PathValue("console")
	tmpl, err := template.ParseFiles("templates/console.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datadir := utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")
	backups, err := os.ReadDir(filepath.Join(datadir, console))
	if err != nil {
		log.Fatal(err)
	}
	backupNumber := len(backups)
	data := ConsoleData{ ConsoleNumber: backupNumber, Consoles: backups}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
