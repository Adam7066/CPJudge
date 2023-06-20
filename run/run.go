package main

import (
	"CPJudge/env"
	"CPJudge/myPath"
	"CPJudge/run/testcase"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mattn/go-shellwords"
	cp "github.com/otiai10/copy"
	"github.com/schaepher/gomics/natsort"
)

var (
	testcasePath = "./testcase"
	outputPath   = "./output"
)

func init() {
	testcasePath, _ = filepath.Abs(testcasePath)
	outputPath, _ = filepath.Abs(outputPath)
}

func removeProblemExec(path string) {
	problems := env.JudgeProblems

	for _, problem := range problems {
		path := filepath.Join(path, problem)
		os.Remove(path)
	}
}

func isMakefile(path string) bool {
	return strings.ToLower(path) == "makefile" || path == "GNUmakefile"
}

func findMakefile(root string) (path string) {
	result := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isMakefile(info.Name()) {
			result = append(result, filepath.Dir(path))
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool {
		len1 := len(strings.Split(result[i], string(os.PathSeparator)))
		len2 := len(strings.Split(result[j], string(os.PathSeparator)))

		if len1 != len2 {
			return len1 < len2
		}
		return natsort.Less(result[i], result[j])
	})

	return result[0]
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

func checkDiskFull(path string, limit int64) chan bool {
	ch := make(chan bool, 1)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if myPath.DiskUsage(path) > limit {
				ch <- true
			}
		}
	}()
	return ch
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
		diskFull := checkDiskFull(dir, env.DiskLimit)
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
		case <-diskFull:
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				fmt.Fprintf(errorFile, "failed to terminate process: %s", err)
				err = cmd.Process.Kill()
				if err != nil {
					fmt.Fprintf(errorFile, "failed to kill process: %s", err)
				}
			}
			return fmt.Errorf("process killed as disk full by size = %dK", myPath.DiskUsage(dir))
		case err := <-done:
			if err != nil {
				return fmt.Errorf("process finished with error = %v", err)
			}
		}
	}
	for _, copyFile := range env.CopyFiles(problem, testcase) {
		paths, err := filepath.Glob(filepath.Join(dir, copyFile))
		if err != nil {
			return fmt.Errorf("failed to find copy file: %s", err)
		}
		for _, path := range paths {
			os.Rename(path, filepath.Join(outputPath, problem, testcase+"_"+filepath.Base(path)))
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
	makefilePath := findMakefile("./stu/")
	stuFileDirPath := makefilePath
	removeProblemExec(stuFileDirPath)
	runMake(stuFileDirPath)
	runJudge(stuFileDirPath)
}
