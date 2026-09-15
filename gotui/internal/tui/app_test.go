package tui

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"

	"github.com/jasperan/deepagents/gotui/internal/cli"
	"github.com/jasperan/deepagents/gotui/internal/threads"
)

// fakeRunner stands in for the deepagents CLI.
//
// The CLI is not installed in this checkout, so every test drives the model
// through this seam. It records the invocations it was given, which is how the
// tests assert that the decision the UI made matches the documented command.
//
// The recorder is mutex-guarded because a run's goroutine reaches Stream while
// the test is still reading: startRun launches that goroutine before returning
// the command that reads from it.
type fakeRunner struct {
	mu        sync.Mutex
	captures  []cli.Invocation
	streams   []cli.Invocation
	captureFn func(inv cli.Invocation) (string, error)
	streamFn  func(inv cli.Invocation, emit func(Line)) error
}

func (f *fakeRunner) Capture(_ context.Context, inv cli.Invocation) (string, error) {
	f.mu.Lock()
	f.captures = append(f.captures, inv)
	fn := f.captureFn
	f.mu.Unlock()
	if fn != nil {
		return fn(inv)
	}
	return "", nil
}

func (f *fakeRunner) Stream(_ context.Context, inv cli.Invocation, emit func(Line)) error {
	f.mu.Lock()
	f.streams = append(f.streams, inv)
	fn := f.streamFn
	f.mu.Unlock()
	if fn != nil {
		return fn(inv, emit)
	}
	return nil
}

// setCapture installs the reply function. A test may swap it between phases.
func (f *fakeRunner) setCapture(fn func(inv cli.Invocation) (string, error)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.captureFn = fn
}

func (f *fakeRunner) captureCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.captures)
}

func (f *fakeRunner) streamCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.streams)
}

// captureLines returns the rendered argv of the captures taken since index.
func (f *fakeRunner) captureLines(since int) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	if since > len(f.captures) {
		return nil
	}
	lines := make([]string, 0, len(f.captures)-since)
	for _, inv := range f.captures[since:] {
		lines = append(lines, inv.CommandLine())
	}
	return lines
}

// lastCapture returns the most recent invocation, or fails the test.
func (f *fakeRunner) lastCapture(t *testing.T) cli.Invocation {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.captures) == 0 {
		t.Fatal("the runner was never asked to capture anything")
	}
	return f.captures[len(f.captures)-1]
}

// streamAt returns a recorded stream invocation.
func (f *fakeRunner) streamAt(t *testing.T, index int) cli.Invocation {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if index >= len(f.streams) {
		t.Fatalf("the runner saw %d streams, want at least %d", len(f.streams), index+1)
	}
	return f.streams[index]
}

const sampleListJSON = `{"schema_version": 1, "command": "threads list", "data": [
  {"thread_id": "t-1", "agent_name": "coder", "updated_at": "2026-09-14T10:04:05Z",
   "git_branch": "main", "initial_prompt": "Summarize the repo", "message_count": 4, "cwd": "/repo"},
  {"thread_id": "t-2", "agent_name": "researcher", "updated_at": "2026-09-14T09:00:00Z",
   "initial_prompt": "Find TODOs", "message_count": 9}
]}`

const sampleAgentsJSON = `{"schema_version": 1, "command": "list", "data": [
  {"name": "coder", "path": "/a", "has_agents_md": true, "is_default": true},
  {"name": "researcher", "path": "/b", "has_agents_md": false, "is_default": false}
]}`

func newModel(runner Runner) Model {
	return New(Options{
		Launcher: cli.Launcher{Bin: "deepagents"},
		Runner:   runner,
		Width:    100,
		Height:   30,
	})
}

// loadModel performs the startup a real program does: Init, then drain.
func loadModel(t *testing.T, m Model) Model {
	t.Helper()
	loaded, ok := Drain(m, m.Init(), 0).(Model)
	if !ok {
		t.Fatal("Drain did not return a Model")
	}
	return loaded
}

