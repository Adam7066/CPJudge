package main

import (
	"CPJudge/env"
	"CPJudge/extract"
	"CPJudge/judge"
	"CPJudge/myPath"
)

func main() {
	rootPath := myPath.GetRootPath()
	env.InitEnv(rootPath)
	extract.ExtractHomework()
	judge.GenJudgeFile(rootPath)
	judge.AutoRun()
}
