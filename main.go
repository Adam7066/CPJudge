package main

import (
	"CPJudge/extract"
	"CPJudge/judge"
	"CPJudge/myPath"
	"CPJudge/ui"
	"fmt"
)

func main() {
	rootPath := myPath.GetRootPath()
	judge.GenJudgeFile(rootPath)
	choice := "N"
	fmt.Print("Run Auto Judge (y/N): ")
	fmt.Scanln(&choice)
	if choice == "y" || choice == "Y" {
		extract.ExtractHomework()
		judge.AutoRun()
	}
	ui.Run()
}
