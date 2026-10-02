package web

import (
	"console_backupper/model"
	"console_backupper/utils"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	configFilePath := utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml")
	backupConfig := model.GetBackupConfig(configFilePath)

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, backupConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
