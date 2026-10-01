// Main file, starts the web interface and the main go routine that manages backups
package main

import (
	"console_backupper/engine"
	"console_backupper/utils"
	"console_backupper/web"
	"log"
	"net/http"
	"time"
	"strconv"
)


func main() {

	backupInterval, err := strconv.Atoi(utils.Getenv("CONSOLE_BACKUPPER_BACKUP_INTERVAL", "1"))
	if err != nil {
		log.Fatalf("Invalid backup interval: %s %v",backupInterval, err)
	}

	backupNotificationChannel := make(chan engine.ConsoleConfig)
	quitChannel := make(chan bool, 1)
	ticker := time.NewTicker(time.Duration(backupInterval) * time.Minute)
	dataDir := utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR","/var/lib/console_backupper")
	cacheDir := utils.Getenv("CONSOLE_BACKUPPER_CACHE_DIR","/var/cache/console_backupper")
	configFilePath := utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml")

	go engine.Engine(configFilePath,backupNotificationChannel,ticker,quitChannel,dataDir,cacheDir)

	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/consoles/{console}", web.ConsoleHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir(utils.Getenv("CONSOLE_BACKUPPER_DATA_DIR", "/var/lib/console_backupper")))))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
