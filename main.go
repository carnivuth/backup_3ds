// Main file, starts the web interface and the main go routine that manages backups
package main
import (
	"log"
	"net/http"
	"console_backupper/engine"
	"console_backupper/web"
)


func main() {
	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/:console", web.BackupHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Starting backup engine")
	go engine.BackupEngine()
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
