package testcase

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/schaepher/gomics/natsort"
)

type FS struct {
	root string
}

func NewFS(root string) *FS {
	return &FS{root: root}
}

func (f *FS) Open(problem, testcase string) (fs.File, error) {
	return os.Open(filepath.Join(f.root, problem, testcase))
}

func (f *FS) Problems() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(f.root, "*"))
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			ret = append(ret, filepath.Base(match))
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return natsort.Less(ret[i], ret[j])
	})
	return ret, nil
}

func (f *FS) Testcases(problem string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(f.root, problem, "*"))
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			ret = append(ret, filepath.Base(match))
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return natsort.Less(ret[i], ret[j])
	})
	return ret, nil
}
