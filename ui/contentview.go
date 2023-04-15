package ui

import (
	"CPJudge/selector"
	"bytes"
	"path/filepath"

	"github.com/rivo/tview"
	"golang.org/x/image/bmp"
)

type ContentView struct {
	*tview.Grid
	text     *tview.TextView
	image    *tview.Image
	selector *selector.Selector
}

func NewContentView(selector *selector.Selector) *ContentView {
	text := tview.NewTextView()
	image := tview.NewImage().
		SetColors(tview.TrueColor)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0)

	c := &ContentView{
		Grid:     grid,
		text:     text,
		image:    image,
		selector: selector,
	}
	c.text.SetDynamicColors(true)
	return c
}

func (c *ContentView) Update() {
	switch filepath.Ext(c.selector.CurFileDir()) {
	case ".bmp":
		r := bytes.NewReader(c.selector.CurFileContent())
		img, err := bmp.Decode(r)
		if err != nil {
			panic(err)
		}
		c.image.SetImage(img)
		c.Clear()
		c.Grid.AddItem(c.image, 0, 0, 1, 1, 0, 0, false)
	default:
		c.text.Clear()
		w := tview.ANSIWriter(c.text)
		w.Write(c.selector.CurFileContent())
		c.Clear()
		c.Grid.AddItem(c.text, 0, 0, 1, 1, 0, 0, false)
	}
}
