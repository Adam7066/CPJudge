package ui

import (
	"CPJudge/env"
	"CPJudge/selector"
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Run() {
	selector, err := selector.NewSelector(env.OutputPath)
	if err != nil {
		panic(err)
	}

	sidebar := NewSidebar(selector)
	sidebar.Update()
	contentView := NewContentView(selector)

	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(22, 0).
		SetBorders(true).
		AddItem(NewHint(), 1, 0, 1, 2, 0, 0, false)

	grid.AddItem(sidebar, 0, 0, 1, 1, 0, 0, false).
		AddItem(contentView, 0, 1, 1, 1, 0, 0, false)

	update := func() {
		sidebar.Update()
		contentView.Update()
	}

	app := tview.NewApplication().SetRoot(grid, true).SetFocus(grid)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Rune() == 'q':
			app.Stop()
		case event.Key() == tcell.KeyLeft:
			selector.PrevStu()
			update()
		case event.Key() == tcell.KeyRight:
			selector.NextStu()
			update()
		case event.Key() == tcell.KeyUp:
			selector.PrevFile()
			update()
		case event.Key() == tcell.KeyDown:
			selector.NextFile()
			update()
		case event.Key() == tcell.KeyEnter:
			app.Suspend(func() {
				fmt.Println(selector.CurFileDir())
				cmd := exec.Command("less", "-S", selector.CurFileDir())
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			})
		}

		return event
	}).
		SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
			switch {
			case event.Buttons() == tcell.WheelUp:
				selector.PrevStu()
				update()
			case event.Buttons() == tcell.WheelDown:
				selector.NextStu()
				update()
			}
			return event, action
		})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
