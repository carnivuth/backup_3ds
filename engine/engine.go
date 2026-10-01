package engine

import (
	"console_backupper/model"
	"console_backupper/utils"
	"log"
	"time"
)

// Main backend function, it handles backup scheduling and pruning
func Engine(
	configFilePath string,
	backupNotificationChannel chan *model.ConsoleConfig,
	startBackupsTimer *time.Ticker,
	resetBackupsTimer *time.Ticker,
	pruneNotificationChannel chan *model.ConsoleConfig,
	quitChannel chan bool,
	dataDir string,
	cacheDir string,
	backupsToKeep int) {

		log.Printf("Starting Engine")
		config := utils.ParseConfig(configFilePath)
		go config.StartBackups(backupNotificationChannel)

		for {
			select {
			case <-startBackupsTimer.C:
				go config.StartBackups(backupNotificationChannel)
			case <-resetBackupsTimer.C:
				go config.ResetBackupStates()
			case consoleToBackup := <- backupNotificationChannel:
				go consoleToBackup.BackupConsole(dataDir,cacheDir,pruneNotificationChannel)
			case consoleToPrune := <- pruneNotificationChannel:
				go consoleToPrune.PruneBackups(dataDir,backupsToKeep)
			case <-quitChannel:
				log.Printf("Shutting down Engine")
				return
			}
		}
}