func pressKey(t *testing.T, m tea.Model, key string) tea.Model {
	t.Helper()
	msg := tea.KeyPressMsg{Code: rune(key[0]), Text: key}
	if len(key) > 1 {
		// Named keys arrive as their own code with no text.
		msg = tea.KeyPressMsg{Code: keyCode(key)}
	}
	updated, cmd := m.Update(msg)
	return Drain(updated, cmd, 0)
}

// keyCode maps the few named keys the tests use onto their constants.
func keyCode(name string) rune {
	switch name {
	case "enter":
		return tea.KeyEnter
	case "esc":
		return tea.KeyEscape
	default:
		return rune(name[0])
	}
}

// --- browser ----------------------------------------------------------------

func TestStartupListsThreadsAndRendersThem(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	if len(m.all) != 2 || len(m.shown) != 2 {
		t.Fatalf("loaded %d threads, want 2", len(m.all))
	}
	view := m.View().Content
	for _, want := range []string{"t-1", "coder", "Summarize the repo", "2 thread(s)"} {
		if !strings.Contains(view, want) {
			t.Errorf("view is missing %q", want)
		}
	}
	if !strings.Contains(fr.lastCapture(t).CommandLine(), "--json") {
		t.Error("the thread listing was not requested as JSON")
	}
}

// A fresh install has no threads; that must read as an empty state with a next
// step, not as a failure.
func TestEmptyListingShowsAnEmptyStateNotAnError(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 1, "command": "threads list", "data": []}`, nil
	}}
	m := loadModel(t, newModel(fr))

	view := m.View().Content
	if !strings.Contains(view, "No threads are stored yet") {
		t.Errorf("view does not explain the empty state: %q", view)
	}
	if m.statusKind == "error" {
		t.Errorf("an empty listing was reported as an error: %q", m.status)
	}
}

// A failed refresh must not wipe the rows already on screen.
func TestFailedRefreshKeepsTheExistingRows(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	fr.setCapture(func(cli.Invocation) (string, error) { return "", errors.New("CLI vanished") })
	updated, _ := m.Update(threadsMsg{err: errors.New("CLI vanished")})
	final := updated.(Model)

	if len(final.shown) != 2 {
		t.Errorf("rows were cleared by a failed refresh: %d left", len(final.shown))
	}
	if final.statusKind != "error" {
		t.Errorf("statusKind = %q, want error", final.statusKind)
	}
}

func TestUnknownSchemaIsReportedAsAnError(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 99, "command": "threads list", "data": []}`, nil
	}}
	m := loadModel(t, newModel(fr))
	if m.statusKind != "error" {
		t.Errorf("statusKind = %q, want error for an unknown schema", m.statusKind)
	}
	if !strings.Contains(m.status, "schema_version") {
		t.Errorf("status %q does not explain the mismatch", m.status)
	}
}

// --- key handling -----------------------------------------------------------

// Regression: a sibling repository shipped v2 handlers that fired twice per
// keystroke because the same key was interpreted both as tea.KeyMsg and as
// tea.KeyPressMsg. One press must produce exactly one action.
func TestOneKeyPressFiresExactlyOneAction(t *testing.T) {
	fr := &fakeRunner{captureFn: func(inv cli.Invocation) (string, error) {
		if strings.Contains(inv.CommandLine(), " list ") || strings.HasSuffix(inv.CommandLine(), "list --json") {
			return sampleAgentsJSON, nil
		}
		return sampleListJSON, nil
	}}
	m := loadModel(t, newModel(fr))
	before := fr.captureCount()

	// "n" opens the run form, which requires one agent listing.
	after := pressKey(t, m, "n").(Model)

	agentsCalls := 0
	for _, inv := range fr.captures[before:] {
		if strings.Contains(inv.CommandLine(), "list --json") && !strings.Contains(inv.CommandLine(), "threads") {
			agentsCalls++
		}
	}
	if agentsCalls != 1 {
		t.Errorf("one press of \u201cn\u201d produced %d agent listings, want 1", agentsCalls)
	}
	if after.runForm == nil {
		t.Error("the run form did not open")
	}
}

