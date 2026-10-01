package engine

import (
	"log"
	"time"
	"console_backupper/model"
)

// Main backend function, it handles backup scheduling and pruning
func Engine(
	configFilePath string,
	backupNotificationChannel chan model.ConsoleConfig,
	scanNotificationChannel *time.Ticker,
	pruneNotificationChannel chan model.ConsoleConfig,
	quitChannel chan bool,
	dataDir string,
	cacheDir string,
	backupsToKeep int) {

		log.Printf("Starting Engine")
		go ConfigScan(configFilePath,backupNotificationChannel)
		for {
			select {
			case <-scanNotificationChannel.C:
				go ConfigScan(configFilePath,backupNotificationChannel)
			case consoleToBackup := <- backupNotificationChannel:
				for _, dir := range consoleToBackup.Dirs {
					go BackupConsole(consoleToBackup, dir,dataDir,cacheDir,pruneNotificationChannel)
				}
			case consoleToPrune := <- pruneNotificationChannel:
				for _, dir := range consoleToPrune.Dirs {
				go pruneBackups(consoleToPrune,dir,dataDir,backupsToKeep)
				}
			case <-quitChannel:
				log.Printf("Shutting down Engine")
				return
			}
		}
}

