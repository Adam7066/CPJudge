package ui

import (
	"CPJudge/selector"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	*tview.Grid
	selector *selector.Selector
	percent  *tview.TextView
	fileTree *tview.TreeView
}

func NewSidebar(selector *selector.Selector) *Sidebar {
	// name := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	percent := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	divider := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(strings.Repeat("─", 30))
	fileTree := tview.NewTreeView()

	grid := tview.NewGrid().
		SetRows(1, 1, 0).
		SetColumns(1, 0, 1).
		// display "← (0/100) →"
		AddItem(tview.NewTextView().SetText("←"), 0, 0, 1, 1, 0, 0, false).
		AddItem(percent, 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetText("→"), 0, 2, 1, 1, 0, 0, false).
		// display "──────────────────────────────"
		AddItem(divider, 1, 0, 1, 3, 0, 0, false).
		AddItem(fileTree, 2, 0, 1, 3, 0, 0, false)

	s := &Sidebar{
		Grid:     grid,
		selector: selector,
		percent:  percent,
		fileTree: fileTree,
	}
	s.Update()
	return s
}

func (s *Sidebar) buildTree() {
	var selected *tview.TreeNode

	var build func(string) *tview.TreeNode
	build = func(dir string) *tview.TreeNode {
		name := filepath.Base(dir)
		node := tview.NewTreeNode(name).SetReference(dir)

		dirEntries, err := os.ReadDir(dir)
		// dir is a file, not a directory
		if err != nil {
			if dir == s.selector.CurFileDir() {
				selected = node
			}
			return node
		}

		node.SetSelectable(false).SetColor(tcell.ColorBlue)
		for _, dirEntry := range dirEntries {
			node.AddChild(build(filepath.Join(dir, dirEntry.Name())))
		}
		return node
	}

	root := build(s.selector.CurStuDir())
	s.fileTree.SetRoot(root).SetCurrentNode(selected)
}

func (s *Sidebar) Update() {
	percent := fmt.Sprintf("(%d/%d)", s.selector.CurStuPos(), s.selector.StuNum())
	s.percent.SetText(percent)
	s.buildTree()
}
