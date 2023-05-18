package env

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var LimitTime int = 1
var MaxWorkers int = 1
var HWZipPath string
var ExtractPath string
var JudgeEnvPath string
var WorkingPath string
var SharePath string
var OutputPath string
var AnsPath string

var HWZip string
var CopyFiles []string

func getHWEnv() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	HWZip = viper.GetString("HWZip")
	CopyFiles = viper.GetStringSlice("CopyFile")
}

func InitEnv(rootPath string) {
	fmt.Print("Please input limit time (s), default=1: ")
	fmt.Scanln(&LimitTime)
	fmt.Print("Please input max worker (s), default=1: ")
	fmt.Scanln(&MaxWorkers)
	getHWEnv()
	HWZipPath = path.Join(rootPath, HWZip)
	ExtractPath = path.Join(rootPath, strings.Split(HWZip, ".")[0]+"/extract/")
	OutputPath = path.Join(rootPath, strings.Split(HWZip, ".")[0]+"/output/")
	AnsPath = path.Join(rootPath, strings.Split(HWZip, ".")[0]+"/ans/")
	JudgeEnvPath = filepath.Join(rootPath, "judgeEnv")
	WorkingPath = filepath.Join(JudgeEnvPath, "working_copy")
	SharePath = filepath.Join(JudgeEnvPath, "share")
}
