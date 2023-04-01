package main

import (
	"CPJudge/env"
	"CPJudge/extract"
	"CPJudge/judge"
	"CPJudge/myPath"
	"path"
	"strings"
)

func main() {
	rootPath := myPath.GetRootPath()
	hwZipPath := path.Join(rootPath, env.HWInfo["HWZip"])
	extractPath := path.Join(rootPath, strings.Split(env.HWInfo["HWZip"], ".")[0]+"/extract/")
	// Extract homework
	extract.ExtractHomework(hwZipPath, extractPath)
	// Auto run judge
	judge.AutoRun(extractPath)
}
