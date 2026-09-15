package tui

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jasperan/deepagents/gotui/internal/cli"
	"github.com/jasperan/deepagents/gotui/internal/threads"
)

// These tests drive the real ShellRunner against a stand-in CLI, because the
// fakeRunner used elsewhere stands in for the runner itself and so cannot test
// it: process spawning, stream attribution, timeouts and error extraction are
// exactly the parts that only a real subprocess exercises.
//
// The real deepagents CLI is not installed in this checkout (its package has no
// virtual environment and dependencies must not be installed here), so the
// stand-in emits the same envelopes its write_json produces.

// writeFakeCLI writes an executable script and returns its path.
func writeFakeCLI(t *testing.T, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in CLI is a POSIX shell script")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "deepagents")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatalf("write stand-in CLI: %v", err)
	}
	return path
}

func TestShellRunnerCaptureRunsTheRealCLIAndParsesItsEnvelope(t *testing.T) {
	binary := writeFakeCLI(t, `printf '%s\n' '{"schema_version": 1, "command": "threads list", "data": [{"thread_id": "t-1", "agent_name": "coder", "message_count": 3}]}'`)

	runner := ShellRunner{CaptureTimeout: 10 * time.Second}
	inv := cli.ThreadsList(cli.Launcher{Bin: binary}, cli.ThreadsListOptions{})

	out, err := runner.Capture(context.Background(), inv)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	list, err := threads.ParseThreadsList([]byte(out))
	if err != nil {
		t.Fatalf("ParseThreadsList: %v", err)
	}
	if len(list) != 1 || list[0].ThreadID != "t-1" || list[0].AgentLabel() != "coder" {
		t.Fatalf("list = %+v, want the thread the CLI emitted", list)
	}
}

// A failing CLI is the expected offline outcome. The user must see the CLI's own
// message rather than a bare exit status.
func TestShellRunnerCaptureSurfacesTheCLIsOwnError(t *testing.T) {
	binary := writeFakeCLI(t, `echo "ValueError: DA_CLI_RECENT_THREADS is invalid" >&2; exit 1`)

	runner := ShellRunner{CaptureTimeout: 10 * time.Second}
	_, err := runner.Capture(context.Background(),
		cli.ThreadsList(cli.Launcher{Bin: binary}, cli.ThreadsListOptions{}))
	if err == nil {
		t.Fatal("Capture reported success for a failing CLI")
	}
	if !strings.Contains(err.Error(), "ValueError") {
		t.Errorf("error %q does not quote the CLI's own message", err)
	}
}

func TestShellRunnerCaptureReportsAMissingBinary(t *testing.T) {
	runner := ShellRunner{CaptureTimeout: 5 * time.Second}
	_, err := runner.Capture(context.Background(), cli.Invocation{Argv: []string{"/nonexistent/deepagents", "threads", "list"}})
	if err == nil {
		t.Fatal("Capture reported success for a missing binary")
	}
}

// Without a bound, an unreachable model endpoint would hang the UI forever,
// which is the failure mode this front-end exists to avoid.
func TestShellRunnerCaptureHonoursItsTimeout(t *testing.T) {
	binary := writeFakeCLI(t, `sleep 30`)

	runner := ShellRunner{CaptureTimeout: 200 * time.Millisecond}
	start := time.Now()
	_, err := runner.Capture(context.Background(),
		cli.ThreadsList(cli.Launcher{Bin: binary}, cli.ThreadsListOptions{}))
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Capture reported success for a command that never finished")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("error %q does not report the timeout", err)
	}
	if elapsed > 10*time.Second {
		t.Errorf("Capture took %s, so the timeout was not enforced", elapsed)
	}
}

// Stream must attribute each line to the stream it came from, or the UI cannot
// tell progress from a problem.
func TestShellRunnerStreamAttributesLinesToTheirStream(t *testing.T) {
	binary := writeFakeCLI(t, `echo "an answer"; echo "a warning" >&2; exit 0`)

	runner := ShellRunner{RunTimeout: 10 * time.Second}
	var got []Line
	err := runner.Stream(context.Background(),
		cli.Invocation{Argv: []string{binary}}, func(l Line) { got = append(got, l) })
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d lines, want 2: %+v", len(got), got)
	}
	var sawStdout, sawStderr bool
	for _, l := range got {
		switch l.Stream {
		case "stdout":
			sawStdout = l.Text == "an answer" && !l.IsStderr()
		case "stderr":
			sawStderr = l.Text == "a warning" && l.IsStderr()
		}
	}
	if !sawStdout || !sawStderr {
		t.Errorf("stream attribution is wrong: %+v", got)
	}
}

