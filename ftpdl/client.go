package ftpdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

// DownloadDir recursively copies remoteDir from the FTP server into localDir.
// Symlinks are skipped.
func DownloadDir(c *ftp.ServerConn, remoteDir, localDir string) error {
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		log.Printf("Failed to create local directory %s", localDir)
		return err
	}

	// List fully before doing anything else: FTP allows only one data
	// connection at a time, so we must not Retr while a listing is open.
	entries, err := c.List(remoteDir)
	if err != nil {
		log.Printf("Error listing %s", remoteDir)
		return err
	}

	for _, e := range entries {
		if e.Name == "." || e.Name == ".." {
			continue
		}
		// Guard against a malicious server sending names like "../../x".
		if strings.ContainsAny(e.Name, `/\`) {
			continue
		}

		remotePath := path.Join(remoteDir, e.Name) // remote paths always use "/"
		localPath := filepath.Join(localDir, e.Name)

		switch e.Type {
		case ftp.EntryTypeFolder:
			if err := DownloadDir(c, remotePath, localPath); err != nil {
				return err
			}
		case ftp.EntryTypeFile:
			if err := downloadFile(c, remotePath, localPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadFile(c *ftp.ServerConn, remotePath, localPath string) (err error) {
	resp, err := c.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("retr %s: %w", remotePath, err)
	}
	defer resp.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(f, resp); err != nil {
		return fmt.Errorf("download %s: %w", remotePath, err)
	}
	return nil
}

func ConnectAndDownloadDir(addr string,port int, remoteDir, localDir string) error{

	connectionString := fmt.Sprintf("%s:%d", addr, port)
	c, err := ftp.Dial(connectionString, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
 		log.Printf("Failed to connect to FTP server: %v", err)
		return err
	}
	defer c.Quit()

	if err := DownloadDir(c, remoteDir, localDir); err != nil {
		log.Printf("Failed to download %s from %s",remoteDir,addr)
		return err
	}
	return nil
}
