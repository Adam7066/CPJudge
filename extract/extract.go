package extract

import (
	"CPJudge/myPath"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Adam7066/golang/log"
	"golang.org/x/text/encoding/unicode"
)

func unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	for _, f := range r.File {
		decoder := unicode.UTF8.NewDecoder()
		fname, err := decoder.String(f.Name)
		if err != nil {
			return filenames, err
		}
		fpath := filepath.Join(dest, fname)
		// Check for ZipSlip
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func ExtractHomework(hwZipPath, extractPath string) {
	// Unzip the homework
	if myPath.Exists(extractPath) {
		os.RemoveAll(extractPath)
	}
	_, err := unzip(hwZipPath, extractPath)
	if err != nil {
		log.Error.Println(err)
		return
	}
	// Rename the student folder
	oldPathSlice := []string{}
	newPathSlice := []string{}
	err = filepath.Walk(extractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != extractPath {
			tmpSlice := strings.Split(path, "/")
			newPath := "/"
			for i := 0; i < len(tmpSlice)-1; i++ {
				newPath = filepath.Join(newPath, tmpSlice[i])
			}
			newPath = filepath.Join(newPath, strings.Split(tmpSlice[len(tmpSlice)-1], "_")[0])
			oldPathSlice = append(oldPathSlice, path)
			newPathSlice = append(newPathSlice, newPath)
		}
		return nil
	})
	if err != nil {
		log.Error.Println(err)
		return
	}
	for i := 0; i < len(oldPathSlice); i++ {
		os.Rename(oldPathSlice[i], newPathSlice[i])
	}
	// Unzip student zip
	err = filepath.Walk(extractPath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".zip" {
			tmpSlice := strings.Split(path, "/")
			stuDir := "/"
			for i := 0; i < len(tmpSlice)-1; i++ {
				stuDir = filepath.Join(stuDir, tmpSlice[i])
			}
			_, err := unzip(path, stuDir)
			if err != nil {
				return err
			}
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		log.Error.Println(err)
		return
	}
}
