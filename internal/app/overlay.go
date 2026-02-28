package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

func PlaceOverlay(x, y int, fg, bg string) string {
	fgLines := strings.Split(fg, "\n")
	bgLines := strings.Split(bg, "\n")
	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, fgLine := range fgLines {
		bgY := y + i
		if bgY < 0 || bgY >= len(bgLines) {
			continue
		}

		fgW := uint(lipgloss.Width(fgLine))
		bgLine := bgLines[bgY]
		bgW := uint(lipgloss.Width(bgLine))

		if uint(x) >= bgW {
			continue
		}

		left := truncate.String(bgLine, uint(x))
		right := ""
		if end := uint(x) + fgW; end < bgW {
			right = rightSlice(bgLine, int(end))
		}

		result[bgY] = left + fgLine + right
	}

	return strings.Join(result, "\n")
}

func rightSlice(s string, fromVisible int) string {
	visible := 0
	inEsc := false
	runes := []rune(s)
	for i, r := range runes {
		if r == '\x1b' {
			inEsc = true
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		if visible == fromVisible {
			return string(runes[i:])
		}
		visible++
	}
	return ""
}

