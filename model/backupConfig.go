package model

import (
	"log"
)

type BackupConfig struct {
	Consoles []ConsoleConfig
}

func (config *BackupConfig) FindConsoleConfigByName(name string) *ConsoleConfig {
	for i, _ := range config.Consoles {
		if config.Consoles[i].Name == name {
			return &config.Consoles[i]
		}
	}
	return nil
}

func (config *BackupConfig) ResetBackupStates(){
	log.Printf("Resetting backup states for all consoles")
	for i,_ := range config.Consoles {
		config.Consoles[i].DoneBackup = false
	}
}

func (config *BackupConfig) StartBackups(backupNotificationChannel chan *ConsoleConfig){
	log.Printf("Starting backups for all consoles")
	for i,_ := range config.Consoles {
		backupNotificationChannel <- &config.Consoles[i]
	}
}

