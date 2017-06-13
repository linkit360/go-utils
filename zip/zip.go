package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// unzip(bytes, size, "/tmp/xxx/")
func Unzip(zipBytes []byte, contentLength int64, target string) error {

	if err := os.MkdirAll(target, 0755); err != nil {
		err = fmt.Errorf("file: %s, os.MkdirAll: %s", target, err.Error())
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(zipBytes), contentLength)
	if err != nil {
		err = fmt.Errorf("zip.NewReader: %s", err.Error())
		return err
	}

	for _, file := range reader.File {
		log.WithFields(log.Fields{
			"file": target + file.Name,
		}).Debug("unzip")

		path := filepath.Join(target, file.Name)

		flagNeedCreateDir := false
		if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
			flagNeedCreateDir = true
		}

		if file.FileInfo().IsDir() || flagNeedCreateDir {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				err = fmt.Errorf("File: %s, unzip path: %s, os.MkdirAll: %s", file.Name, path, err.Error())
				return err
			}
			log.WithFields(log.Fields{
				"dir":      filepath.Dir(path),
				"for file": file.Name,
			}).Debug("create dir")
		}

		fileReader, err := file.Open()
		if err != nil {
			err = fmt.Errorf("name: %s, file.Open: %s", file.Name, err.Error())
			return err
		}
		defer fileReader.Close()
		//log.WithFields(log.Fields{}).Debug("file opened")

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			err = fmt.Errorf("name: %s, file.OpenFile: %s", file.Name, err.Error())
			return err
		}
		defer targetFile.Close()

		//log.WithFields(log.Fields{}).Debug("prepare to copy")
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			err = fmt.Errorf("name: %s, io.Copy: %s", file.Name, err.Error())
			return err
		}
		//log.WithFields(log.Fields{}).Debug("copied")
	}

	return nil
}
