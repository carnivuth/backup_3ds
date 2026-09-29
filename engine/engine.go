package engine

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"time"
)
// run backup config every minute, return a channel to interrupt the cron
func ConfigScanEngine(configFilePath string, backupNotificationChannel chan ConsoleConfig,backupInterval int) chan bool {
	ticker := time.NewTicker(time.Duration(backupInterval) * time.Minute)
	quit := make(chan bool, 1)
	go func() {
		log.Printf("Starting ConfigScanEngine")
		for {
			select {
			case <-ticker.C:
				ConfigScan(configFilePath,backupNotificationChannel)
			case <-quit:

				ticker.Stop()
			}
		}
	}()
	return quit
}

func ConfigScan(configFilePath string, backupNotificationChannel chan ConsoleConfig){
	//parse configuration file
	var backupConfig BackupConfig
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("error in opening configuration file at %s", configFilePath )
	}
	err = yaml.Unmarshal(configFile, &backupConfig)
	if err != nil {
		log.Fatalf("error in Unmarshaling configuration file at %s", configFilePath )
	}
	log.Printf("Starting console backups with the following configuration:")
	log.Printf("%+v",backupConfig)
	for _, console := range backupConfig.Consoles {
		backupNotificationChannel <- console
	}
}

func BackupEngine(backupNotificationChannel chan ConsoleConfig) {
	log.Printf("Starting backup engine")
	for {
		consoleToBackup := <-backupNotificationChannel
		for _, dir := range consoleToBackup.Dirs {
			go BackupConsole(consoleToBackup, dir)
		}
	}
}

func BackupConsole(console ConsoleConfig, dir string){
	log.Printf("coping %s from %s",dir,console.Name)
	// TODO: copy content recursively from ftp server
	log.Printf("archiving %s",dir)
	// TODO: archive the content in a tar.gz file
}
