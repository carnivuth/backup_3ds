package engine

import (
	"console_backupper/ftpdl"
	"log"
	"os"
	"path/filepath"
	"time"
	"strings"
	"gopkg.in/yaml.v3"
	"console_backupper/compress"
)

// run backup config every minute, return a channel to interrupt the cron
func Engine(configFilePath string, backupNotificationChannel chan ConsoleConfig, scanNotificationChannel *time.Ticker, quitChannel chan bool, dataDir string, cacheDir string ){
		log.Printf("Starting Engine")
		go ConfigScan(configFilePath,backupNotificationChannel)
		log.Printf("asd")
		for {
			select {
			case <-scanNotificationChannel.C:
				go ConfigScan(configFilePath,backupNotificationChannel)
			case consoleToBackup := <- backupNotificationChannel:
				for _, dir := range consoleToBackup.Dirs {
					go BackupConsole(consoleToBackup, dir,dataDir,cacheDir)
				}
			case <-quitChannel:
				log.Printf("Shutting down Engine")
				return
			}
		}
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
	log.Printf("parsed configuration:")
	log.Printf("%+v",backupConfig)
	for _, console := range backupConfig.Consoles {
		log.Printf("super asd")
		backupNotificationChannel <- console
	}
}

func BackupConsole(console ConsoleConfig, dir string,dataDir string,cacheDir string){
	log.Printf("coping %s from %s",dir,console.Name)
	if err := ftpdl.ConnectAndDownloadDir(console.Name,console.Port,console.User,console.Password,dir,filepath.Join(cacheDir,console.Name,dir)); err != nil {
		log.Printf("error in downloading %s from %s error: %v",dir,console.Name,err )
		return
	}

	log.Printf("archiving %s",dir)
	zipFileDir := filepath.Join(dataDir,console.Name)
	zipFileName := console.Name + "-" + strings.Replace(dir,"/","-",-1) + "-" + time.Now().Format("2006-01-02-15-04-05") + ".zip"

	if err := os.MkdirAll(zipFileDir, 0o755); err != nil {
		log.Printf("Failed to create local directory %s error: %v", zipFileDir, err)
		return

	}

	if err := compress.ZipDir(filepath.Join(cacheDir,console.Name,dir), filepath.Join(zipFileDir,zipFileName)); err != nil {
		log.Printf("error in archiving %s, Error: %v", dir, err)
		return
	}

	log.Printf("archive %s created",zipFileName)

}
