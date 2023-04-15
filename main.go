package main

import (
	"CPJudge/env"
	"CPJudge/extract"
	"CPJudge/judge"
	"CPJudge/myPath"
	"CPJudge/ui"
	"fmt"
)

func main() {
	rootPath := myPath.GetRootPath()
	env.InitEnv(rootPath)
	choice := "N"
	fmt.Print("Run Auto Judge (y/N): ")
	fmt.Scanln(&choice)
	if choice == "y" || choice == "Y" {
		extract.ExtractHomework()
		judge.GenJudgeFile(rootPath)
		judge.AutoRun()
	}
	ui.Run()
}
