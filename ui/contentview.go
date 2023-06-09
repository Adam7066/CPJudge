package ui

import (
	"CPJudge/ui/uiutil"
	"fmt"
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
		c.SetTitle(dir)
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
	c.SetTitle(dir)
}

func (c *ContentView) LoadString(s, name string) {
	if c.HasPage(name) {
		c.SwitchToPage(name)
		c.SetTitle(name)
		return
	}
	textView := tview.NewTextView().
		SetDynamicColors(true)
	textView.SetText(tview.TranslateANSI(s))
	c.AddAndSwitchToPage(name, textView, true)
	c.SetTitle(name)
}

func (c *ContentView) LoadDiff(srcDir, dstDir, checkpoint string) {
	name := fmt.Sprintf("%s -> %s", srcDir, dstDir)
	if c.HasPage(name) {
		c.SwitchToPage(name)
		c.SetTitle(name)
		return
	}
	src, err := os.ReadFile(srcDir)
	if err != nil {
		panic(err)
	}
	dst, err := os.ReadFile(dstDir)
	if err != nil {
		panic(err)
	}
	diff := uiutil.Diff(string(src), string(dst), checkpoint)
	c.LoadString(diff, name)
	c.SetTitle(name)
}

func (c *ContentView) LoadReader(r io.Reader, name string) {
	if c.HasPage(name) {
		c.SwitchToPage(name)
		c.SetTitle(name)
		return
	}
	textView := tview.NewTextView().
		SetDynamicColors(true)
	w := tview.ANSIWriter(textView)
	io.Copy(w, r)
	c.AddAndSwitchToPage(name, textView, true)
	c.SetTitle(name)
}
