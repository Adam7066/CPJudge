package ui

import "github.com/rivo/tview"

const hint = "[←] prevStu / [→] nextStu / [↑] prevFile / [↓] nextFile / [q] quit / [d] diff / [r] reload / [x] execJudge"

func NewHint() *tview.TextView {
	h := tview.NewTextView().SetText(hint)
	return h
}
