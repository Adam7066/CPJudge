package myPath

import (
	"os"
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
