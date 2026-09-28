package engine

import (
	"console_backupper/utils"
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

func BackupEngine() {
	//parse configuration file
	var backupConfig BackupConfig
	configFilePath := utils.Getenv("CONSOLE_BACKUPPER_CONFIG_FILE", "/etc/console_backupper/config.yml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("error in opening configuration file at %s", configFilePath )
	}

	err = yaml.Unmarshal(configFile, &backupConfig)
	if err != nil {
		log.Fatalf("error in Unmarshaling configuration file at %s", configFilePath )
	}

	log.Printf("Starting backup engine with the following configuration:")
	log.Printf("%+v",backupConfig)
	for _, console := range backupConfig.Consoles {
		for _, dir := range console.Dirs {
			go BackupConsole(console,dir)
		}
	}
}
func BackupConsole(console ConsoleConfig, dir string){
	log.Printf("coping %s from %s",dir,console.Name)
	// TODO: copy content recursively from ftp server
	log.Printf("archiving %s",dir)
	// TODO: archive the content in a tar.gz file
}
