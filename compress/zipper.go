package compress

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ZipDir creates a zip archive at destZip containing the contents of srcDir.
// Paths inside the archive are relative to srcDir.
func ZipDir(srcDir, destZip string) (err error) {
	info, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", srcDir)
	}

	out, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	zw := zip.NewWriter(out)
	defer func() {
		// Close flushes the central directory; it must succeed for a valid zip.
		if cerr := zw.Close(); err == nil {
			err = cerr
		}
	}()

	// Resolve absolute paths so we can skip the archive itself if it
	// lives inside the directory being zipped.
	absDest, _ := filepath.Abs(destZip)

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if absPath, _ := filepath.Abs(path); absPath == absDest {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil // skip the root itself
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// Zip entries always use forward slashes.
		header.Name = filepath.ToSlash(rel)

		if d.IsDir() {
			header.Name += "/"
			_, err = zw.CreateHeader(header)
			return err
		}

		// Skip anything that isn't a regular file (symlinks, sockets, etc.)
		if !info.Mode().IsRegular() {
			return nil
		}

		header.Method = zip.Deflate // compress the file contents

		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		return err
	})
}


func main() {
	if err := ZipDir("./my_folder", "./my_folder.zip"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("archive created")
}