// Regression: a terminal can be resized to a few columns, and every width in
// the render path is derived by subtracting chrome from it. A panic below 6
// columns broke a sibling repository's TUI.
func TestRenderingIsSafeAtEveryTerminalWidth(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	for _, width := range []int{-5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 20, 40, 80, 200} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("width %d panicked: %v", width, r)
				}
			}()
			updated, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 3})
			final := updated.(Model)
			if final.View().Content == "" && width > 0 {
				t.Errorf("width %d rendered nothing", width)
			}
			// The filter field, detail pane and output pane must survive too.
			final.detail.SetContent(final.detailText(threads.ThreadInfo{ThreadID: "t"}))
			_ = final.detail.View()
			final.appendLine(Line{Text: "a very long line of output that will not fit"})
			_ = final.output.View()
			_ = final.render()
		}()
	}
}

func TestTruncateAndPadAreSafeAtEveryWidth(t *testing.T) {
	for _, n := range []int{-3, 0, 1, 2, 5} {
		if got := Truncate("abcdef", n); len([]rune(got)) > maxInt(0, n) {
			t.Errorf("Truncate(_, %d) = %q, which is longer than requested", n, got)
		}
		_ = Pad("abc", n)
	}
}

func TestShortKeysMapToTheDocumentedActions(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))
	before := fr.captureCount()

	// r refreshes.
	m = pressKey(t, m, "r").(Model)
	if fr.captureCount() != before+1 {
		t.Errorf("r did not trigger a refresh: %d captures", fr.captureCount())
	}

	// / focuses the filter, and typing narrows the list.
	m = pressKey(t, m, "/").(Model)
	if !m.filtering {
		t.Fatal("/ did not focus the filter")
	}
	for _, r := range "researcher" {
		updated, cmd := m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
		m = Drain(updated, cmd, 0).(Model)
	}
	if len(m.shown) != 1 || m.shown[0].ThreadID != "t-2" {
		t.Errorf("filtering to researcher showed %d rows", len(m.shown))
	}

	// esc clears it again.
	m = pressKey(t, m, "esc").(Model)
	if m.filtering || len(m.shown) != 2 {
		t.Errorf("esc left filtering=%v with %d rows", m.filtering, len(m.shown))
	}
}

// --- detail -----------------------------------------------------------------

func TestEnterOpensTheDetailPaneWithTheSelectedThread(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	m = pressKey(t, m, "enter").(Model)
	if m.screen != screenDetail {
		t.Fatalf("screen = %v, want screenDetail", m.screen)
	}
	if m.target.ThreadID != "t-1" {
		t.Errorf("target = %q, want the row under the cursor", m.target.ThreadID)
	}
	view := m.View().Content
	for _, want := range []string{"Thread detail", "t-1", "Summarize the repo", "main"} {
		if !strings.Contains(view, want) {
			t.Errorf("detail view is missing %q", want)
		}
	}

	// q returns to the browser.
	m = pressKey(t, m, "q").(Model)
	if m.screen != screenThreads {
		t.Errorf("screen = %v, want screenThreads after q", m.screen)
	}
}

// A filter that narrows the list must drag the cursor back into range, or the
// next action would target a row that is no longer displayed.
func TestCursorStaysInRangeWhenTheFilterNarrowsTheList(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	m.table.SetCursor(1)
	m.filter.SetValue("Summarize")
	m.applyFilter()

	if m.table.Cursor() >= len(m.shown) && len(m.shown) > 0 {
		t.Fatalf("cursor %d is past the end of %d rows", m.table.Cursor(), len(m.shown))
	}
	selected, ok := m.selected()
	if !ok || selected.ThreadID != "t-1" {
		t.Errorf("selected = %+v (ok=%v), want the only visible row", selected, ok)
	}
}

// --- deletion ---------------------------------------------------------------

