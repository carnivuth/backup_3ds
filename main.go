// Main file, starts the web interface and the main go routine that manages backups
package main

import (
	"console_backupper/engine"
	"console_backupper/utils"
	"console_backupper/web"
	"log"
	"net/http"
)


func main() {
	consoleBackupNotificationChannel := make(chan engine.ConsoleConfig)

	go engine.BackupEngine(consoleBackupNotificationChannel,utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR","/var/lib/console_backupper"),utils.Getenv("CONSOLE_BACKUPPER_CACHE_DIR","/var/cache/console_backupper"))

	engine.ConfigScanEngine(utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml"), consoleBackupNotificationChannel,1)

	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/consoles/{console}", web.ConsoleHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir(utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")))))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
