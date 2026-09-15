package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestThreadsListAlwaysAsksForJSON(t *testing.T) {
	inv := ThreadsList(Launcher{Bin: "deepagents"}, ThreadsListOptions{})
	got := inv.CommandLine()
	if !strings.Contains(got, "--json") {
		t.Fatalf("argv %q does not request JSON; this front-end must never scrape the Rich table", got)
	}
	if strings.HasPrefix(got, "deepagents threads list") == false {
		t.Errorf("argv %q is not the documented subcommand", got)
	}
}

func TestThreadsListOnlyPassesFlagsItHasValuesFor(t *testing.T) {
	// An omitted flag must stay omitted: the CLI reads its own config for sort
	// order and its own environment for the limit, and sending a value would
	// silently override the user's saved preference.
	empty := ThreadsList(Launcher{Bin: "deepagents"}, ThreadsListOptions{}).CommandLine()
	for _, unwanted := range []string{"--agent", "--limit", "--sort", "--branch", "-n ", "-v"} {
		if strings.Contains(empty, unwanted) {
			t.Errorf("empty options produced %q, which contains %q", empty, unwanted)
		}
	}

	full := ThreadsList(Launcher{Bin: "deepagents"}, ThreadsListOptions{
		Agent:   "coder",
		Limit:   7,
		Sort:    "created",
		Branch:  "main",
		Verbose: true,
	}).CommandLine()
	for _, want := range []string{"--agent coder", "-n 7", "--sort created", "--branch main", "-v"} {
		if !strings.Contains(full, want) {
			t.Errorf("argv %q is missing %q", full, want)
		}
	}
}

func TestThreadsDeleteAndAgentsListUseTheDocumentedSubcommands(t *testing.T) {
	got := ThreadsDelete(Launcher{Bin: "deepagents"}, "abc-123").CommandLine()
	if !strings.Contains(got, "threads delete abc-123") || !strings.Contains(got, "--json") {
		t.Errorf("argv %q is not `threads delete <id> --json`", got)
	}
	if got := AgentsList(Launcher{Bin: "deepagents"}).CommandLine(); got != "deepagents list --json" {
		t.Errorf("argv = %q, want `deepagents list --json`", got)
	}
}

func TestRunIsNonInteractiveAndQuiet(t *testing.T) {
	inv, err := Run(Launcher{Bin: "deepagents"}, RunOptions{Task: "summarize", Quiet: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	got := inv.CommandLine()
	// -n is what makes the CLI run one task and exit; -q is what keeps the
	// answer alone on stdout. Both are required for this front-end to be usable
	// in a pipe.
	for _, want := range []string{"-n summarize", "-q"} {
		if !strings.Contains(got, want) {
			t.Errorf("argv %q is missing %q", got, want)
		}
	}
}

func TestRunPassesTheTaskAsOneArgument(t *testing.T) {
	// A task containing spaces must survive as a single argv element; quoting it
	// here would send the quotes to the CLI as part of the task.
	task := `find "all" the TODOs and --watch out`
	inv, err := Run(Launcher{Bin: "deepagents"}, RunOptions{Task: task})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(inv.Argv) != 3 || inv.Argv[1] != "-n" || inv.Argv[2] != task {
		t.Fatalf("argv = %q, want [deepagents -n <task>]", inv.Argv)
	}
}

func TestRunRefusesAnEmptyTask(t *testing.T) {
	for _, task := range []string{"", "   ", "\t"} {
		if _, err := Run(Launcher{Bin: "deepagents"}, RunOptions{Task: task}); err == nil {
			t.Errorf("Run accepted the empty task %q", task)
		}
	}
}

func TestRunOmitsFlagsLeftAtTheirZeroValue(t *testing.T) {
	inv, err := Run(Launcher{Bin: "deepagents"}, RunOptions{Task: "t"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	got := inv.CommandLine()
	// -y is the dangerous one: it must never be sent unless the user asked.
	for _, unwanted := range []string{"-a ", "-M ", "-y", "-S ", "--no-stream"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("argv %q contains %q, which the user did not request", got, unwanted)
		}
	}
}

func TestRunSendsEveryRequestedOption(t *testing.T) {
	inv, err := Run(Launcher{Bin: "deepagents"}, RunOptions{
		Task:           "t",
		Agent:          "coder",
		Model:          "claude-sonnet-4-6",
		Quiet:          true,
		NoStream:       true,
		AutoApprove:    true,
		ShellAllowList: "recommended",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	got := inv.CommandLine()
	for _, want := range []string{"-a coder", "-M claude-sonnet-4-6", "-y", "-S recommended", "--no-stream"} {
		if !strings.Contains(got, want) {
			t.Errorf("argv %q is missing %q", got, want)
		}
	}
}

// The uv path is what makes a source checkout usable without installing the CLI
// globally, so the prefix has to land between the program and the arguments.
func TestLauncherPrefixIsSplicedAfterTheProgram(t *testing.T) {
	l := Launcher{Bin: "uv", Prefix: []string{"run", "deepagents"}, Dir: "/repo/libs/cli"}
	inv := ThreadsList(l, ThreadsListOptions{})
	want := "uv run deepagents threads list --json"
	if inv.CommandLine() != want {
		t.Errorf("argv = %q, want %q", inv.CommandLine(), want)
	}
	if inv.Dir != "/repo/libs/cli" {
		t.Errorf("Dir = %q, want the project directory", inv.Dir)
	}
}

func TestResolveHonoursTheExplicitOverride(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "deepagents")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write fake: %v", err)
	}

	l, err := Resolve(fake)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", fake, err)
	}
	if l.Bin != fake {
		t.Errorf("Bin = %q, want %q", l.Bin, fake)
	}
}

func TestResolveHonoursTheEnvironmentOverride(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "deepagents-elsewhere")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write fake: %v", err)
	}
	t.Setenv(EnvBinary, fake)

	l, err := Resolve("")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if l.Bin != fake {
		t.Errorf("Bin = %q, want the %s value %q", l.Bin, EnvBinary, fake)
	}
}

