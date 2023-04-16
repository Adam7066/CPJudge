package ui

import (
	"CPJudge/env"
	"os"
	"os/exec"

	"github.com/rivo/tview"
)

var (
	app      *tview.Application
	explorer *Explorer
)

func Run() {
	contentView := NewContentView()

	explorer = NewExplorer(env.OutputPath)
	explorer.SetChangedFunc(func(node *tview.TreeNode) {
		path := node.GetReference().(string)
		contentView.Load(path)
	})
	explorer.SetSelectedFunc(func(node *tview.TreeNode) {
		path := node.GetReference().(string)
		app.Suspend(func() {
			cmd := exec.Command("less", "-S", path)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		})
	})

	main := tview.NewFlex().
		AddItem(explorer, 24, 1, true).
		AddItem(contentView, 0, 1, false)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(main, 0, 1, true).
		AddItem(NewHint(), 1, 1, false)

	root.SetBorder(true)

	app = tview.NewApplication().SetRoot(root, true).SetFocus(root)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
