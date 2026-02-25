package ufw

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type UFWStatus struct {
	Active        bool
	DefaultIn     string
	DefaultOut    string
	DefaultRouted string
	Logging       string
	RawVerbose    string
	RawNumbered   string
	Error         error
	UpTime        string
}

var (
	reStatus   = regexp.MustCompile(`(?i)Status:\s*(\w+)`)
	reDefault  = regexp.MustCompile(`(?im)^Default:\s*(.+)$`)
	reLogging  = regexp.MustCompile(`(?im)^Logging:\s*(.+)$`)
	reRuleLine = regexp.MustCompile(`^\s*\[\s*(\d+)\]\s+(.+)$`)
)

func RunCmd(name string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	return out.String(), errOut.String(), err
}

func RunSudo(args ...string) (stdout, stderr string, err error) {
	return RunCmd("sudo", append([]string{"ufw"}, args...)...)
}

func extractPolicy(s string) string {
	parts := strings.Fields(s)
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func GetStatus() (UFWStatus, []Rule) {
	var s UFWStatus

	numOut, _, errNum := RunSudo("status", "numbered")
	verbOut, _, errVerb := RunSudo("status", "verbose")
	s.RawNumbered = numOut
	s.RawVerbose = verbOut
	if errNum != nil && errVerb != nil {
		s.Error = fmt.Errorf("ufw status: %w", errNum)
		return s, nil
	}

	if m := reStatus.FindStringSubmatch(numOut); len(m) > 1 {
		s.Active = strings.EqualFold(m[1], "active")
	}
	if m := reStatus.FindStringSubmatch(verbOut); len(m) > 1 {
		s.Active = strings.EqualFold(m[1], "active")
	}

	if m := reDefault.FindStringSubmatch(verbOut); len(m) > 1 {
		parts := strings.SplitSeq(strings.TrimSpace(m[1]), ",")
		for p := range parts {
			p = strings.TrimSpace(p)
			lower := strings.ToLower(p)
			if strings.Contains(lower, "incoming") {
				s.DefaultIn = strings.ToUpper(extractPolicy(p))
			}
			if strings.Contains(lower, "outgoing") {
				s.DefaultOut = strings.ToUpper(extractPolicy(p))
			}
			if strings.Contains(lower, "routed") {
				s.DefaultRouted = strings.ToUpper(extractPolicy(p))
			}
		}
	}
	if m := reLogging.FindStringSubmatch(verbOut); len(m) > 1 {
		s.Logging = m[1]
	}
	rules := ParseRules(numOut)
	return s, rules
}

func Enable() (stdout, stderr string, err error) {
	return RunSudo("enable")
}

func Disable() (stdout, stderr string, err error) {
	return RunSudo("disable")
}

func DeleteRule(num int) (stdout, stderr string, err error) {
	return RunSudo("delete", fmt.Sprintf("%d", num))
}

func DefaultIncoming(allow bool) (stdout, stderr string, err error) {
	pol := "deny"
	if allow {
		pol = "allow"
	}
	return RunSudo("default", pol, "incoming")
}

func DefaultOutgoing(allow bool) (stdout, stderr string, err error) {
	pol := "deny"
	if allow {
		pol = "allow"
	}
	return RunSudo("default", pol, "outgoing")
}

func AddRule(rule string) (stdout, stderr string, err error) {
	parts := strings.Fields(rule)
	if len(parts) == 0 {
		return "", "", fmt.Errorf("empty rule")
	}
	return RunSudo(append([]string{}, parts...)...)
}
