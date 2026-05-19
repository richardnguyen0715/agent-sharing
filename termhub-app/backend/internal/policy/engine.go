// Package policy provides a simple command policy engine for evaluating
// terminal input against configurable rules. It uses plain string containment
// matching — not regex — because pretending command parsing is real security
// would be dishonest.
package policy

import (
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Action represents the policy decision for a matched rule.
type Action string

const (
	// ActionBlock prevents the command from being sent to the PTY.
	ActionBlock Action = "block"
	// ActionWarn allows the command but logs a warning audit event.
	ActionWarn Action = "warn"
	// ActionAllow permits the command without any special handling.
	ActionAllow Action = "allow"
)

// PolicyResult holds the outcome of evaluating input against policy rules.
type PolicyResult struct {
	Action  Action
	Rule    string
	Message string
}

// Rule defines a single policy rule with a pattern to match and an action to take.
type Rule struct {
	Pattern string
	Action  Action
	Message string
}

// Engine evaluates terminal input against a set of policy rules.
type Engine struct {
	mu     sync.RWMutex
	rules  []Rule
	logger zerolog.Logger
}

// NewEngine creates a new Engine pre-loaded with default safety rules.
func NewEngine() *Engine {
	e := &Engine{
		logger: log.With().Str("component", "policy-engine").Logger(),
	}

	// Block rules — dangerous destructive commands.
	e.rules = []Rule{
		{Pattern: "rm -rf /", Action: ActionBlock, Message: "destructive command: rm -rf / is not allowed"},
		{Pattern: "shutdown", Action: ActionBlock, Message: "system shutdown is not allowed"},
		{Pattern: "reboot", Action: ActionBlock, Message: "system reboot is not allowed"},
		{Pattern: "mkfs", Action: ActionBlock, Message: "filesystem formatting is not allowed"},
		{Pattern: "dd if=", Action: ActionBlock, Message: "raw disk write via dd is not allowed"},
		{Pattern: ":(){ :|:& };:", Action: ActionBlock, Message: "fork bomb is not allowed"},

		// Warn rules — potentially dangerous but not blocked.
		{Pattern: "chmod 777", Action: ActionWarn, Message: "chmod 777 grants excessive permissions"},
		{Pattern: "rm -rf", Action: ActionWarn, Message: "recursive force delete detected"},
		{Pattern: "kill -9 1", Action: ActionWarn, Message: "killing PID 1 may destabilize the system"},
		{Pattern: "sudo su", Action: ActionWarn, Message: "escalating to root shell"},
	}

	e.logger.Info().Int("rule_count", len(e.rules)).Msg("policy engine initialized with default rules")
	return e
}

// Evaluate checks input against all rules and returns the first matching result.
// If no rule matches, ActionAllow is returned.
func (e *Engine) Evaluate(input string) PolicyResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return PolicyResult{Action: ActionAllow}
	}

	for _, rule := range e.rules {
		if strings.Contains(trimmed, rule.Pattern) {
			e.logger.Debug().
				Str("pattern", rule.Pattern).
				Str("action", string(rule.Action)).
				Str("input", trimmed).
				Msg("policy rule matched")
			return PolicyResult{
				Action:  rule.Action,
				Rule:    rule.Pattern,
				Message: rule.Message,
			}
		}
	}

	return PolicyResult{Action: ActionAllow}
}

// AddRule appends a new rule to the engine.
func (e *Engine) AddRule(pattern string, action Action, message string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, Rule{
		Pattern: pattern,
		Action:  action,
		Message: message,
	})
	e.logger.Info().
		Str("pattern", pattern).
		Str("action", string(action)).
		Msg("policy rule added")
}
