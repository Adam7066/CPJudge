package judge

import (
	"CPJudge/env"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Adam7066/golang/log"
	cp "github.com/otiai10/copy"
)

func JudgeStu(stuExtractPath string) error {
	var shareStuDir = filepath.Join(env.SharePath, "stu")
	// Copy student code to share/stu
	err := os.RemoveAll(shareStuDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(shareStuDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = cp.Copy(stuExtractPath, shareStuDir)
	if err != nil {
		return err
	}
	// Run judge
	stu := strings.Split(stuExtractPath, "/")[len(strings.Split(stuExtractPath, "/"))-1]
	stuOutputDir := filepath.Join(env.ExtractPath, "../output", stu)
	os.RemoveAll(stuOutputDir)
	os.MkdirAll(stuOutputDir, os.ModePerm)
	err = os.RemoveAll(env.WorkingPath)
	if err != nil {
		return err
	}
	err = cp.Copy(env.SharePath, env.WorkingPath)
	if err != nil {
		return err
	}
	log.Info.Println("Run judge: " + stuExtractPath)
	outFile, err := os.Create(filepath.Join(stuOutputDir, "out"))
	if err != nil {
		return err
	}
	defer outFile.Close()
	errorFile, err := os.Create(filepath.Join(stuOutputDir, "error"))
	if err != nil {
		return err
	}
	defer errorFile.Close()
	cmd := exec.Command(
		"docker-compose", "run", "--rm",
		"--name", "cpjudge", "homework",
		"/bin/bash", "-c",
		fmt.Sprintf("./autoJudge --limitTime=%d --maxWorkers=%d",
			env.LimitTime,
			env.MaxWorkers,
		),
	)
	cmd.Dir = env.JudgeEnvPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = outFile
	cmd.Stderr = errorFile
	if err = cmd.Run(); err != nil {
		fmt.Println(err)
	}
	// move output to output folder
	cp.Copy(filepath.Join(env.WorkingPath, "output"), stuOutputDir)
	return nil
}

func AutoRun() {
	os.RemoveAll(filepath.Join(env.ExtractPath, "../output"))
	err := filepath.Walk(env.ExtractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != env.ExtractPath {
			if err := JudgeStu(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Error.Println(err)
		return
	}
}
