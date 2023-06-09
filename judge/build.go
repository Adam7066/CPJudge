package judge

import (
	"os"
	"os/exec"
	"runtime"
)

func GenJudgeFile(rootPath string) {
	cmd := exec.Command(
		"go", "build",
		"-o", "autoJudge",
		"./run/run.go",
	)
	cmd.Dir = rootPath
	cmd.Env = append(os.Environ(), "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH="+runtime.GOARCH)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	os.Rename(rootPath+"/autoJudge", rootPath+"/judgeEnv/share/autoJudge")
}
