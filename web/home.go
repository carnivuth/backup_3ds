package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"console_backupper/utils"
)

type HomeData struct {
	ConsoleNumber int
	Consoles []os.DirEntry
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get consoles
	datadir := utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")
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
