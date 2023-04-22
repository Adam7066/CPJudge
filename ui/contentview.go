package ui

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
	"golang.org/x/image/bmp"
)

type ContentView struct {
	*tview.Pages
}

func NewContentView() *ContentView {
	pages := tview.NewPages()
	pages.SetBorder(true)
	return &ContentView{
		Pages: pages,
	}
}

func (c *ContentView) Load(dir string) {
	if c.HasPage(dir) {
		c.SwitchToPage(dir)
		return
	}
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	switch filepath.Ext(dir) {
	case ".bmp":
		img, err := bmp.Decode(f)
		if err != nil {
			panic(err)
		}
		imageView := tview.NewImage().
			SetColors(tview.TrueColor).
			SetImage(img)
		c.AddAndSwitchToPage(dir, imageView, true)
	default:
		c.LoadReader(f, dir)
	}
}

func (c *ContentView) LoadReader(r io.Reader, name string) {
	textView := tview.NewTextView().
		SetDynamicColors(true)
	w := tview.ANSIWriter(textView)
	io.Copy(w, r)
	c.AddAndSwitchToPage(name, textView, true)
}
