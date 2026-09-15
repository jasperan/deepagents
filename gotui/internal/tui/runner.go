package tui

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jasperan/deepagents/gotui/internal/cli"
)

// Line is one streamed line of subprocess output.
type Line struct {
	// Stream is "stdout" or "stderr".
	Stream string
	// Text is the line without its trailing newline.
	Text string
}

// IsStderr reports whether the line came from the error stream, so the UI can
// style it as a warning without re-deriving it.
func (l Line) IsStderr() bool { return l.Stream == "stderr" }

// Runner executes the real deepagents CLI.
//
// It is an interface so the browser and the run view can be tested without the
// CLI installed, which is the only way they can be tested in a checkout that has
// no CLI virtual environment.
type Runner interface {
	// Capture runs a short read-only command and returns its combined output.
	Capture(ctx context.Context, inv cli.Invocation) (string, error)
	// Stream runs a command and calls emit for every line as it arrives, so a
	// long agent run shows progress instead of a frozen screen.
	Stream(ctx context.Context, inv cli.Invocation, emit func(Line)) error
}

// DefaultTimeout bounds a read-only call such as `threads list`.
const DefaultTimeout = 60 * time.Second

// DefaultRunTimeout bounds one agent run.
//
// An agent run can legitimately take many minutes, but an unbounded wait would
// turn "the model endpoint is unreachable" into a hang, which is the failure
// mode this front-end exists to avoid.
const DefaultRunTimeout = 30 * time.Minute

// ShellRunner runs the CLI as a subprocess.
type ShellRunner struct {
	// CaptureTimeout bounds Capture; zero means DefaultTimeout.
	CaptureTimeout time.Duration
	// RunTimeout bounds Stream; zero means DefaultRunTimeout.
	RunTimeout time.Duration
}

// NewShellRunner returns a runner with the default bounds.
func NewShellRunner() ShellRunner {
	return ShellRunner{CaptureTimeout: DefaultTimeout, RunTimeout: DefaultRunTimeout}
}

func (r ShellRunner) captureTimeout() time.Duration {
	if r.CaptureTimeout <= 0 {
		return DefaultTimeout
	}
	return r.CaptureTimeout
}

func (r ShellRunner) runTimeout() time.Duration {
	if r.RunTimeout <= 0 {
		return DefaultRunTimeout
	}
	return r.RunTimeout
}

// Capture runs the invocation and returns its combined output.
func (r ShellRunner) Capture(ctx context.Context, inv cli.Invocation) (string, error) {
	if len(inv.Argv) == 0 {
		return "", errors.New("empty invocation")
	}
	ctx, cancel := context.WithTimeout(ctx, r.captureTimeout())
	defer cancel()

	cmd := command(ctx, inv)
	out, err := cmd.CombinedOutput()
	text := string(out)
	if ctx.Err() == context.DeadlineExceeded {
		return text, fmt.Errorf("timed out after %s", r.captureTimeout())
	}
	if err != nil {
		return text, describeFailure(cmd, err, text)
	}
	return text, nil
}

// Stream runs the invocation and emits each line as it arrives.
//
// emit is called serially, never concurrently, so a caller may safely append to
// a slice or mutate a model from it. Ordering between stdout and stderr is
// best-effort, since they are separate pipes.
func (r ShellRunner) Stream(ctx context.Context, inv cli.Invocation, emit func(Line)) error {
	if len(inv.Argv) == 0 {
		return errors.New("empty invocation")
	}
	if emit == nil {
		emit = func(Line) {}
	}
	ctx, cancel := context.WithTimeout(ctx, r.runTimeout())
	defer cancel()

	cmd := command(ctx, inv)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start %s: %w", inv.Argv[0], err)
	}

	var (
		wg sync.WaitGroup
		// mu serializes emit and guards tail.
		//
		// stdout and stderr are pumped by two goroutines, so without this a caller's
		// callback would run concurrently and a caller that appends to a slice --
		// the obvious thing to write -- would lose lines to a data race.
		mu   sync.Mutex
		tail []string
	)
	// record keeps a few stderr lines so a failure can quote the CLI's own
	// explanation instead of reporting a bare exit status, and one holds both it
	// and emit under the same lock.
	record := func(l Line) {
		if l.IsStderr() && len(tail) < 5 {
			tail = append(tail, strings.TrimSpace(l.Text))
		}
		emit(l)
	}
	one := func(l Line) {
		mu.Lock()
		defer mu.Unlock()
		record(l)
	}

	wg.Add(2)
	go func() { defer wg.Done(); pump(stdout, "stdout", one) }()
	go func() { defer wg.Done(); pump(stderr, "stderr", one) }()
	wg.Wait()

	waitErr := cmd.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timed out after %s", r.runTimeout())
	}
	if waitErr != nil {
		mu.Lock()
		summary := firstNonEmpty(tail)
		mu.Unlock()
		return exitError(inv, waitErr, summary)
	}
	return nil
}