// A missing CLI must produce the actionable error, not an exec failure later.
func TestResolveReportsAMissingCLI(t *testing.T) {
	t.Setenv(EnvBinary, "deepagents-definitely-not-installed-xyz")
	_, err := Resolve("")
	if err == nil {
		t.Fatal("Resolve accepted a binary that does not exist")
	}
	if !strings.Contains(err.Error(), "deepagents") {
		t.Errorf("error %q does not name the missing CLI", err)
	}
}

// No invocation may carry a secret, because argv is visible in a process list.
func TestNoInvocationLeaksAConfiguredPassword(t *testing.T) {
	t.Setenv("DEEPAGENTS_ORACLE_PASSWORD", "hunter2secret")
	inv, err := Run(Launcher{Bin: "deepagents"}, RunOptions{Task: "t"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if strings.Contains(inv.CommandLine(), "hunter2secret") {
		t.Fatalf("argv leaks the password: %q", inv.CommandLine())
	}
	// The environment is inherited, which is how the child gets the credential.
	if !containsEntry(inv.Env, "DEEPAGENTS_ORACLE_PASSWORD=hunter2secret") {
		t.Error("the child did not inherit DEEPAGENTS_ORACLE_PASSWORD, so it cannot authenticate")
	}
}

func TestRedactMasksSecretsButNotOrdinaryValues(t *testing.T) {
	env := []string{
		"DEEPAGENTS_ORACLE_PASSWORD=hunter2secret",
		"DEEPAGENTS_ORACLE_USER=deepagents",
		"PATH=/usr/bin",
	}
	got := Redact("connecting as deepagents with hunter2secret", env)
	if strings.Contains(got, "hunter2secret") {
		t.Errorf("Redact left the password in %q", got)
	}
	if !strings.Contains(got, "deepagents") {
		t.Errorf("Redact removed a non-secret value: %q", got)
	}
}

func TestEnvWithReplacesRatherThanDuplicates(t *testing.T) {
	t.Setenv("DEEPAGENTS_ORACLE_MODE", "freepdb")
	env := EnvWith(map[string]string{"DEEPAGENTS_ORACLE_MODE": "adb"})
	count := 0
	for _, kv := range env {
		if strings.HasPrefix(kv, "DEEPAGENTS_ORACLE_MODE=") {
			count++
			if kv != "DEEPAGENTS_ORACLE_MODE=adb" {
				t.Errorf("entry = %q, want the override", kv)
			}
		}
	}
	if count != 1 {
		t.Errorf("found %d entries for the overridden variable, want 1", count)
	}
}

func containsEntry(env []string, want string) bool {
	for _, kv := range env {
		if kv == want {
			return true
		}
	}
	return false
}
