package main

import (
	"CPJudge/env"
	"CPJudge/myPath"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

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

type Problem struct {
	Name      string
	Testcases []*Testcase
}

type Testcase struct {
	Name string
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

func getProblems() []*Problem {
	problemDirEntries, err := os.ReadDir(testcasePath)
	if err != nil {
		return nil
	}
	problems := make([]*Problem, 0, len(problemDirEntries))
	for _, problemDirEntry := range problemDirEntries {
		testcaseDirEntries, err := os.ReadDir(filepath.Join(testcasePath, problemDirEntry.Name()))
		if err != nil {
			return nil
		}
		testcases := make([]*Testcase, 0, len(testcaseDirEntries))
		for _, testcaseDirEntry := range testcaseDirEntries {
			testcases = append(testcases, &Testcase{Name: testcaseDirEntry.Name()})
		}
		sort.Slice(testcases, func(i, j int) bool {
			return natsort.Less(testcases[i].Name, testcases[j].Name)
		})
		problems = append(problems, &Problem{Name: problemDirEntry.Name(), Testcases: testcases})
	}
	sort.Slice(problems, func(i, j int) bool {
		return natsort.Less(problems[i].Name, problems[j].Name)
	})
	return problems
}

func execJudge(execPath, inputDir, outputDir, errorDir string, limitTime int) error {
	execDir := filepath.Dir(execPath)
	execName := filepath.Base(execPath)
	inputFile, err := os.Open(inputDir)
	if err != nil {
		return err
	}
	defer inputFile.Close()
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
	cmd := exec.Command(
		"./" + execName,
	)
	cmd.Dir = execDir
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
	return nil
}

type testJob struct {
	idx      int
	testcase *Testcase
}

type testResult struct {
	idx int
	err error
}

func runJudge(stuFileDirPath string) {
	problems := getProblems()

	for _, problem := range problems {
		os.MkdirAll(filepath.Join(outputPath, problem.Name), os.ModePerm)
		problemErrorFile, err := os.Create(filepath.Join(outputPath, problem.Name, "err"))
		if err != nil {
			continue
		}
		defer problemErrorFile.Close()
		if !myPath.Exists(filepath.Join(stuFileDirPath, problem.Name)) {
			fmt.Fprintf(problemErrorFile, "Can't find %s file", problem.Name)
			continue
		}

		wg := &sync.WaitGroup{}
		worker := func(jobs <-chan testJob, results chan<- testResult) {
			for job := range jobs {
				testcaseName := job.testcase.Name
				err := execJudge(
					filepath.Join(stuFileDirPath, problem.Name),
					filepath.Join(testcasePath, problem.Name, testcaseName),
					filepath.Join(outputPath, problem.Name, testcaseName),
					filepath.Join(outputPath, problem.Name, "err_"+testcaseName),
					env.LimitTime(problem.Name, testcaseName),
				)
				results <- testResult{job.idx, err}
			}
			wg.Done()
		}

		testJobs := make(chan testJob, len(problem.Testcases))
		testResults := make(chan testResult, len(problem.Testcases))
		for w := 0; w < env.NumWorkers(problem.Name); w++ {
			wg.Add(1)
			go worker(testJobs, testResults)
		}

		for i, testcase := range problem.Testcases {
			testJobs <- testJob{i, testcase}
		}
		close(testJobs)
		wg.Wait()

		execErrors := make([]error, len(problem.Testcases))
		for range execErrors {
			result := <-testResults
			if result.err != nil {
				execErrors[result.idx] = fmt.Errorf("Testcase %s: %s", problem.Testcases[result.idx].Name, result.err)
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