func TestShellRunnerStreamReportsAFailingRun(t *testing.T) {
	binary := writeFakeCLI(t, `echo "RuntimeError: model unavailable" >&2; exit 3`)

	runner := ShellRunner{RunTimeout: 10 * time.Second}
	err := runner.Stream(context.Background(), cli.Invocation{Argv: []string{binary}}, nil)
	if err == nil {
		t.Fatal("Stream reported success for a failing run")
	}
	if !strings.Contains(err.Error(), "status 3") {
		t.Errorf("error %q does not report the exit status", err)
	}
	if !strings.Contains(err.Error(), "model unavailable") {
		t.Errorf("error %q does not quote the CLI's own message", err)
	}
}

// A cancelled context must take the whole tree down, not just the direct child:
// the CLI spawns a server subprocess that would otherwise hold the pipe open and
// leave the reader blocked past the deadline.
func TestShellRunnerStreamIsCancelledPromptly(t *testing.T) {
	binary := writeFakeCLI(t, `sleep 30 & sleep 30`)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	runner := ShellRunner{RunTimeout: time.Minute}
	start := time.Now()
	err := runner.Stream(ctx, cli.Invocation{Argv: []string{binary}}, nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Stream reported success for a cancelled run")
	}
	if elapsed > 15*time.Second {
		t.Errorf("Stream took %s to stop, so cancellation did not reach the process tree", elapsed)
	}
}

func TestShellRunnerRefusesAnEmptyInvocation(t *testing.T) {
	runner := ShellRunner{}
	if _, err := runner.Capture(context.Background(), cli.Invocation{}); err == nil {
		t.Error("Capture accepted an empty invocation")
	}
	if err := runner.Stream(context.Background(), cli.Invocation{}, nil); err == nil {
		t.Error("Stream accepted an empty invocation")
	}
}

// A long line must not silently truncate the stream: the scanner's default
// buffer is 64KiB and agent output can exceed it.
func TestShellRunnerStreamHandlesVeryLongLines(t *testing.T) {
	binary := writeFakeCLI(t, `printf '%s' "$(head -c 200000 /dev/zero | tr '\0' 'x')"; echo`)

	runner := ShellRunner{RunTimeout: 30 * time.Second}
	longest := 0
	err := runner.Stream(context.Background(),
		cli.Invocation{Argv: []string{binary}}, func(l Line) {
			if len(l.Text) > longest {
				longest = len(l.Text)
			}
		})
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if longest < 100000 {
		t.Errorf("longest line seen was %d bytes; the scanner buffer truncated the stream", longest)
	}
}

// firstMeaningfulLine must prefer the CLI's own exception header over a traceback
// frame, which is a code location rather than the failure.
func TestFirstMeaningfulLinePrefersTheExceptionHeader(t *testing.T) {
	output := strings.Join([]string{
		"Traceback (most recent call last):",
		`  File "sessions.py", line 1116, in list_threads_command`,
		"    raise ValueError(message)",
		"ValueError: DA_CLI_RECENT_THREADS is invalid",
	}, "\n")

	got := firstMeaningfulLine(output)
	if !strings.Contains(got, "ValueError: DA_CLI_RECENT_THREADS is invalid") {
		t.Errorf("firstMeaningfulLine = %q, want the exception header", got)
	}
}

// Rich output is box-drawn, so borders must not be mistaken for content.
func TestFirstMeaningfulLineStripsRichBorders(t *testing.T) {
	output := "\x1b[31m│\x1b[0m Error: the service is unreachable\n"
	got := firstMeaningfulLine(output)
	if strings.ContainsAny(got, "│\x1b") {
		t.Errorf("firstMeaningfulLine = %q, want free of borders and escape codes", got)
	}
	if !strings.Contains(got, "Error: the service is unreachable") {
		t.Errorf("firstMeaningfulLine = %q, want the message", got)
	}
}

func TestFirstMeaningfulLineFallsBackToTheFirstContentLine(t *testing.T) {
	if got := firstMeaningfulLine("\n\n  something went wrong  \n"); got != "something went wrong" {
		t.Errorf("firstMeaningfulLine = %q", got)
	}
	if got := firstMeaningfulLine("   \n\n"); got != "" {
		t.Errorf("firstMeaningfulLine(blank) = %q, want empty", got)
	}
}

func TestFirstNonEmptySkipsLeadingBlanks(t *testing.T) {
	if got := firstNonEmpty([]string{"", "  ", "real"}); got != "real" {
		t.Errorf("firstNonEmpty = %q, want real", got)
	}
	if got := firstNonEmpty([]string{"", "  "}); got != "" {
		t.Errorf("firstNonEmpty(blank) = %q, want empty", got)
	}
}