// The CLI deletes without asking, so the confirmation lives in this front-end.
// The deletion must not run until it is confirmed, and when it does run it must
// name the selected thread.
func TestDeleteConfirmationGatesTheDeletionCommand(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	m = pressKey(t, m, "d").(Model)
	if m.screen != screenDelete || m.deleteForm == nil {
		t.Fatalf("d did not open the confirmation (screen=%v)", m.screen)
	}
	if fr.captureCount() != 1 {
		t.Fatalf("opening the confirmation ran %d commands; it must not delete anything yet", fr.captureCount())
	}

	// Cancelling must not delete.
	cancelled, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if fr.captureCount() != 1 {
		t.Error("cancelling the confirmation still ran a command")
	}
	if cancelled.(Model).screen != screenThreads {
		t.Error("cancelling did not return to the browser")
	}
}

func TestDeleteCommandNamesTheThreadAndAsksForJSON(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))
	m.target = m.shown[0]

	// The seam: the form cannot be driven headlessly, but the command it would
	// produce can be executed directly.
	result := Drain(m, m.deleteThread(m.target.ThreadID), 0).(Model)
	got := fr.lastCapture(t).CommandLine()
	if !strings.Contains(got, "threads delete t-1") || !strings.Contains(got, "--json") {
		t.Errorf("argv = %q, want `threads delete t-1 --json`", got)
	}
	if result.screen != screenThreads {
		t.Errorf("screen = %v, want the browser after a delete", result.screen)
	}
}

func TestDeleteResultIsReported(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	updated, _ := m.Update(deleteMsg{result: threads.DeleteResult{ThreadID: "t-1", Deleted: true}})
	final := updated.(Model)
	if final.statusKind != "success" || !strings.Contains(final.status, "t-1") {
		t.Errorf("status = %q (%s), want a success naming the thread", final.status, final.statusKind)
	}

	// A thread the CLI could not find is reported, not treated as a failure.
	updated, _ = m.Update(deleteMsg{result: threads.DeleteResult{ThreadID: "t-9"}})
	final = updated.(Model)
	if final.statusKind != "warning" {
		t.Errorf("statusKind = %q, want warning for a not-found thread", final.statusKind)
	}
}

// --- run --------------------------------------------------------------------

// waitForStreams polls until the runner has seen n streams.
//
// startRun launches its goroutine before returning the command that reads from
// it, so the invocation reaches the runner asynchronously.
func waitForStreams(t *testing.T, fr *fakeRunner, n int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if fr.streamCount() >= n {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("the runner saw %d streams, want at least %d", fr.streamCount(), n)
}

// The forms cannot be driven headlessly, so startRun is where the
// answers-become-an-invocation decision is pinned.
func TestStartRunBuildsTheDocumentedCommandFromTheAnswers(t *testing.T) {
	fr := &fakeRunner{}
	m := newModel(fr)
	m.runAnswers = &runAnswers{
		Task:        "summarize the repo",
		Agent:       "coder",
		Model:       "claude-sonnet-4-6",
		AutoApprove: true,
		ShellAllow:  "recommended",
	}

	cmd := m.startRun()
	if m.screen != screenRunning {
		t.Errorf("screen = %v, want screenRunning", m.screen)
	}
	for _, want := range []string{"-n summarize the repo", "-q", "-a coder", "-M claude-sonnet-4-6", "-y", "-S recommended"} {
		if !strings.Contains(m.invocation, want) {
			t.Errorf("invocation %q is missing %q", m.invocation, want)
		}
	}
	if cmd == nil {
		t.Fatal("startRun returned no command, so nothing would be read")
	}

	// The same argv must be what actually reaches the runner.
	waitForStreams(t, fr, 1)
	if got := fr.streamAt(t, 0).CommandLine(); got != m.invocation {
		t.Errorf("runner argv = %q, want the displayed %q", got, m.invocation)
	}
}

// An empty task must be refused before anything is spawned.
func TestStartRunRefusesAnEmptyTask(t *testing.T) {
	fr := &fakeRunner{}
	m := newModel(fr)
	m.runAnswers = &runAnswers{Task: "   "}

	m.startRun()
	time.Sleep(10 * time.Millisecond)
	if fr.streamCount() != 0 {
		t.Error("an empty task still spawned the CLI")
	}
	if m.screen != screenOutput || m.runErr == nil {
		t.Fatalf("screen=%v err=%v, want an error state", m.screen, m.runErr)
	}
	if !strings.Contains(m.View().Content, "task is required") {
		t.Errorf("view does not explain the failure: %q", m.View().Content)
	}
}

// step feeds one message to the model and returns the updated model.
func step(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	updated, _ := m.Update(msg)
	final, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", updated)
	}
	return final
}

