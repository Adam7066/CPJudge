package env

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Adam7066/golang/log"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/viper"
)

//go:embed config.toml
var fs embed.FS

var (
	HW            string
	DockerCmd     string
	DockerArgs    []string
	JudgeFileName string
	HWZipPath     string
	ExtractPath   string
	JudgeEnvPath  string
	WorkingPath   string
	SharePath     string
	OutputPath    string
	AnsPath       string
)

func init() {
	f, err := fs.Open("config.toml")
	if err != nil {
		log.Error.Fatalln("Cannot open config.toml")
	}
	defer f.Close()
	viper.SetConfigType("toml")
	err = viper.ReadConfig(f)
	if err != nil {
		log.Error.Fatalln("Cannot read config.toml")
	}
	viper.SetDefault("HW", "*")
	viper.SetDefault("docker-compose.cmd", "docker-compose")
	viper.SetDefault("judge.filename", "./autoJudge")
	viper.SetDefault("judge.timeLimit", 1)
	viper.SetDefault("judge.numWorkers", 1)
	viper.SetDefault("judge.cmds", []string{"./{name}"})
	viper.SetDefault("judge.copyFiles", []string{})

	HW = viper.GetString("HW")
	HWZipPath = filepath.Join(HW + ".zip")
	dockerCommand := viper.GetString("docker-compose.cmd")
	JudgeFileName = viper.GetString("judge.filename")
	DockerCmd, DockerArgs = parseDockerCommand(dockerCommand, JudgeFileName)
	ExtractPath = filepath.Join(HW, "extract")
	OutputPath = filepath.Join(HW, "output")
	AnsPath = filepath.Join(HW, "ans")
	JudgeEnvPath = "judgeEnv"
	WorkingPath = filepath.Join(JudgeEnvPath, "working_copy")
	SharePath = filepath.Join(JudgeEnvPath, "share")
}

func replaceCommand(command, problem, testcase string) string {
	command = strings.ReplaceAll(command, "{name}", problem)
	command = strings.ReplaceAll(command, "{case}", testcase)
	return command
}

func replaceCommands(commands []string, problem, testcase string) []string {
	ret := make([]string, 0, len(commands))
	for _, command := range commands {
		ret = append(ret, replaceCommand(command, problem, testcase))
	}
	return ret
}

func parseDockerCommand(dockerCommand, judgeFileName string) (name string, args []string) {
	cmd := fmt.Sprintf("%s run --rm --name cpjudge homework /bin/bash -c %s", dockerCommand, judgeFileName)
	ret, err := shellwords.Parse(cmd)
	if err != nil {
		log.Error.Fatalln("Cannot parse docker command")
	}
	return ret[0], ret[1:]
}

func NumWorkers(problem string) int {
	if viper.IsSet("judge." + problem + ".numWorkers") {
		return viper.GetInt("judge." + problem + ".numWorkers")
	}
	return viper.GetInt("judge.numWorkers")
}

func LimitTime(problem, testcase string) int {
	if viper.IsSet("judge." + problem + "." + testcase + ".timeLimit") {
		return viper.GetInt("judge." + problem + "." + testcase + ".timeLimit")
	}
	if viper.IsSet("judge." + problem + ".timeLimit") {
		return viper.GetInt("judge." + problem + ".timeLimit")
	}
	return viper.GetInt("judge.timeLimit")
}

func ExecCommands(problem, testcase string) []string {
	if viper.IsSet("judge." + problem + "." + testcase + ".cmds") {
		return replaceCommands(viper.GetStringSlice("judge."+problem+"."+testcase+".cmds"), problem, testcase)
	}
	if viper.IsSet("judge." + problem + ".cmds") {
		return replaceCommands(viper.GetStringSlice("judge."+problem+".cmds"), problem, testcase)
	}
	return replaceCommands(viper.GetStringSlice("judge.cmds"), problem, testcase)
}

func CopyFiles(problem, testcase string) []string {
	if viper.IsSet("judge." + problem + "." + testcase + ".copyFiles") {
		return viper.GetStringSlice("judge." + problem + "." + testcase + ".copyFiles")
	}
	if viper.IsSet("judge." + problem + ".copyFiles") {
		return viper.GetStringSlice("judge." + problem + ".copyFiles")
	}
	return viper.GetStringSlice("judge.copyFiles")
}
