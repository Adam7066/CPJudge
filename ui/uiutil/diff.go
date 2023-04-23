package uiutil

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(src, dst, checkpoint string) string {
	src = stripansi.Strip(src)
	dst = stripansi.Strip(dst)
	parts1 := strings.Split(src, checkpoint)
	parts2 := strings.Split(dst, checkpoint)

	dmp := diffmatchpatch.New()
	b := &strings.Builder{}
	for i := 0; i < len(parts1) || i < len(parts2); i++ {
		var part1, part2 string
		if i < len(parts1) {
			part1 = parts1[i]
		}
		if i < len(parts2) {
			part2 = parts2[i]
		}
		diffs := dmp.DiffMain(part1, part2, false)
		b.WriteString(dmp.DiffPrettyText(diffs))
	}
	return b.String()
}
