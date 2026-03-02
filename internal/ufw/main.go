package ufw

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

type Stats struct {
	Active     bool
	Logging    string
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
	IPv6   bool 

	ToDest     string
	ToPort     string
	ToProtocol string

	FromSource string
	FromPort   string

	Comment string
	Raw     string
}

type ufwData struct {
	Stats     Stats
	Policy    Policy
	Rules     []Rule 
	IPv4Rules []Rule
	IPv6Rules []Rule
	Error     error
}

// UFW action constants
const (
	ActionAllow  = "ALLOW"
	ActionDeny   = "DENY"
	ActionReject = "REJECT"
	ActionLimit  = "LIMIT"
)

var Actions = []string{ActionAllow, ActionDeny, ActionReject, ActionLimit}

var (
	reStatus  = regexp.MustCompile(`(?i)Status:\s*(\w+)`)
	reDefault = regexp.MustCompile(`(?im)^Default:\s*(.+)$`)
	reLogging = regexp.MustCompile(`(?im)^Logging:\s*(.+)$`)
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

func GetUFWData() ufwData {
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

	for _, rule := range data.Rules {
		if rule.IPv6 {
			data.IPv6Rules = append(data.IPv6Rules, rule)
		} else {
			data.IPv4Rules = append(data.IPv4Rules, rule)
		}
	}

	return data
}

func Enable() (stdout, stderr string, err error) {
	return RunSudo("enable")
}

func Disable() (stdout, stderr string, err error) {
	return RunSudo("disable")
}

func DeleteRule(num int) (stdout, stderr string, err error) {
	return RunSudo("--force", "delete", fmt.Sprintf("%d", num))
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

func DefaultRouted(allow bool) (stdout, stderr string, err error) {
	pol := "deny"
	if allow {
		pol = "allow"
	}
	return RunSudo("default", pol, "routed")
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

func InsertRule(position int, args ...string) (stdout, stderr string, err error) {
	cmdArgs := append([]string{"insert", fmt.Sprintf("%d", position)}, args...)
	log.Printf("InsertRule cmdArgs: %v", cmdArgs)
	return RunSudo(cmdArgs...)
}

// normalizeAddress cleans up UFW display addresses for use in commands
// "Anywhere (v6)" -> "" (let UFW handle it)
// "Anywhere" -> "" (let UFW handle it)
// "192.168.1.0/24" -> "192.168.1.0/24"
func normalizeAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	addr = strings.TrimSuffix(addr, " (v6)")
	if addr == "Anywhere" || addr == "" {
		return ""
	}
	return addr
}

// normalizePort cleans up UFW display port values for use in commands
// "80 (v6)" -> "80"
// "22/tcp (v6)" -> "22/tcp"
func normalizePort(port string) string {
	port = strings.TrimSpace(port)
	port = strings.TrimSuffix(port, " (v6)")
	return port
}

func InsertRuleFromExisting(position int, action string, rule Rule) (stdout, stderr string, err error) {
	args := buildRuleArgs(action, rule)

	log.Printf("InsertRuleFromExisting: position=%d, action=%s, ipv6=%v, args=%v", position, action, rule.IPv6, args)

	if position <= 0 {
		return AppendRule(args...)
	}
	return InsertRule(position, args...)
}

func AppendRule(args ...string) (stdout, stderr string, err error) {
	log.Printf("AppendRule args: %v", args)
	return RunSudo(args...)
}

func buildRuleArgs(action string, rule Rule) []string {
	var args []string
	args = append(args, strings.ToLower(action))

	toDest := normalizeAddress(rule.ToDest)
	toPort := normalizePort(rule.ToPort)
	toProto := rule.ToProtocol

	fromSource := normalizeAddress(rule.FromSource)
	fromPort := normalizePort(rule.FromPort)

	// For IPv6 rules, we need to use specific addresses or "any" will default to IPv4
	// UFW creates IPv6 rules when using IPv6 addresses or when the original rule was IPv6
	// We use "::/0" for IPv6 "anywhere" equivalent
	anyAddr := "any"
	if rule.IPv6 {
		anyAddr = "::/0" // IPv6 equivalent of "anywhere"
	}

	if fromSource != "" {
		args = append(args, "from", fromSource)
	} else {
		args = append(args, "from", anyAddr)
	}

	if toDest != "" {
		args = append(args, "to", toDest)
	} else {
		args = append(args, "to", anyAddr)
	}

	if toPort != "" && toPort != "any" {
		args = append(args, "port", toPort)
	}

	if toProto != "" && toProto != "any" {
		args = append(args, "proto", strings.ToLower(toProto))
	}

	if fromPort != "" && fromPort != "any" {
		args = append(args, "sport", fromPort)
	}

	return args
}

func GetCurrentRules() []Rule {
	numOut, _, err := RunSudo("status", "numbered")
	if err != nil {
		log.Printf("Error fetching rules: %v", err)
		return nil
	}
	return ParseRules(numOut)
}

func FindInsertPosition(rules []Rule, prevNextNum int, isIPv6 bool) int {
	if prevNextNum <= 0 {
		return 0
	}

	for _, r := range rules {
		if r.IPv6 == isIPv6 && r.Num >= 1 {
			return r.Num
		}
	}
	return 0
}

func MoveRule(rule Rule, direction int, newAction string) error {
	originalNum := rule.Num
	isIPv6 := rule.IPv6

	log.Printf("MoveRule: rule #%d, direction=%d, newAction=%s, isIPv6=%v", originalNum, direction, newAction, isIPv6)

	// Step 1: Get current rules to understand the landscape
	currentRules := GetCurrentRules()
	if currentRules == nil {
		return fmt.Errorf("failed to fetch current rules")
	}

	// Find our rule's neighbors of the same IP version
	var sameTypeRules []Rule
	for _, r := range currentRules {
		if r.IPv6 == isIPv6 {
			sameTypeRules = append(sameTypeRules, r)
		}
	}

	// Find index of our rule in the same-type list
	ourIndex := -1
	for i, r := range sameTypeRules {
		if r.Num == originalNum {
			ourIndex = i
			break
		}
	}

	if ourIndex == -1 {
		return fmt.Errorf("rule #%d not found in current rules", originalNum)
	}

	// Calculate target index in same-type list
	targetIndex := max(ourIndex + direction, 0)
	if targetIndex >= len(sameTypeRules) {
		targetIndex = len(sameTypeRules) - 1
	}

	// Step 2: Delete the rule
	_, stderr, err := DeleteRule(originalNum)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %s", stderr)
	}

	// Step 3: Re-fetch rules to get accurate positions
	newRules := GetCurrentRules()
	if newRules == nil {
		return fmt.Errorf("failed to fetch rules after delete")
	}

	// Find same-type rules again after deletion
	var newSameTypeRules []Rule
	for _, r := range newRules {
		if r.IPv6 == isIPv6 {
			newSameTypeRules = append(newSameTypeRules, r)
		}
	}

	// Step 4: Calculate actual insert position
	var insertPos int
	if len(newSameTypeRules) == 0 {
		insertPos = 0
	} else if targetIndex >= len(newSameTypeRules) {
		insertPos = 0
	} else {
		insertPos = newSameTypeRules[targetIndex].Num
	}

	log.Printf("MoveRule: targetIndex=%d, insertPos=%d, newSameTypeRules=%d", targetIndex, insertPos, len(newSameTypeRules))

	// Step 5: Validate insert position
	// UFW insert position must be between 1 and (total_rules + 1)
	totalRulesAfterDelete := len(newRules)
	if insertPos > totalRulesAfterDelete+1 {
		log.Printf("MoveRule: insertPos %d exceeds valid range (1-%d), appending instead", insertPos, totalRulesAfterDelete+1)
		insertPos = 0 // Signal to append
	}

	// Step 6: Insert the rule at the calculated position
	_, stderr, err = InsertRuleFromExisting(insertPos, newAction, rule)
	if err != nil {
		return fmt.Errorf("failed to insert rule (pos=%d, total=%d): %s", insertPos, totalRulesAfterDelete, stderr)
	}

	return nil
}

func GetInterfaces() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting interfaces: %v", err)
		return []string{}
	}

	var names []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		names = append(names, iface.Name)
	}
	return names
}

