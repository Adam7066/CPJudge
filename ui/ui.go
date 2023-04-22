package ui

import (
	"CPJudge/env"
	"os"
	"os/exec"

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
	initDir(env.AnsPath)
	explorer2 := NewExplorer(env.AnsPath)

	contentView := NewContentView()

	explorer.SetChangedFunc(func(node *tview.TreeNode) {
		path := node.GetReference().(string)
		contentView.Load(path)
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
		case event.Rune() == 'd':
			node1 := explorer.GetCurrentNode()
			path1 := node1.GetReference().(string)
			node2 := explorer2.GetCurrentNode()
			path2 := node2.GetReference().(string)
			contentView.LoadDiff(path2, path1, "※※※※※※※※※※")
		default:
			return event
		}
		return nil
	})

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(explorer, 0, 1, true).
		AddItem(explorer2, 0, 1, true)

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
		case event.Key() == tcell.KeyTab:
			if explorer.HasFocus() {
				app.SetFocus(explorer2)
			} else {
				app.SetFocus(explorer)
			}
		default:
			return event
		}
		return nil
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
