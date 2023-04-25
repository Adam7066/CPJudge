package ui

import (
	"CPJudge/env"
	"CPJudge/myPath"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func initDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

func Run() {
	app := tview.NewApplication()

	initDir(env.OutputPath)
	explorer := NewExplorer(env.OutputPath)

	contentView := NewContentView()

	showingDiff := false

	loadView := func() {
		node := explorer.GetCurrentNode()
		path := node.GetReference().(string)

		problem, testcase := func() (string, string) {
			parts := strings.Split(path, "/")
			return parts[len(parts)-2], parts[len(parts)-1]
		}()
		switch showingDiff {
		case true:
			ansPath := filepath.Join(env.AnsPath, problem, testcase)
			if myPath.Exists(ansPath) {
				contentView.LoadDiff(ansPath, path, "※※※※※※※※※※")
				break
			}
			fallthrough
		case false:
			contentView.Load(path)
		}
	}

	explorer.SetChangedFunc(func(*tview.TreeNode) {
		loadView()
	})

	explorer.SetSelectedFunc(func(*tview.TreeNode) {
		app.SetFocus(contentView)
	})

	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape || event.Rune() == 'q':
			app.SetFocus(explorer)
		case event.Key() == tcell.KeyEnter:
			node := explorer.GetCurrentNode()
			path := node.GetReference().(string)
			app.Suspend(func() {
				cmd := exec.Command("less", "-S", path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			})

		default:
			return event
		}
		return nil
	})

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(explorer, 0, 1, true)

	main := tview.NewFlex().
		AddItem(sidebar, 24, 1, true).
		AddItem(contentView, 0, 1, false)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(main, 0, 1, true).
		AddItem(NewHint(), 1, 1, false)

	root.SetBorder(true)

	app.SetRoot(root, true).SetFocus(root)

	sidebar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape || event.Rune() == 'q':
			app.Stop()
		case event.Rune() == 'd':
			showingDiff = !showingDiff
			loadView()
		default:
			return event
		}
		return nil
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
