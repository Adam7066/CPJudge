package main

import (
	"CPJudge/env"
	"CPJudge/extract"
	"CPJudge/judge"
	"CPJudge/myPath"
	"fmt"
	"path"
	"strings"
)

func main() {
	rootPath := myPath.GetRootPath()
	hwZipPath := path.Join(rootPath, env.HWInfo["HWZip"])
	extractPath := path.Join(rootPath, strings.Split(env.HWInfo["HWZip"], ".")[0]+"/extract/")
	// Extract homework
	extract.ExtractHomework(hwZipPath, extractPath)
	// Generate judge file
	judge.GenJudgeFile(rootPath)
	// Auto run judge
	fmt.Print("Please input limit time (s), default=1: ")
	limitTime := 1
	fmt.Scanln(&limitTime)
	judge.AutoRun(extractPath, limitTime)
}
