package look

import (
	"CPJudge/env"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var stuExtractNames []string
var stuOutputNames []string

func getExtract() {
	fileInfo, err := os.ReadDir(env.ExtractPath)
	if err != nil {
		panic(err)
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			stuExtractNames = append(stuExtractNames, file.Name())
		}
	}
	sort.Strings(stuExtractNames)
}

func getOutput() {
	outputPath := filepath.Join(env.ExtractPath, "../output")
	fileInfo, err := os.ReadDir(outputPath)
	if err != nil {
		panic(err)
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			stuOutputNames = append(stuOutputNames, file.Name())
		}
	}
	sort.Strings(stuOutputNames)
}

func LookOutput() {
	getExtract()
	getOutput()
	nowStuPos := 0
	nowOutPos := 0

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	LH := []int{3}
	LV := []int{20}
	for {
		width, height := s.Size()
		DrawBox(s, 0, 0, width-1, height-1, defStyle, fmt.Sprint(nowStuPos, "/", len(stuExtractNames)-1))
		DrawHLine(s, 1, width-2, LH[0], defStyle)
		DrawVLine(s, 1, height-2, LV[0], defStyle)
		s.SetContent(LV[0], LH[0], tcell.RunePlus, nil, defStyle)
		// block right-top
		drawStr := []string{
			fmt.Sprint(stuExtractNames[nowStuPos]),
			"[<-] preStu / [->] nextStu / [^] preFile / [v] nextFile / [c] clear / [q] quit",
		}
		for i := 0; i < len(drawStr); i++ {
			DrawText(s, LV[0]+1, i+1, width-2, i+1, defStyle, drawStr[i])
		}
		// block left-bottom
		stuOutPath := filepath.Join(env.ExtractPath, "../output", stuOutputNames[nowStuPos])
		outFilePaths := []string{}
		drawStr2 := []string{}
		filepath.Walk(stuOutPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && path != stuOutPath {
				drawStr2 = append(drawStr2, strings.Split(path, stuOutPath)[1][1:])
				outFilePaths = append(outFilePaths, path)
			}
			return nil
		})
		for i := 0; i < len(drawStr2); i++ {
			tmpStyle := defStyle
			if i == nowOutPos {
				tmpStyle = defStyle.Background(tcell.ColorYellow)
			}
			DrawText(s, 1, LH[0]+i+1, LV[0]-1, height-2, tmpStyle, drawStr2[i])
		}
		// block right-bottom
		file, err := os.Open(outFilePaths[nowOutPos])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for i := 0; scanner.Scan() && i < 50; i++ {
			line := scanner.Text()
			DrawText(s, LV[0]+1, LH[0]+i+1, width-2, height-2, defStyle, line)
		}
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
				return
			} else if ev.Rune() == 'c' {
				s.Clear()
			} else if ev.Key() == tcell.KeyLeft {
				if nowStuPos > 0 {
					nowStuPos--
					nowOutPos = 0
				}
				s.Clear()
			} else if ev.Key() == tcell.KeyRight {
				if nowStuPos < len(stuExtractNames)-1 {
					nowStuPos++
					nowOutPos = 0
				}
				s.Clear()
			} else if ev.Key() == tcell.KeyUp {
				if nowOutPos > 0 {
					nowOutPos--
				}
				s.Clear()
			} else if ev.Key() == tcell.KeyDown {
				if nowOutPos < len(outFilePaths)-1 {
					nowOutPos++
				}
				s.Clear()
			}
		}
	}
}