func TestStreamedRunOutputLandsInThePaneAndFinishes(t *testing.T) {
	fr := &fakeRunner{}
	m := newModel(fr)
	m.runAnswers = &runAnswers{Task: "t"}
	m.startRun()

	m = step(t, m, streamLineMsg(Line{Text: "line one"}))
	m = step(t, m, streamLineMsg(Line{Stream: "stderr", Text: "warning"}))

	if len(m.lines) != 3 {
		t.Fatalf("collected %d lines, want the invocation plus 2", len(m.lines))
	}
	if !strings.Contains(m.View().Content, "line one") {
		t.Error("streamed output is not in the view")
	}

	m = step(t, m, streamDoneMsg{})
	if m.screen != screenOutput {
		t.Errorf("screen = %v, want screenOutput", m.screen)
	}
	if m.statusKind != "success" {
		t.Errorf("statusKind = %q, want success", m.statusKind)
	}
}

// A failed run is the expected offline outcome and must surface as an error
// state with the CLI's own message, never as a hang.
func TestFailedRunSurfacesTheError(t *testing.T) {
	fr := &fakeRunner{}
	m := newModel(fr)
	m.runAnswers = &runAnswers{Task: "t"}
	m.startRun()

	m = step(t, m, streamDoneMsg{err: errors.New("deepagents exited with status 1: model unavailable")})
	if m.runErr == nil || m.statusKind != "error" {
		t.Fatalf("err=%v kind=%q, want an error state", m.runErr, m.statusKind)
	}
	view := m.View().Content
	if !strings.Contains(view, "Run failed") || !strings.Contains(view, "model unavailable") {
		t.Errorf("view does not show the failure: %q", view)
	}
}

// --- forms ------------------------------------------------------------------

func TestRunFormAsksForTheTaskAndTheOptions(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleAgentsJSON, nil }}
	m := loadModel(t, newModel(fr))
	m = pressKey(t, m, "n").(Model)

	started, ok := Drain(m, m.runForm.Init(), 0).(Model)
	if !ok {
		t.Fatal("Drain did not return a Model")
	}
	first := started.View().Content
	for _, want := range []string{"Task"} {
		if !strings.Contains(first, want) {
			t.Errorf("the first page is missing the %q field", want)
		}
	}

	// The task field is required, so a blank answer must not complete the form.
	advanced, _ := started.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	after := Drain(advanced, nil, 0).(Model)
	if after.runForm == nil {
		t.Fatal("the form closed on a blank task")
	}
	if after.runForm.State == huh.StateCompleted {
		t.Error("the form completed with a blank task")
	}
	if strings.TrimSpace(after.runAnswers.Task) != "" {
		t.Errorf("the task was silently defaulted to %q", after.runAnswers.Task)
	}
}

// The agent list is a convenience: a query that fails must still leave a usable
// form, because omitting -a uses the CLI's own default agent.
func TestRunFormStillOpensWhenTheAgentListFails(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return "", errors.New("no agents directory")
	}}
	m := loadModel(t, newModel(fr))
	m = pressKey(t, m, "n").(Model)

	if m.runForm == nil {
		t.Fatal("the run form did not open after a failed agent listing")
	}
	if m.statusKind != "warning" {
		t.Errorf("statusKind = %q, want a warning", m.statusKind)
	}
}

