package ui

func InsertCursor(line string, selected bool, styles Styles) string {
	prefix := "  "
	style := styles.Value
	if selected {
		prefix = "▶ "
		style = styles.ActiveStatus.Bold(true)
	}
	line = prefix + style.Render(line)
	return line
}
