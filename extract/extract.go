package extract

import (
	"CPJudge/env"
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

func ExtractHomework() {
	// Unzip the homework
	if myPath.Exists(env.ExtractPath) {
		os.RemoveAll(env.ExtractPath)
	}
	err := os.MkdirAll(env.ExtractPath, os.ModePerm)
	if err != nil {
		log.Error.Println(err)
		return
	}
	_, err = unzip(env.HWZipPath, env.ExtractPath)
	if err != nil {
		log.Error.Println(err)
		return
	}

	// Rename the homework
	dirs, err := filepath.Glob(filepath.Join(env.ExtractPath, "*"))
	if err != nil {
		log.Error.Println(err)
		return
	}
	for _, dir := range dirs {
		underscoreIdx := strings.Index(dir, "_")
		newDir := dir[:underscoreIdx]
		os.Rename(dir, newDir)
	}

	// Unzip the zip files
	zips, err := filepath.Glob(filepath.Join(env.ExtractPath, "*", "*.zip"))
	if err != nil {
		log.Error.Println(err)
		return
	}

	for _, zip := range zips {
		_, err := unzip(zip, filepath.Dir(zip))
		if err != nil {
			log.Error.Println(err)
			return
		}
		os.Remove(zip)
	}
}
