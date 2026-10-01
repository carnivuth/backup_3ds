package model

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

type backupConfig struct {
	Consoles []ConsoleConfig
}

var config *backupConfig = nil

func GetBackupConfig(configFilePath string) * backupConfig {
	if config == nil {
		configFile, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatalf("error in opening configuration file at %s", configFilePath )
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			log.Fatalf("error in Unmarshaling configuration file at %s", configFilePath )
		}
		log.Printf("parsed configuration:")
		log.Printf("%+v",config)
	}
		return config

}
func (config *backupConfig) FindConsoleConfigByName(name string) *ConsoleConfig {
	for i, _ := range config.Consoles {
		if config.Consoles[i].Name == name {
			return &config.Consoles[i]
		}
	}
	return nil
}

func (config *backupConfig) ResetBackupStates(){
	log.Printf("Resetting backup states for all consoles")
	for i,_ := range config.Consoles {
		config.Consoles[i].DoneBackup = false
	}
}

func (config *backupConfig) StartBackups(backupNotificationChannel chan *ConsoleConfig){
	log.Printf("Starting backups for all consoles")
	for i,_ := range config.Consoles {
		backupNotificationChannel <- &config.Consoles[i]
	}
}

