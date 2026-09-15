// Package threads decodes the deepagents CLI's machine-readable output.
//
// Every subcommand invoked with `--json` writes one single-line envelope,
// produced by deepagents_cli.output.write_json:
//
//	{"schema_version": 1, "command": "...", "data": ...}
//
// The envelope is the CLI's documented, stable contract, so this package parses
// it rather than scraping the Rich tables. A schema version this front-end does
// not know is refused explicitly instead of being interpreted on a guess.
package threads

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SupportedSchemaVersion is the only envelope version this package understands.
const SupportedSchemaVersion = 1

// Command names, exactly as write_json records them.
const (
	CommandThreadsList   = "threads list"
	CommandThreadsDelete = "threads delete"
	CommandAgentsList    = "list"
)

// envelope mirrors write_json's output.
type envelope struct {
	SchemaVersion int             `json:"schema_version"`
	Command       string          `json:"command"`
	Data          json.RawMessage `json:"data"`
}

// decode validates the envelope and returns its payload.
//
// wantCommand pins the command: a response to a different question must never be
// read as an answer to this one, because the payload shapes differ.
func decode(raw []byte, wantCommand string) (json.RawMessage, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return nil, errors.New("the deepagents CLI produced no output; run it directly to see why")
	}
	var env envelope
	if err := json.Unmarshal([]byte(trimmed), &env); err != nil {
		// The most likely cause is that the CLI ignored --json and printed a
		// human table, so say that rather than reporting a bare syntax error.
		return nil, fmt.Errorf("the CLI output was not the expected JSON envelope "+
			"(is this the deepagents CLI?): %w", err)
	}
	if env.SchemaVersion != SupportedSchemaVersion {
		return nil, fmt.Errorf("unsupported JSON schema_version %d (this front-end speaks %d)",
			env.SchemaVersion, SupportedSchemaVersion)
	}
	if env.Command != wantCommand {
		return nil, fmt.Errorf("expected the %q response, got %q", wantCommand, env.Command)
	}
	return env.Data, nil
}

// ThreadInfo mirrors the ThreadInfo TypedDict returned by `threads list`.
//
// Optional members are pointers because the CLI omits keys it does not know,
// and a zero value would be indistinguishable from a real 0 (message_count) or
// an empty name.
type ThreadInfo struct {
	ThreadID           string  `json:"thread_id"`
	AgentName          *string `json:"agent_name"`
	UpdatedAt          *string `json:"updated_at"`
	CreatedAt          *string `json:"created_at"`
	GitBranch          *string `json:"git_branch"`
	InitialPrompt      *string `json:"initial_prompt"`
	MessageCount       *int    `json:"message_count"`
	LatestCheckpointID *string `json:"latest_checkpoint_id"`
	Cwd                *string `json:"cwd"`
}

// ParseThreadsList decodes a `threads list --json` response.
//
// An empty list is a normal, and common, result: a fresh install has no
// threads. It is returned as an empty slice rather than an error so callers
// render an empty state instead of a failure.
func ParseThreadsList(raw []byte) ([]ThreadInfo, error) {
	data, err := decode(raw, CommandThreadsList)
	if err != nil {
		return nil, err
	}
	var list []ThreadInfo
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("could not read the thread list: %w", err)
	}
	return list, nil
}

// DeleteResult mirrors the payload of `threads delete --json`.
//
// Deleted is false when the id was not found, which the CLI reports as a normal
// outcome and not an error.
type DeleteResult struct {
	ThreadID string `json:"thread_id"`
	Deleted  bool   `json:"deleted"`
}

// ParseThreadsDelete decodes a `threads delete --json` response.
func ParseThreadsDelete(raw []byte) (DeleteResult, error) {
	data, err := decode(raw, CommandThreadsDelete)
	if err != nil {
		return DeleteResult{}, err
	}
	var result DeleteResult
	if err := json.Unmarshal(data, &result); err != nil {
		return DeleteResult{}, fmt.Errorf("could not read the delete result: %w", err)
	}
	return result, nil
}

// AgentInfo mirrors one row of `deepagents list --json`.
type AgentInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	HasAgentsMD bool   `json:"has_agents_md"`
	IsDefault   bool   `json:"is_default"`
}

