package ui

import (
	"container/list"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Explorer struct {
	*tview.TreeView
	list *list.List
	cur  *list.Element
	pos  int
}

func buildFileTree(dir string) *tview.TreeNode {
	m := make(map[string]*tview.TreeNode)
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		node := tview.NewTreeNode(d.Name()).SetReference(path)
		if d.IsDir() {
			node.SetSelectable(false).SetColor(tcell.ColorBlue)
		}
		if path != dir {
			m[filepath.Dir(path)].AddChild(node)
		}
		m[path] = node
		return nil
	})
	return m[dir]
}

func (e *Explorer) init(dir string) {
	list := list.New()
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, entry := range dirEntries {
		path := filepath.Join(dir, entry.Name())
		root := buildFileTree(path)
		list.PushBack(root)
	}

	e.list = list
	e.cur = list.Front()
	e.pos = 0

	e.UpdateRoot()
	e.UpdateTitle()
}

func NewExplorer(dir string) *Explorer {
	treeView := tview.NewTreeView()
	treeView.SetBorder(true)

	e := &Explorer{TreeView: treeView}
	e.init(dir)

	e.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyLeft:
			if e.cur == nil || e.cur.Prev() == nil {
				break
			}
			e.cur = e.cur.Prev()
			e.pos--
			e.UpdateRoot()
			e.UpdateTitle()
		case event.Key() == tcell.KeyRight:
			if e.cur == nil || e.cur.Next() == nil {
				break
			}
			e.cur = e.cur.Next()
			e.pos++
			e.UpdateRoot()
			e.UpdateTitle()
		default:
			return event
		}
		return nil
	})

	return e
}

func (e *Explorer) UpdateTitle() {
	e.SetTitle(fmt.Sprintf("(%d/%d)", e.pos, e.list.Len()))
}

func (e *Explorer) UpdateRoot() {
	if e.cur == nil {
		return
	}
	root := e.cur.Value.(*tview.TreeNode)
	e.SetRoot(root).SetCurrentNode(root)
	// force refresh
	e.TreeView.Move(-1)
}

func (e *Explorer) Reload() {
	e.init(filepath.Dir(e.GetRoot().GetReference().(string)))
}

func (e *Explorer) JumpTo(dir string) {
	for pos, cur := 0, e.list.Front(); cur != nil; pos, cur = pos+1, cur.Next() {
		root := cur.Value.(*tview.TreeNode)
		if !strings.HasPrefix(dir, root.GetReference().(string)) {
			continue
		}
		e.cur = cur
		e.pos = pos
		e.UpdateRoot()
		e.UpdateTitle()
	L:
		for cur := root; ; {
			if cur.GetReference().(string) == dir {
				e.SetCurrentNode(cur)
				break
			}
			for _, child := range cur.GetChildren() {
				if strings.HasPrefix(dir, child.GetReference().(string)) {
					cur = child
					continue L
				}
			}
			break
		}
		break
	}
}
