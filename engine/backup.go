package engine

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"console_backupper/compress"
	"console_backupper/ftpdl"
	"console_backupper/model"
)

func BackupConsole(console model.ConsoleConfig, dir string,dataDir string,cacheDir string,pruneNotificationChannel chan model.ConsoleConfig){
	log.Printf("coping %s from %s",dir,console.Name)
	if err := ftpdl.ConnectAndDownloadDir(console.Name,console.Port,console.User,console.Password,dir,filepath.Join(cacheDir,console.Name,dir)); err != nil {
		log.Printf("error in downloading %s from %s error: %v",dir,console.Name,err )
		return
	}

	log.Printf("archiving %s",dir)
	zipFileDir := filepath.Join(dataDir,console.Name,strings.Replace(dir,"/","-",-1))
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
	pruneNotificationChannel <- console
}
