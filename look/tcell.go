package look

import "github.com/gdamore/tcell/v2"

func DrawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for i, ch := range text {
		s.SetContent(x1+i, y1, ch, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func DrawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
	DrawHLine(s, x1, x2, y1, style)
	DrawHLine(s, x1, x2, y2, style)
	DrawVLine(s, y1, y2, x1, style)
	DrawVLine(s, y1, y2, x2, style)
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	DrawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func DrawHLine(s tcell.Screen, x1, x2, y int, style tcell.Style) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y, tcell.RuneHLine, nil, style)
	}
}

func DrawVLine(s tcell.Screen, y1, y2, x int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	for row := y1; row <= y2; row++ {
		s.SetContent(x, row, tcell.RuneVLine, nil, style)
	}
}
