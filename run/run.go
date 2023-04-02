package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
)

func findMakefile(findPath string) (name, path string) {
	retName := ""
	retPath := ""
	filepath.Walk(findPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if retPath != "" {
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.ToLower(info.Name()) == "makefile" || info.Name() == "GNUmakefile") {
			retName = info.Name()
			retPath = path
			return nil
		}
		return nil
	})
	return retName, retPath
}

func runMake(stuFileDirPath string) {
	// copy ta files to stu dir
	cp.Copy("./copy/", stuFileDirPath)
	cmd := exec.Command("make")
	cmd.Dir = stuFileDirPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runJudge(stuFileDirPath string, limitTime int) {
	testcasePath := "./testcase"
	outputPath := "./output"
	filepath.Walk(testcasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(path, "hw") {
			problem := strings.Split(path, "/")[1]
			testcase := strings.Split(path, "/")[2]
			os.MkdirAll(filepath.Join(outputPath, problem), os.ModePerm)
			if _, err := os.Stat(filepath.Join(stuFileDirPath, problem)); os.IsNotExist(err) {
				errorFile, err := os.Create(filepath.Join(outputPath, problem, "err"))
				if err != nil {
					return err
				}
				defer errorFile.Close()
				errorFile.WriteString("Can't find " + problem + " file")
			} else {
				inputFile, err := os.Open(filepath.Join(testcasePath, problem, testcase))
				if err != nil {
					return err
				}
				defer inputFile.Close()
				outputFile, err := os.Create(filepath.Join(outputPath, problem, testcase))
				if err != nil {
					return err
				}
				defer outputFile.Close()
				errorFile, err := os.Create(filepath.Join(outputPath, problem, "err_"+testcase))
				if err != nil {
					return err
				}
				defer errorFile.Close()
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(limitTime)*time.Second)
				defer cancel()
				cmd := exec.CommandContext(
					ctx,
					"./"+filepath.Join(stuFileDirPath, problem),
				)
				cmd.Stdin = inputFile
				cmd.Stdout = outputFile
				cmd.Stderr = errorFile
				err = cmd.Run()
				if err != nil {
					file, err := os.OpenFile(filepath.Join(outputPath, problem, "err"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
					fmt.Fprintln(file, testcase, err)
				}
			}
		}
		return nil
	})
}

func main() {
	makefileName, makefilePath := findMakefile("./stu/")
	stuFileDirPath := strings.Split(makefilePath, "/"+makefileName)[0]
	runMake(stuFileDirPath)
	limitTime := flag.Int("limitTime", 1, "limit time for each testcase")
	flag.Parse()
	runJudge(stuFileDirPath, *limitTime)
}
