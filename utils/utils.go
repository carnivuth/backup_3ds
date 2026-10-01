package utils

import (
	"os"
	"log"
	"console_backupper/model"
	"gopkg.in/yaml.v3"
)

func Getenv(key, fallback string) string {
		value := os.Getenv(key)
		if len(value) == 0 {
				return fallback
		}
		return value
}


func ParseConfig(configFilePath string) model.BackupConfig {
	var backupConfig model.BackupConfig
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
	return backupConfig

}
