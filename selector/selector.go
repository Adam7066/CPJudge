package selector

import (
	"io"
	"os"
	"path/filepath"
)

type Selector struct {
	fs      *studentFs
	stuPos  int
	filePos int
	buf     []byte
}

type studentFs struct {
	Root     string
	Students []student
}

type student struct {
	Name  string
	Files []string
}

func NewSelector(name string) (*Selector, error) {
	fs, err := newStudenFS(name)
	if err != nil {
		return nil, err
	}
	return &Selector{
		fs:      fs,
		stuPos:  0,
		filePos: 0,
	}, nil
}

// ====== Pos Methods ======

func (s *Selector) NextStu() {
	if s.stuPos < len(s.fs.Students)-1 {
		s.stuPos++
		s.filePos = 0
		s.buf = nil
	}
}

func (s *Selector) PrevStu() {
	if s.stuPos > 0 {
		s.stuPos--
		s.filePos = 0
		s.buf = nil
	}
}

func (s *Selector) NextFile() {
	if s.filePos < len(s.fs.Students[s.stuPos].Files)-1 {
		s.filePos++
		s.buf = nil
	}
}

func (s *Selector) PrevFile() {
	if s.filePos > 0 {
		s.filePos--
		s.buf = nil
	}
}

func (s *Selector) CurStuName() string {
	return s.fs.Students[s.stuPos].Name
}

func (s *Selector) CurFileName() string {
	return s.fs.Students[s.stuPos].Files[s.filePos]
}

func (s *Selector) CurStuPath() string {
	root := s.fs.Root
	stu := s.CurStuName()
	return filepath.Join(root, stu)
}

func (s *Selector) CurFilePath() string {
	root := s.fs.Root
	stu := s.CurStuName()
	file := s.CurFileName()
	return filepath.Join(root, stu, file)
}

func (s *Selector) CurStuPos() int {
	return s.stuPos
}

func (s *Selector) CurFilePos() int {
	return s.filePos
}

func (s *Selector) StuNum() int {
	return len(s.fs.Students)
}

func (s *Selector) FileNum() int {
	return len(s.fs.Students[s.stuPos].Files)
}

func (s *Selector) StuNames() []string {
	names := make([]string, 0)
	for _, stu := range s.fs.Students {
		names = append(names, stu.Name)
	}
	return names
}

func (s *Selector) FileNames() []string {
	return s.fs.Students[s.stuPos].Files
}

func (s *Selector) Open() (*os.File, error) {
	root := s.fs.Root
	stu := s.CurStuName()
	file := s.CurFileName()
	return os.Open(filepath.Join(root, stu, file))
}

func (s *Selector) CurFileContent() ([]byte, error) {
	if s.buf != nil {
		return s.buf, nil
	}
	f, err := s.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	s.buf = buf
	return buf, nil
}

func newStudenFS(root string) (*studentFs, error) {
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	students := make([]student, 0)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			stu := student{
				Name: dirEntry.Name(),
			}
			filepath.Walk(filepath.Join(root, dirEntry.Name()), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !info.IsDir() {
					rel, _ := filepath.Rel(filepath.Join(root, dirEntry.Name()), path)
					stu.Files = append(stu.Files, rel)
				}
				return nil
			})
			students = append(students, stu)
		}
	}

	return &studentFs{
		Root:     root,
		Students: students,
	}, nil
}
