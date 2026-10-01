package engine

import (
	"console_backupper/model"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func pruneBackups(console model.ConsoleConfig,dir string, dataDir string, backupsToKeep int) error{
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
	return nil
}
