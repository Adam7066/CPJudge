package judge

import (
	"CPJudge/myPath"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Adam7066/golang/log"
	cp "github.com/otiai10/copy"
)

func AutoRun(extractPath string, limitTime int) {
	os.RemoveAll(filepath.Join(extractPath, "..", "output"))
	err := filepath.Walk(extractPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != extractPath {
			// Create share/stu
			judgeEnvDir := filepath.Join(myPath.GetRootPath(), "judgeEnv")
			workingDir := filepath.Join(judgeEnvDir, "working_copy")
			shareDir := filepath.Join(judgeEnvDir, "share")
			shareStuDir := filepath.Join(shareDir, "stu")
			// Copy student code to share
			err := os.RemoveAll(shareStuDir)
			if err != nil {
				return err
			}
			err = os.MkdirAll(shareStuDir, os.ModePerm)
			if err != nil {
				return err
			}
			err = cp.Copy(path, shareStuDir)
			if err != nil {
				return err
			}
			// Run judge
			stu := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
			stuOutputDir := filepath.Join(extractPath, "..", "output", stu)
			os.MkdirAll(stuOutputDir, os.ModePerm)
			err = os.RemoveAll(workingDir)
			if err != nil {
				return err
			}
			err = cp.Copy(shareDir, workingDir)
			if err != nil {
				return err
			}
			log.Info.Println("Run judge: " + path)
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
				"./autoJudge", "--limitTime="+fmt.Sprint(limitTime),
			)
			cmd.Dir = judgeEnvDir
			cmd.Stdin = os.Stdin
			cmd.Stdout = outFile
			cmd.Stderr = errorFile
			if err = cmd.Run(); err != nil {
				fmt.Println(err)
			}
			// move output to output folder
			cp.Copy(filepath.Join(workingDir, "output"), stuOutputDir)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Error.Println(err)
		return
	}
}
