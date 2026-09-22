// Main file, starts the web interface and the main go routine that manages backups
package main
import (
	"log"
	"net/http"
	"engine"
	"web"
)


func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/:console", backupHandler)
	log.Println("Starting backup engine")
	go BackupEngine()
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