const (
	DirectionIn  = "in"
	DirectionOut = "out"
)

const (
	ProtoTCP = "tcp"
	ProtoUDP = "udp"
	ProtoAny = "any"
)

var Directions = []string{DirectionIn, DirectionOut}
var Protocols = []string{ProtoAny, ProtoTCP, ProtoUDP}
var CommonPorts = []struct {
	Port string
	Name string
}{
	{"22", "SSH"},
	{"80", "HTTP"},
	{"443", "HTTPS"},
	{"53", "DNS"},
	{"5432", "PostgreSQL"},
	{"6379", "Redis"},
	{"3000", "Dev Server"},
}

type AddRuleParams struct {
	Action    string // allow, deny, reject, limit
	Direction string // in, out, or empty for both
	Protocol  string // tcp, udp, any
	Port      string // port number or range
	FromAddr  string // source address or empty for any
	ToAddr    string // destination address or empty for any
	Interface string // network interface or empty for all
	Comment   string // optional comment
}

func BuildAddRuleCommand(p AddRuleParams) []string {
	var args []string

	args = append(args, strings.ToLower(p.Action))

	if p.Direction != "" {
		args = append(args, p.Direction)
	}

	if p.Interface != "" {
		args = append(args, "on", p.Interface)
	}

	if p.Port != "" {
		if p.FromAddr != "" && p.FromAddr != "any" {
			args = append(args, "from", p.FromAddr)
		} else {
			args = append(args, "from", "any")
		}

		if p.ToAddr != "" && p.ToAddr != "any" {
			args = append(args, "to", p.ToAddr)
		} else {
			args = append(args, "to", "any")
		}

		args = append(args, "port", p.Port)

		if p.Protocol != "" && p.Protocol != ProtoAny {
			args = append(args, "proto", p.Protocol)
		}
	} else {
		if p.FromAddr != "" && p.FromAddr != "any" {
			args = append(args, "from", p.FromAddr)
		}
		if p.ToAddr != "" && p.ToAddr != "any" {
			args = append(args, "to", p.ToAddr)
		}
	}

	if p.Comment != "" {
		args = append(args, "comment", p.Comment)
	}

	return args
}

func AddNewRule(p AddRuleParams) (stdout, stderr string, err error) {
	args := BuildAddRuleCommand(p)
	log.Printf("AddNewRule: %v", args)
	return RunSudo(args...)
}
