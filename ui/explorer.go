package ui

import (
	"container/list"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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

func NewExplorer(dir string) *Explorer {
	tree := tview.NewTreeView()
	tree.SetBorder(true)
	list := list.New()
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, entry := range dirEntries {
		path := filepath.Join(dir, entry.Name())
		tree := buildFileTree(path)
		list.PushBack(tree)
	}
	root := list.Front().Value.(*tview.TreeNode)
	tree.SetRoot(root).SetCurrentNode(root)

	explorer := &Explorer{
		TreeView: tree,
		list:     list,
		cur:      list.Front(),
	}
	explorer.UpdateTitle()

	explorer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyLeft && explorer.cur.Prev() != nil:
			explorer.cur = explorer.cur.Prev()
			explorer.pos--
			root := explorer.cur.Value.(*tview.TreeNode)
			explorer.SetRoot(root).SetCurrentNode(root)
			explorer.UpdateTitle()
		case event.Key() == tcell.KeyRight && explorer.cur.Next() != nil:
			explorer.cur = explorer.cur.Next()
			explorer.pos++
			root := explorer.cur.Value.(*tview.TreeNode)
			explorer.SetRoot(root).SetCurrentNode(root)
			explorer.UpdateTitle()
		default:
			return event
		}
		return nil
	})
	return explorer
}

func (e *Explorer) UpdateTitle() {
	e.SetTitle(fmt.Sprintf("(%d/%d)", e.pos, e.list.Len()))
}
