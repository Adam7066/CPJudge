package env

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Adam7066/golang/log"
	"github.com/joho/godotenv"
)

var LimitTime int
var HWZipPath string
var ExtractPath string
var JudgeEnvPath string
var WorkingPath string
var SharePath string
var OutputPath string
var AnsPath string

var hwInfo = getHWEnv()

func getHWEnv() map[string]string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error.Fatal("Error loading .env file")
	}
	ret := make(map[string]string)
	ret["HWZip"] = os.Getenv("HWZip")
	return ret
}

func InitEnv(rootPath string) {
	fmt.Print("Please input limit time (s), default=1: ")
	LimitTime = 1
	fmt.Scanln(&LimitTime)
	HWZipPath = path.Join(rootPath, hwInfo["HWZip"])
	ExtractPath = path.Join(rootPath, strings.Split(hwInfo["HWZip"], ".")[0]+"/extract/")
	OutputPath = path.Join(rootPath, strings.Split(hwInfo["HWZip"], ".")[0]+"/output/")
	AnsPath = path.Join(rootPath, strings.Split(hwInfo["HWZip"], ".")[0]+"/ans/")
	JudgeEnvPath = filepath.Join(rootPath, "judgeEnv")
	WorkingPath = filepath.Join(JudgeEnvPath, "working_copy")
	SharePath = filepath.Join(JudgeEnvPath, "share")
}