// command builds a child that can be cancelled as a whole.
func command(ctx context.Context, inv cli.Invocation) *exec.Cmd {
	cmd := exec.CommandContext(ctx, inv.Argv[0], inv.Argv[1:]...)
	cmd.Env = inv.Env
	cmd.Dir = inv.Dir

	// CommandContext alone kills only the direct child. The CLI spawns a
	// langgraph server subprocess, and a surviving grandchild would hold the
	// stdout pipe open, so the reader would block past the deadline and the UI
	// would appear frozen. Signalling the process group closes those pipes
	// immediately.
	setProcessGroup(cmd)
	cmd.Cancel = func() error { return killProcessGroup(cmd) }
	// Belt and braces: Wait stops waiting on the pipes after this delay even if
	// no signal could be delivered, so the timeout is always honoured.
	cmd.WaitDelay = 3 * time.Second
	return cmd
}

// describeFailure turns an exec error into something a user can act on.
func describeFailure(cmd *exec.Cmd, err error, output string) error {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if _, lookErr := exec.LookPath(cmd.Path); lookErr != nil {
			return fmt.Errorf("%s is not on PATH: %w", cmd.Path, lookErr)
		}
		if line := firstMeaningfulLine(output); line != "" {
			return errors.New(line)
		}
		return fmt.Errorf("%s exited with status %d", cmd.Path, exit.ExitCode())
	}
	return err
}

// exitError quotes the child's own message where it has one.
func exitError(inv cli.Invocation, err error, summary string) error {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		message := fmt.Sprintf("%s exited with status %d", strings.Join(inv.Argv, " "), exit.ExitCode())
		if summary != "" {
			message += ": " + summary
		}
		return errors.New(message)
	}
	return err
}

// errorLinePattern matches a real exception header ("ValueError: ...") and
// deliberately not a traceback source line such as "raise ValueError(msg)",
// which is a code frame rather than the failure.
var errorLinePattern = func(s string) bool {
	for _, marker := range []string{"Error:", "Exception:", "error:"} {
		if strings.Contains(s, marker) {
			return true
		}
	}
	return false
}

// firstMeaningfulLine picks the most useful line of a failed command's output.
//
// The CLI is Python, so a failure normally ends with an exception header, which
// explains far more than "exit status 1". Rich-formatted output is box-drawn, so
// borders are stripped first.
func firstMeaningfulLine(output string) string {
	fallback := ""
	for _, raw := range strings.Split(strings.TrimSpace(output), "\n") {
		line := strings.TrimSpace(stripANSI(raw))
		line = strings.Trim(line, "│┃ ")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if errorLinePattern(line) {
			return line
		}
		if fallback == "" {
			fallback = line
		}
	}
	return fallback
}

// firstNonEmpty returns the first line that carries content, so a leading blank
// or whitespace-only line does not hide the real message.
func firstNonEmpty(lines []string) string {
	for _, l := range lines {
		if trimmed := strings.TrimSpace(l); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// pump forwards a pipe line by line.
//
// The scanner buffer is enlarged because agent output can contain very long
// lines, and the default 64KiB limit would silently stop the stream.
func pump(reader io.Reader, stream string, emit func(Line)) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		emit(Line{Stream: stream, Text: scanner.Text()})
	}
}

// stripANSI removes SGR escape sequences so Rich colouring does not leak into a
// line this front-end renders itself.
func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && !(s[j] >= 0x40 && s[j] <= 0x7e) {
				j++
			}
			i = j + 1
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