// ParseAgents decodes a `list --json` response.
func ParseAgents(raw []byte) ([]AgentInfo, error) {
	data, err := decode(raw, CommandAgentsList)
	if err != nil {
		return nil, err
	}
	var list []AgentInfo
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("could not read the agent list: %w", err)
	}
	return list, nil
}

// AgentNames returns the agent names, for populating a select field.
func AgentNames(agents []AgentInfo) []string {
	names := make([]string, 0, len(agents))
	for _, a := range agents {
		if a.Name != "" {
			names = append(names, a.Name)
		}
	}
	return names
}

// DefaultAgentName returns the agent the CLI would pick on its own, or "".
func DefaultAgentName(agents []AgentInfo) string {
	for _, a := range agents {
		if a.IsDefault && a.Name != "" {
			return a.Name
		}
	}
	return ""
}

// AgentLabel is the display name for a thread's agent.
func (t ThreadInfo) AgentLabel() string {
	if t.AgentName == nil || strings.TrimSpace(*t.AgentName) == "" {
		return "(no agent recorded)"
	}
	return *t.AgentName
}

// Prompt returns the thread's opening prompt, or a placeholder.
func (t ThreadInfo) Prompt() string {
	if t.InitialPrompt == nil || strings.TrimSpace(*t.InitialPrompt) == "" {
		return "(no prompt recorded)"
	}
	return strings.Join(strings.Fields(*t.InitialPrompt), " ")
}

// Messages returns the message count as display text.
func (t ThreadInfo) Messages() string {
	if t.MessageCount == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *t.MessageCount)
}

// Branch returns the git branch, or a placeholder.
func (t ThreadInfo) Branch() string {
	if t.GitBranch == nil || strings.TrimSpace(*t.GitBranch) == "" {
		return "-"
	}
	return *t.GitBranch
}

// WorkingDir returns the recorded working directory, or a placeholder.
func (t ThreadInfo) WorkingDir() string {
	if t.Cwd == nil || strings.TrimSpace(*t.Cwd) == "" {
		return "-"
	}
	return *t.Cwd
}

// Updated returns the last-update timestamp, or a placeholder.
func (t ThreadInfo) Updated() string { return formatOptional(t.UpdatedAt) }

// Created returns the creation timestamp, or a placeholder.
func (t ThreadInfo) Created() string { return formatOptional(t.CreatedAt) }

func formatOptional(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "-"
	}
	return FormatTimestamp(*value)
}

// timestampLayouts are the shapes a timestamp can arrive in.
//
// write_json is called with default=str, so a datetime that has not already been
// formatted by the CLI is serialized with str(), which uses a space separator
// rather than RFC 3339's "T" and may carry microseconds. Both are accepted.
var timestampLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// FormatTimestamp renders a timestamp compactly for a terminal, falling back to
// the raw string when it cannot be parsed: an unparseable timestamp is still
// worth showing, and is never worth crashing over.
func FormatTimestamp(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "-"
	}
	normalized := trimmed
	if len(normalized) > 10 && normalized[10] == 'T' {
		// Already RFC 3339; leave it alone.
	} else if len(normalized) > 10 && normalized[10] == ' ' {
		normalized = normalized[:10] + "T" + normalized[11:]
	}
	for _, layout := range timestampLayouts {
		if parsed, err := time.Parse(layout, normalized); err == nil {
			return parsed.Local().Format("2006-01-02 15:04")
		}
	}
	return trimmed
}

// FilterThreads keeps the threads matching a case-insensitive substring of the
// prompt, id or agent. An empty query matches everything.
//
// This is a view filter over data the CLI already returned; it never changes
// what is fetched, so it cannot disagree with the CLI about what exists.
func FilterThreads(list []ThreadInfo, query string) []ThreadInfo {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return list
	}
	out := make([]ThreadInfo, 0, len(list))
	for _, t := range list {
		haystack := strings.ToLower(t.ThreadID + " " + t.AgentLabel() + " " + t.Prompt() + " " + t.Branch())
		if strings.Contains(haystack, q) {
			out = append(out, t)
		}
	}
	return out
}
