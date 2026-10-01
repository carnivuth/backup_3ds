package web

import (
	"html/template"
	"net/http"
	"console_backupper/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	configFilePath := utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml")
	backupConfig := utils.ParseConfig(configFilePath)

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, backupConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
