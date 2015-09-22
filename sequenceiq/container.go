package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func debug(s ...string) {
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG]")
		fmt.Fprintln(os.Stderr, s)
	}
}

func extractTarArchiveFile(header *tar.Header, dest string, input io.Reader) error {
	filePath := filepath.Join(dest, header.Name)
	fileInfo := header.FileInfo()

	if fileInfo.IsDir() {
		err := os.MkdirAll(filePath, fileInfo.Mode())
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return os.Symlink(header.Linkname, filePath)
		}

		fileCopy, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		_, err = io.Copy(fileCopy, input)
		if err != nil {
			return err
		}
	}

	return nil
}

func Untar(src, dest string) error {

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	tr := tar.NewReader(r)
	os.MkdirAll(dest, 0755)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		debug("extracting: ", hdr.Name)

		if hdr.Name == "." {
			continue
		}

		err = extractTarArchiveFile(hdr, dest, tr)
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	fmt.Println("Creating Runnable Container ...")

	err := RestoreAssets("/tmp/", "")
	if err != nil {
		fmt.Println(err)
	}

	err = Untar("/tmp/kontainer.tar", "/tmp/kontainer/")
	if err != nil {
		fmt.Println(err)
	}

	os.Chdir("/tmp/kontainer/")

	p := exec.Command("sudo", "./runc", "start")
	debug("execurting process:", strings.Join(p.Args, " "))

	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Stdin = os.Stdin

	if err := p.Start(); err != nil {
		fmt.Println("An error occured: ", err)
	}

	p.Wait()
}
