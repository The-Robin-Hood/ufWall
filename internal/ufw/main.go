package ufw

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Stats struct {
	Active  bool
	Logging string
	TotalRules int
}

type Policy struct {
	DefaultIncoming string
	DefaultOutgoing string
	DefaultRouted   string
}

type Rule struct {
	Num    int
	Action string

	ToDest     string
	ToPort     string
	ToProtocol string

	FromSource string
	FromPort   string

	Comment string
	Raw     string
}

type ufwData struct {
	Stats  Stats
	Policy Policy
	Rules  []Rule
	Error  error
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

func GetUFWData() (ufwData) {
	var data ufwData

	numOut, _, errNum := RunSudo("status", "numbered")
	verbOut, _, errVerb := RunSudo("status", "verbose")
	
	if errNum != nil && errVerb != nil {
		data.Error = fmt.Errorf("ufw status: %w", errNum)
		return data
	}

	if m := reStatus.FindStringSubmatch(verbOut); len(m) > 1 {
		data.Stats.Active = strings.EqualFold(m[1], "active")
	}

	if m := reDefault.FindStringSubmatch(verbOut); len(m) > 1 {
		parts := strings.SplitSeq(strings.TrimSpace(m[1]), ",")
		for p := range parts {
			p = strings.TrimSpace(p)
			lower := strings.ToLower(p)
			if strings.Contains(lower, "incoming") {
				data.Policy.DefaultIncoming = strings.ToUpper(extractPolicy(p))
			}
			if strings.Contains(lower, "outgoing") {
				data.Policy.DefaultOutgoing = strings.ToUpper(extractPolicy(p))
			}
			if strings.Contains(lower, "routed") {
				data.Policy.DefaultRouted = strings.ToUpper(extractPolicy(p))
			}
		}
	}
	if m := reLogging.FindStringSubmatch(verbOut); len(m) > 1 {
		data.Stats.Logging = m[1]
	}

	data.Rules = ParseRules(numOut)
	data.Stats.TotalRules = len(data.Rules)

	return data
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

func SetLogging(level string) (stdout, stderr string, err error) {
	return RunSudo("logging", level)
}

func AddRule(rule string) (stdout, stderr string, err error) {
	parts := strings.Fields(rule)
	if len(parts) == 0 {
		return "", "", fmt.Errorf("empty rule")
	}
	return RunSudo(append([]string{}, parts...)...)
}
