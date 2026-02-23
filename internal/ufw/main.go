package ufw

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type UFWStatus struct {
	Active        bool
	DefaultIn     string
	DefaultOut    string
	DefaultRouted string
	Logging       string
	Rules         []Rule
	RawVerbose    string
	RawNumbered   string
	Error         error
	UpTime        string
}

type Rule struct {
	Num    int
	To     string
	Action string
	From   string
	Raw    string
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

func GetStatus() UFWStatus {
	var s UFWStatus
	// Prefer numbered for rule list; use verbose for defaults/logging.
	numOut, _, errNum := RunSudo("status", "numbered")
	verbOut, _, errVerb := RunSudo("status", "verbose")
	s.RawNumbered = numOut
	s.RawVerbose = verbOut
	if errNum != nil && errVerb != nil {
		s.Error = fmt.Errorf("ufw status: %w", errNum)
		return s
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
				s.DefaultIn = strings.TrimSuffix(p, "(incoming)")
			}
			if strings.Contains(lower, "outgoing") {
				s.DefaultOut = strings.TrimSuffix(p, "(outgoing)")
			}
			if strings.Contains(lower, "routed") {
				s.DefaultRouted = strings.TrimSuffix(p, "(routed)")
			}
		}
	}
	if m := reLogging.FindStringSubmatch(verbOut); len(m) > 1 {
		s.Logging = m[1]
	}

	sc := bufio.NewScanner(strings.NewReader(numOut))
	for sc.Scan() {
		line := sc.Text()
		if subm := reRuleLine.FindStringSubmatch(line); len(subm) >= 3 {
			var num int
			fmt.Sscanf(subm[1], "%d", &num)
			rest := strings.TrimSpace(subm[2])
			fields := splitRuleFields(rest)
			r := Rule{Num: num, Raw: rest}
			if len(fields) >= 3 {
				r.To = fields[0]
				r.Action = fields[1]
				r.From = strings.Join(fields[2:], " ")
			} else {
				r.To = rest
			}
			s.Rules = append(s.Rules, r)
		}
	}
	out, _, _ := RunCmd("systemctl", "show", "ufw", "--property=ActiveEnterTimestamp")
	line := strings.TrimSpace(out)
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 || parts[1] == "" {
		s.Error = fmt.Errorf("could not parse timestamp")
		return s
	}
	timestamp := parts[1]
	startTime, err := time.Parse("Mon 2006-01-02 15:04:05 MST", timestamp)
	if err != nil {
		s.Error = err
		return s
	}

	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	
	s.UpTime = fmt.Sprintf("%dh %dm", hours, minutes)	
	return s
}

func splitRuleFields(s string) []string {
	var out []string
	for f := range strings.FieldsSeq(s) {
		out = append(out, f)
	}
	return out
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

// AddRule runs e.g. "sudo ufw allow 22/tcp".
func AddRule(rule string) (stdout, stderr string, err error) {
	parts := strings.Fields(rule)
	if len(parts) == 0 {
		return "", "", fmt.Errorf("empty rule")
	}
	return RunSudo(append([]string{}, parts...)...)
}
