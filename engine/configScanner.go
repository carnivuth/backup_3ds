package engine

import (
	"console_backupper/model"
	"console_backupper/utils"
)

func ConfigScan(configFilePath string, backupNotificationChannel chan *model.ConsoleConfig){

	backupConfig := utils.ParseConfig(configFilePath)
	for i,_ := range backupConfig.Consoles {
		backupNotificationChannel <- &backupConfig.Consoles[i]
	}
}


