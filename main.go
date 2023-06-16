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
	fmt.Printf("Problems: %v\n", env.JudgeProblems)
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
