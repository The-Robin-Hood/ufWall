package ufw

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

func ParseRules(output string) []Rule {
	var rules []Rule
	reNumbered := regexp.MustCompile(`^\[\s*(\d+)\]\s+(.+)$`)
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()

		if shouldSkipLine(line) {
			continue
		}

		matches := reNumbered.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		num, _ := strconv.Atoi(matches[1])
		ruleLine := strings.TrimSpace(matches[2])

		rule := parseRuleLine(num, ruleLine)
		rules = append(rules, rule)
	}

	return rules
}

func shouldSkipLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return true
	}

	if strings.Contains(line, "Status:") {
		return true
	}

	if strings.Contains(line, "To") && strings.Contains(line, "Action") && strings.Contains(line, "From") {
		return true
	}

	if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "==") {
		return true
	}

	return false
}

func parseRuleLine(num int, line string) Rule {
	rule := Rule{
		Num: num,
		Raw: line,
	}

	commentIdx := strings.Index(line, "#")
	if commentIdx >= 0 {
		rule.Comment = strings.TrimSpace(line[commentIdx+1:])
		line = strings.TrimSpace(line[:commentIdx])
	}

	parts := splitByMultipleSpaces(line)

	if len(parts) < 3 {
		return rule
	}

	actionIndex := findActionIndex(parts)
	if actionIndex == -1 {
		return rule
	}

	actionParts := strings.Fields(parts[actionIndex])
	if len(actionParts) > 0 {
		rule.Action = actionParts[0]
	}

	if actionIndex > 0 {
		toField := strings.Join(parts[:actionIndex], " ")
		rule.ToDest, rule.ToPort, rule.ToProtocol = parseDestination(toField)
	}

	if actionIndex+1 < len(parts) {
		fromField := strings.Join(parts[actionIndex+1:], " ")
		rule.FromSource, rule.FromPort = parseSource(fromField)
	}

	return rule
}

func findActionIndex(parts []string) int {
	for i, part := range parts {
		upper := strings.ToUpper(strings.Fields(part)[0])
		if upper == "ALLOW" || upper == "DENY" || upper == "REJECT" || upper == "LIMIT" {
			return i
		}
	}
	return -1
}

func splitByMultipleSpaces(s string) []string {
	re := regexp.MustCompile(`\s{2,}`)
	parts := re.Split(s, -1)

	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

func parseDestination(s string) (dest, port, protocol string) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "Anywhere") {
		return s, "any", "any"
	}

	spaceParts := strings.Fields(s)
	if len(spaceParts) >= 2 && isIPOrCIDR(spaceParts[0]) {
		dest = spaceParts[0]
		port, protocol = parsePortSpec(strings.Join(spaceParts[1:], " "))
		return dest, port, protocol
	}

	if isIPOrCIDR(s) {
		return s, "any", "any"
	}

	port, protocol = parsePortSpec(s)
	return "Anywhere", port, protocol
}

func parseSource(s string) (source, port string) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "Anywhere") {
		return s, "any"
	}

	if strings.Contains(s, " ") {
		parts := strings.Fields(s)
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
	}

	return s, "any"
}

func isIPOrCIDR(s string) bool {
	if strings.Contains(s, ".") {
		parts := strings.Split(s, "/")
		if strings.Contains(parts[0], ".") {
			return true
		}
	}

	colonCount := strings.Count(s, ":")
	if colonCount > 1 {
		return true
	}

	return false
}

func parsePortSpec(s string) (port, protocol string) {
	s = strings.TrimSpace(s)

	if strings.Contains(s, "/") {
		parts := strings.SplitN(s, "/", 2)
		if len(parts) == 2 {
			portNum := parts[0]
			protoWithExtra := parts[1]

			protoWords := strings.Fields(protoWithExtra)
			if len(protoWords) > 0 {
				protocol = strings.ToLower(protoWords[0])
				return portNum, protocol
			}
		}
	}

	if strings.Contains(s, ":") && !strings.Contains(s, "::") {
		return s, "any"
	}

	return s, "any"
}
