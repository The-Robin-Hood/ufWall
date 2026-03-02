package ui

import (
	"fmt"
	"strings"
)

func InsertCursor(line string, selected bool, styles Styles) string {
	prefix := "  "
	style := styles.Value
	if selected {
		prefix = "▶ "
		style = styles.ActiveCursorPointer
	}
	line = prefix + style.Render(line)
	return line
}

func InsertCursorRulesSection(line string, selected bool, styles Styles, action string) string {
	prefix := "  "

	if selected {
		return "▶ " + styles.ActiveCursorPointer.Render(line)
	}

	policyStyle := GetPolicyStyle(styles, action)

	actionFormatted := fmt.Sprintf("%-6s", action)
	styledAction := policyStyle.Render(actionFormatted)

	line = strings.Replace(line, actionFormatted, styledAction, 1)

	return prefix + styles.Value.Render(line)
}