// --- infrastructure ---------------------------------------------------------

// Regression: huh re-arms the text input's blink on every update, so a drain
// that feeds BlinkMsg back into Update never terminates. This hung a sibling
// repository's suite for minutes.
func TestDrainTerminatesOnABlinkLoop(t *testing.T) {
	m := newModel(&fakeRunner{})
	blink := func() tea.Msg { return cursor.BlinkMsg{} }

	done := make(chan tea.Model, 1)
	go func() { done <- Drain(m, blink, 0) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Drain did not terminate on a cursor.BlinkMsg stream")
	}
}

func TestCtrlCQuitsFromEveryScreen(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	screens := []Model{m}
	detail := pressKey(t, m, "enter").(Model)
	screens = append(screens, detail)
	runForm := pressKey(t, m, "n").(Model)
	screens = append(screens, runForm)
	del := pressKey(t, m, "d").(Model)
	screens = append(screens, del)

	for _, s := range screens {
		_, cmd := s.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
		if cmd == nil {
			t.Errorf("screen %v did not quit on ctrl+c", s.screen)
		}
	}
}

// --- escape from an open form -------------------------------------------------

// huh's default keymap binds Quit to ctrl+c ONLY, and StateAborted is reachable
// only through keymap.Quit (verified in charm.land/huh/v2@v2.0.3: keymap.go binds
// Quit to ctrl+c alone, and form.go is the single place that sets StateAborted).
// Escape therefore never aborts a huh form.
//
// The failure mode that makes this worth a test: a front-end that delegates "esc"
// to the form leaves it installed, and the orphaned form then swallows every later
// keystroke while the user believes they cancelled and the footer still reads
// "esc cancel". All three symptoms are asserted below, because the swallowed input
// is exactly what a screen-only assertion would miss.
func TestEscClosesTheRunForm(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	m = pressKey(t, m, "n").(Model)
	if m.runForm == nil || m.screen != screenRunForm {
		t.Fatalf("n did not open the run form (form=%v screen=%v)", m.runForm != nil, m.screen)
	}

	m = pressKey(t, m, "esc").(Model)

	if m.runForm != nil {
		t.Fatal("esc left the run form installed; it would keep swallowing keystrokes")
	}
	if m.screen != screenThreads {
		t.Errorf("esc left the screen at %v, want the thread browser", m.screen)
	}

	// The part that catches the real damage: an orphaned form eats every later
	// key. "/" only starts filtering if the browser is genuinely back in charge.
	m = pressKey(t, m, "/").(Model)
	if !m.filtering {
		t.Error("the browser refused a key after esc; a form is still swallowing input")
	}
}

// TestEscClosesTheDeleteForm covers the confirmation, where a surviving form is
// worst of all: it would keep asking for a destructive confirmation the user has
// already cancelled, and the very next keystroke would answer it.
func TestEscClosesTheDeleteForm(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	m := loadModel(t, newModel(fr))

	m = pressKey(t, m, "d").(Model)
	if m.deleteForm == nil || m.screen != screenDelete {
		t.Fatalf("d did not open the confirmation (form=%v screen=%v)", m.deleteForm != nil, m.screen)
	}
	afterOpen := fr.captureCount()

	m = pressKey(t, m, "esc").(Model)

	if m.deleteForm != nil {
		t.Fatal("esc left the confirmation installed")
	}
	if m.screen != screenThreads {
		t.Errorf("esc left the screen at %v, want the thread browser", m.screen)
	}
	if fr.captureCount() != afterOpen {
		t.Error("esc ran a command; cancelling must not delete anything")
	}

	// Keys must reach the browser again, and must not answer the stale form.
	m = pressKey(t, m, "/").(Model)
	if !m.filtering {
		t.Error("the browser refused a key after esc; the confirmation is still swallowing input")
	}
	if fr.captureCount() != afterOpen {
		t.Error("a key after esc ran a command; the confirmation is still live")
	}
}
