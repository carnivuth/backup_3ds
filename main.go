// Main file, starts the web interface and the main go routine that manages backups
package console_backupper
import (
	"log"
	"net/http"
	"os"
)

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/:console", backupHandler)
	log.Println("Starting backup engine")
	go backupEngine()
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
