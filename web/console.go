package web

import (
	"console_backupper/model"
	"console_backupper/utils"
	"html/template"
	"net/http"
)

func ConsoleHandler(w http.ResponseWriter, r *http.Request) {

	configFilePath := utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml")
	backupConfig := model.GetBackupConfig(configFilePath)
	consoleName := r.PathValue("console")
	console := backupConfig.FindConsoleConfigByName(consoleName)

	if console == nil {
		http.Error(w, "Console not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/console.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, console); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
