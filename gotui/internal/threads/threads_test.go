package threads

import (
	"strings"
	"testing"
)

// listEnvelope is a real `threads list --json` response, shaped exactly as
// write_json produces it in deepagents_cli/output.py:
// {"schema_version": 1, "command": "threads list", "data": [...]}.
const listEnvelope = `{"schema_version": 1, "command": "threads list", "data": [
  {"thread_id": "018f-uuid-1", "agent_name": "coder", "updated_at": "2026-09-14 10:04:05.123456+00:00",
   "created_at": "2026-09-14 09:00:00+00:00", "git_branch": "main", "initial_prompt": "Summarize the repo",
   "message_count": 12, "latest_checkpoint_id": "1ef0", "cwd": "/repo"},
  {"thread_id": "018f-uuid-2", "agent_name": null, "updated_at": null, "message_count": 0}
]}`

const emptyEnvelope = `{"schema_version": 1, "command": "threads list", "data": []}`

func TestParseThreadsListReadsEveryField(t *testing.T) {
	list, err := ParseThreadsList([]byte(listEnvelope))
	if err != nil {
		t.Fatalf("ParseThreadsList: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("got %d threads, want 2", len(list))
	}
	first := list[0]
	if first.ThreadID != "018f-uuid-1" {
		t.Errorf("ThreadID = %q", first.ThreadID)
	}
	if first.AgentLabel() != "coder" {
		t.Errorf("AgentLabel = %q", first.AgentLabel())
	}
	if first.Prompt() != "Summarize the repo" {
		t.Errorf("Prompt = %q", first.Prompt())
	}
	if first.Messages() != "12" {
		t.Errorf("Messages = %q, want 12", first.Messages())
	}
	if first.Branch() != "main" {
		t.Errorf("Branch = %q", first.Branch())
	}
	if first.WorkingDir() != "/repo" {
		t.Errorf("WorkingDir = %q", first.WorkingDir())
	}
}

// The CLI omits members it does not know, and nulls are real. Neither may be
// mistaken for a zero value: a thread with 0 messages is not a missing count.
func TestParseThreadsListToleratesAbsentAndNullMembers(t *testing.T) {
	list, err := ParseThreadsList([]byte(listEnvelope))
	if err != nil {
		t.Fatalf("ParseThreadsList: %v", err)
	}
	second := list[1]
	if second.AgentLabel() == "" {
		t.Error("a null agent rendered as an empty label; a placeholder is needed")
	}
	if second.Updated() != "-" {
		t.Errorf("Updated = %q, want a placeholder for null", second.Updated())
	}
	if second.Messages() != "0" {
		t.Errorf("Messages = %q, want 0 for a present zero", second.Messages())
	}
	if second.Branch() != "-" {
		t.Errorf("Branch = %q, want a placeholder for an absent key", second.Branch())
	}
}

// A fresh install has no threads. That is an empty state, not a failure.
func TestParseThreadsListAcceptsAnEmptyList(t *testing.T) {
	list, err := ParseThreadsList([]byte(emptyEnvelope))
	if err != nil {
		t.Fatalf("ParseThreadsList: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("got %d threads, want 0", len(list))
	}
}

// A response to a different question must never be read as an answer to this
// one, because the payload shapes differ.
func TestParseRefusesTheWrongCommand(t *testing.T) {
	wrong := `{"schema_version": 1, "command": "list", "data": []}`
	if _, err := ParseThreadsList([]byte(wrong)); err == nil {
		t.Fatal("ParseThreadsList accepted a `list` response")
	}
}

// A future envelope must be refused, not interpreted on a guess.
func TestParseRefusesAnUnknownSchemaVersion(t *testing.T) {
	future := `{"schema_version": 2, "command": "threads list", "data": []}`
	_, err := ParseThreadsList([]byte(future))
	if err == nil {
		t.Fatal("ParseThreadsList accepted schema_version 2")
	}
	if !strings.Contains(err.Error(), "schema_version") {
		t.Errorf("error %q does not name the version mismatch", err)
	}
}

// If the CLI printed a human table, the error should say so rather than
// reporting a bare JSON syntax error.
func TestParseExplainsNonJSONOutput(t *testing.T) {
	_, err := ParseThreadsList([]byte("╭──────────────╮\n│ Threads      │\n"))
	if err == nil {
		t.Fatal("ParseThreadsList accepted a Rich table")
	}
	if !strings.Contains(err.Error(), "deepagents") {
		t.Errorf("error %q does not mention the CLI", err)
	}
}

func TestParseExplainsEmptyOutput(t *testing.T) {
	if _, err := ParseThreadsList([]byte("  \n")); err == nil {
		t.Fatal("ParseThreadsList accepted empty output")
	}
}

func TestParseThreadsDelete(t *testing.T) {
	ok := `{"schema_version": 1, "command": "threads delete", "data": {"thread_id": "abc", "deleted": true}}`
	result, err := ParseThreadsDelete([]byte(ok))
	if err != nil {
		t.Fatalf("ParseThreadsDelete: %v", err)
	}
	if !result.Deleted || result.ThreadID != "abc" {
		t.Errorf("result = %+v, want abc deleted", result)
	}

	// "not found" is a normal outcome the CLI reports with deleted=false, so it
	// must not be treated as a parse failure.
	notFound := `{"schema_version": 1, "command": "threads delete", "data": {"thread_id": "abc", "deleted": false}}`
	result, err = ParseThreadsDelete([]byte(notFound))
	if err != nil {
		t.Fatalf("ParseThreadsDelete: %v", err)
	}
	if result.Deleted {
		t.Error("deleted = true for a not-found response")
	}
}

func TestParseAgentsAndDefaultDetection(t *testing.T) {
	raw := `{"schema_version": 1, "command": "list", "data": [
	  {"name": "coder", "path": "/home/u/.deepagents/coder", "has_agents_md": true, "is_default": true},
	  {"name": "researcher", "path": "/home/u/.deepagents/researcher", "has_agents_md": false, "is_default": false}
	]}`
	agents, err := ParseAgents([]byte(raw))
	if err != nil {
		t.Fatalf("ParseAgents: %v", err)
	}
	names := AgentNames(agents)
	if len(names) != 2 || names[0] != "coder" || names[1] != "researcher" {
		t.Errorf("AgentNames = %v", names)
	}
	if got := DefaultAgentName(agents); got != "coder" {
		t.Errorf("DefaultAgentName = %q, want coder", got)
	}
}

func TestDefaultAgentNameIsEmptyWhenNothingIsMarkedDefault(t *testing.T) {
	agents := []AgentInfo{{Name: "coder"}}
	if got := DefaultAgentName(agents); got != "" {
		t.Errorf("DefaultAgentName = %q, want empty", got)
	}
}

// write_json uses default=str, so an unformatted datetime arrives with a space
// separator rather than RFC 3339's "T", and may carry microseconds.
func TestFormatTimestampAcceptsEveryShapeTheCLIEmits(t *testing.T) {
	cases := map[string]string{
		"2026-09-14T10:04:05Z":             "2026-09-14",
		"2026-09-14T10:04:05+00:00":        "2026-09-14",
		"2026-09-14 10:04:05.123456+00:00": "2026-09-14",
		"2026-09-14 10:04:05":              "2026-09-14",
		"2026-09-14 10:04:05.123456":       "2026-09-14",
		"2026-09-14":                       "2026-09-14",
	}
	for input, wantPrefix := range cases {
		got := FormatTimestamp(input)
		if !strings.HasPrefix(got, wantPrefix) {
			t.Errorf("FormatTimestamp(%q) = %q, want it to start with %q", input, got, wantPrefix)
		}
	}
}

// An unparseable timestamp is still worth showing, and is never worth crashing
// over.
func TestFormatTimestampFallsBackToTheRawValue(t *testing.T) {
	if got := FormatTimestamp("not a timestamp"); got != "not a timestamp" {
		t.Errorf("FormatTimestamp = %q, want the raw value", got)
	}
	if got := FormatTimestamp(""); got != "-" {
		t.Errorf("FormatTimestamp(\"\") = %q, want -", got)
	}
}

func TestFilterThreadsMatchesPromptIDAndAgent(t *testing.T) {
	list, err := ParseThreadsList([]byte(listEnvelope))
	if err != nil {
		t.Fatalf("ParseThreadsList: %v", err)
	}

	if got := FilterThreads(list, ""); len(got) != 2 {
		t.Errorf("an empty query returned %d threads, want all 2", len(got))
	}
	if got := FilterThreads(list, "SUMMARIZE"); len(got) != 1 || got[0].ThreadID != "018f-uuid-1" {
		t.Errorf("case-insensitive prompt match returned %v", got)
	}
	if got := FilterThreads(list, "uuid-2"); len(got) != 1 || got[0].ThreadID != "018f-uuid-2" {
		t.Errorf("id match returned %v", got)
	}
	if got := FilterThreads(list, "coder"); len(got) != 1 {
		t.Errorf("agent match returned %v", got)
	}
	if got := FilterThreads(list, "no-such-thread"); len(got) != 0 {
		t.Errorf("a non-matching query returned %v", got)
	}
}
