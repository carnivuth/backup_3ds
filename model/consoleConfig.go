package model

import (
	"sync"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"sort"
	"fmt"

	"console_backupper/compress"
	"console_backupper/ftpdl"
)

type ConsoleConfig struct {
	Name string `yaml:"name"`
	User string `yaml:"user"`
	Password string `yaml:"password" `
	Port int `yaml:"port"`
	Dirs []string `yaml:"dirs"`
	RunningBackup sync.Mutex `yaml:"-"`
}
func (console *ConsoleConfig) BackupConsole(dataDir string,cacheDir string,pruneNotificationChannel chan *ConsoleConfig){

	console.RunningBackup.Lock()
	defer console.RunningBackup.Unlock()
	for _, dir := range console.Dirs {
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
}

func (console *ConsoleConfig) PruneBackups(dataDir string, backupsToKeep int) error{
	for _, dir := range console.Dirs {
		log.Printf("Starting backup pruning")
		dirToPrune := filepath.Join(dataDir,console.Name,strings.Replace(dir,"/","-",-1))
		if backupsToKeep < 0 {
			return fmt.Errorf("backupsToKeep must be >= 0 got %d")
		}

		entries, err := os.ReadDir(dirToPrune)
		if err != nil {
			log.Printf("error reading dir %s", dirToPrune)
			return err
		}

		var files []os.FileInfo
		for _, e := range entries {

			if !e.Type().IsRegular() {
				os.Remove(e.Name())
			}else{

				info, err := e.Info()
				if err != nil {
					log.Printf("error reading %s", e.Name())
					return err
				}
				files = append(files, info)
			}
		}

		if len(files) <= backupsToKeep {
			return nil
		}

		// Newest first
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().UnixNano() > files[j].ModTime().UnixNano()
		})

		for _, f := range files[backupsToKeep:] {
			path := filepath.Join(dirToPrune, f.Name())
			log.Printf("pruning %s: ", path)
			if err := os.Remove(path); err != nil {
				log.Printf("error removing %s: ", path)
				return err
			}
		}
	}
	return nil
}
