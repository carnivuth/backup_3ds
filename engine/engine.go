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
		go BackupConsole(console)
	}
}
func BackupConsole(console ConsoleConfig ){

}
