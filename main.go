package main

import (
	"CPJudge/env"
	"CPJudge/extract"
	"CPJudge/myPath"
	"path"
	"strings"
)

func main() {
	rootPath := myPath.GetRootPath()
	hwZipPath := path.Join(rootPath, env.HWInfo["HWZip"])
	extractPath := path.Join(rootPath, strings.Split(env.HWInfo["HWZip"], ".")[0]+"/extract/")
	extract.ExtractHomework(hwZipPath, extractPath)
}
