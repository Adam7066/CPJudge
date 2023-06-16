package main

import (
	"CPJudge/env"
	"CPJudge/myPath"
	"CPJudge/run/testcase"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mattn/go-shellwords"
	cp "github.com/otiai10/copy"
)

var (
	testcasePath = "./testcase"
	outputPath   = "./output"
)

func init() {
	testcasePath, _ = filepath.Abs(testcasePath)
	outputPath, _ = filepath.Abs(outputPath)
}

func findMakefile(findPath string) (name, path string) {
	retName := ""
	retPath := ""
	filepath.Walk(findPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if retPath != "" {
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.ToLower(info.Name()) == "makefile" || info.Name() == "GNUmakefile") {
			retName = info.Name()
			retPath = path
			return nil
		}
		return nil
	})
	return retName, retPath
}

func runMake(stuFileDirPath string) {
	// copy ta files to stu dir
	cp.Copy("./copy/", stuFileDirPath)
	cmd := exec.Command("make")
	cmd.Dir = stuFileDirPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func execJudge(fs *testcase.FS, dir, problem, testcase string) error {
	// input
	inputFile, err := fs.Open(problem, testcase)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// output and error
	outputDir := filepath.Join(outputPath, problem, testcase)
	errorDir := filepath.Join(outputPath, problem, "err_"+testcase)
	outputFile, err := os.Create(outputDir)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	errorFile, err := os.Create(errorDir)
	if err != nil {
		return err
	}
	defer errorFile.Close()

	limitTime := env.LimitTime(problem, testcase)
	commands := env.ExecCommands(problem, testcase)

	for _, command := range commands {
		args, err := shellwords.Parse(command)
		if err != nil {
			return err
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdin = inputFile
		cmd.Stdout = outputFile
		cmd.Stderr = errorFile
		err = cmd.Start()
		if err != nil {
			return err
		}
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(time.Duration(limitTime) * time.Second):
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				fmt.Fprintf(errorFile, "failed to terminate process: %s", err)
				err = cmd.Process.Kill()
				if err != nil {
					fmt.Fprintf(errorFile, "failed to kill process: %s", err)
				}
			}
			return fmt.Errorf("process killed as timeout reached")
		case err := <-done:
			if err != nil {
				return fmt.Errorf("process finished with error = %v", err)
			}
		}
	}
	return nil
}

type testJob struct {
	idx      int
	problem  string
	testcase string
}

type testResult struct {
	idx int
	err error
}

func runJudge(stuFileDirPath string) {
	fs := testcase.NewFS(testcasePath)
	problems := env.JudgeProblems

	for _, problem := range problems {
		os.MkdirAll(filepath.Join(outputPath, problem), os.ModePerm)
		problemErrorFile, err := os.Create(filepath.Join(outputPath, problem, "err"))
		if err != nil {
			continue
		}
		defer problemErrorFile.Close()
		if !myPath.Exists(filepath.Join(stuFileDirPath, problem)) {
			fmt.Fprintf(problemErrorFile, "can't find %s file", problem)
			continue
		}

		wg := &sync.WaitGroup{}
		worker := func(jobs <-chan testJob, results chan<- testResult) {
			for job := range jobs {
				err := execJudge(
					fs,
					stuFileDirPath,
					job.problem,
					job.testcase,
				)
				results <- testResult{job.idx, err}
			}
			wg.Done()
		}

		testcases, err := fs.Testcases(problem)
		if err != nil {
			fmt.Fprintf(problemErrorFile, "can't find %s file", problem)
			continue
		}

		testJobs := make(chan testJob, len(testcases))
		testResults := make(chan testResult, len(testcases))
		for w := 0; w < env.NumWorkers(problem); w++ {
			wg.Add(1)
			go worker(testJobs, testResults)
		}

		for i, testcase := range testcases {
			testJobs <- testJob{i, problem, testcase}
		}
		close(testJobs)
		wg.Wait()

		execErrors := make([]error, len(testcases))
		for range execErrors {
			result := <-testResults
			if result.err != nil {
				execErrors[result.idx] = fmt.Errorf("testcase %s: %s", testcases[result.idx], result.err)
			}
		}
		for _, err := range execErrors {
			if err != nil {
				fmt.Fprintln(problemErrorFile, err)
			}
		}
	}
}

func main() {
	makefileName, makefilePath := findMakefile("./stu/")
	stuFileDirPath := strings.Split(makefilePath, "/"+makefileName)[0]
	runMake(stuFileDirPath)
	runJudge(stuFileDirPath)
}
