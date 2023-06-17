package myPath

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func GetRootPath() string {
	absPath := ""
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		absPath = path.Dir(filename)
	}
	return strings.Split(absPath, "/CPJudge")[0] + "/CPJudge/"
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DiskUsage(path string) int64 {
	var size int64

	buf := bytes.NewBuffer(nil)

	cmd := exec.Command("du", "-sk", path)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return -1
	}

	if _, err := fmt.Fscanf(buf, "%d", &size); err != nil {
		return -1
	}
	return size
}
