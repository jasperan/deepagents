package tui

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jasperan/deepagents/gotui/internal/cli"
)

func plainLauncher() cli.Launcher { return cli.Launcher{Bin: "deepagents"} }

// A scripted action must be chosen before any prompt is considered, so a
// pipeline never blocks on a question. The most destructive flag also has to
// win: "--threads --delete x" must not silently degrade into a read.
func TestParseActionFlagsPrefersTheMostSpecificAndMostDestructive(t *testing.T) {
	cases := []struct {
		task, threadID          string
		threadsFlag, agentsFlag bool
		want                    string
	}{
		{"summarize", "", false, false, ActionRun},
		{"", "abc", false, false, ActionDelete},
		{"", "", true, false, ActionThreads},
		{"", "", false, true, ActionAgents},
		// Specificity order: run beats delete beats threads beats agents.
		{"task", "abc", true, true, ActionRun},
		{"", "abc", true, true, ActionDelete},
		{"", "", true, true, ActionThreads},
	}
	for _, c := range cases {
		got, scripted := ParseActionFlags(c.task, c.threadID, c.threadsFlag, c.agentsFlag)
		if !scripted || got != c.want {
			t.Errorf("ParseActionFlags(%q, %q, %v, %v) = (%q, %v), want (%q, true)",
				c.task, c.threadID, c.threadsFlag, c.agentsFlag, got, scripted, c.want)
		}
	}
	if _, scripted := ParseActionFlags("", "", false, false); scripted {
		t.Error("no flags still reported a scripted action")
	}
}

// --json must forward the CLI's bytes untouched, because re-encoding a decoded
// value risks dropping a field the CLI added.
func TestThreadsActionForwardsJSONVerbatim(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	var out bytes.Buffer

	err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionThreads, JSON: true}, &out)
	if err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(sampleListJSON) {
		t.Errorf("JSON was not forwarded verbatim:\n%s", out.String())
	}
}

func TestThreadsActionRendersATableWithoutJSON(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleListJSON, nil }}
	var out bytes.Buffer

	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionThreads}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	text := out.String()
	for _, want := range []string{"2 thread(s)", "t-1", "coder", "Summarize the repo"} {
		if !strings.Contains(text, want) {
			t.Errorf("output is missing %q:\n%s", want, text)
		}
	}
}

func TestThreadsActionReportsAnEmptyListing(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 1, "command": "threads list", "data": []}`, nil
	}}
	var out bytes.Buffer

	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionThreads}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	if !strings.Contains(out.String(), "No threads are stored yet") {
		t.Errorf("output does not explain the empty state: %q", out.String())
	}
}

// A prompt is never the only route through a delete, and a pipe is never
// mistaken for consent: the deletion must refuse without an explicit --yes.
func TestDeleteRequiresAnExplicitConfirmation(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 1, "command": "threads delete", "data": {"thread_id": "abc", "deleted": true}}`, nil
	}}
	var out bytes.Buffer

	err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionDelete, ThreadID: "abc"}, &out)
	if err == nil {
		t.Fatal("delete ran without --yes")
	}
	if !strings.Contains(err.Error(), "--yes") {
		t.Errorf("error %q does not say how to confirm", err)
	}
	if fr.captureCount() != 0 {
		t.Error("the CLI was invoked despite the refusal")
	}
}

func TestDeleteWithoutAnIDIsRefused(t *testing.T) {
	fr := &fakeRunner{}
	var out bytes.Buffer
	err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionDelete, Confirmed: true}, &out)
	if err == nil || fr.captureCount() != 0 {
		t.Fatalf("err=%v captures=%d, want a refusal with no invocation", err, fr.captureCount())
	}
}

func TestConfirmedDeleteReportsBothOutcomes(t *testing.T) {
	// Found: a success line naming the thread.
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 1, "command": "threads delete", "data": {"thread_id": "abc", "deleted": true}}`, nil
	}}
	var out bytes.Buffer
	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionDelete, ThreadID: "abc", Confirmed: true}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	if !strings.Contains(out.String(), "deleted thread abc") {
		t.Errorf("output = %q", out.String())
	}
	if got := fr.lastCapture(t).CommandLine(); !strings.Contains(got, "threads delete abc") {
		t.Errorf("argv = %q", got)
	}

	// Not found: reported as a normal outcome, not as an error.
	fr = &fakeRunner{captureFn: func(cli.Invocation) (string, error) {
		return `{"schema_version": 1, "command": "threads delete", "data": {"thread_id": "abc", "deleted": false}}`, nil
	}}
	out.Reset()
	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionDelete, ThreadID: "abc", Confirmed: true}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	if !strings.Contains(out.String(), "was not found") {
		t.Errorf("output = %q", out.String())
	}
}

func TestRunActionStreamsTheRunToStdout(t *testing.T) {
	fr := &fakeRunner{streamFn: func(_ cli.Invocation, emit func(Line)) error {
		emit(Line{Text: "the answer"})
		emit(Line{Stream: "stderr", Text: "a warning"})
		return nil
	}}
	var out bytes.Buffer

	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionRun, Task: "summarize"}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	for _, want := range []string{"the answer", "a warning"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("streamed output is missing %q: %q", want, out.String())
		}
	}
	// The scripted path must always be quiet, or a pipe would carry the CLI's
	// progress chatter as well as the answer.
	if got := fr.streamAt(t, 0).CommandLine(); !strings.Contains(got, "-q") {
		t.Errorf("argv %q is not quiet", got)
	}
}

func TestRunActionPropagatesARunFailure(t *testing.T) {
	fr := &fakeRunner{streamFn: func(cli.Invocation, func(Line)) error {
		return errors.New("deepagents exited with status 1: model unavailable")
	}}
	var out bytes.Buffer
	err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionRun, Task: "t"}, &out)
	if err == nil {
		t.Fatal("a failed run reported success")
	}
}

func TestRunActionRefusesAnEmptyTask(t *testing.T) {
	fr := &fakeRunner{}
	var out bytes.Buffer
	err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionRun, Task: "  "}, &out)
	if err == nil || fr.streamCount() != 0 {
		t.Fatalf("err=%v streams=%d, want a refusal with no invocation", err, fr.streamCount())
	}
}

func TestAgentsActionListsNamesAndMarksTheDefault(t *testing.T) {
	fr := &fakeRunner{captureFn: func(cli.Invocation) (string, error) { return sampleAgentsJSON, nil }}
	var out bytes.Buffer

	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: ActionAgents}, &out); err != nil {
		t.Fatalf("RunAction: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "2 agent(s)") {
		t.Errorf("output = %q", text)
	}
	if !strings.Contains(text, "* coder") {
		t.Errorf("the default agent is not marked:\n%s", text)
	}
}

func TestUnknownActionIsRefused(t *testing.T) {
	fr := &fakeRunner{}
	var out bytes.Buffer
	if err := RunAction(context.Background(), fr, plainLauncher(),
		ActionRequest{Action: "nonsense"}, &out); err == nil {
		t.Fatal("an unknown action was accepted")
	}
}

// writeRaw must refuse anything that is not JSON, so a caller cannot mistake a
// human table for a machine contract.
func TestWriteRawRefusesNonJSON(t *testing.T) {
	var out bytes.Buffer
	if err := writeRaw(&out, "╭───╮\n│ x │\n"); err == nil {
		t.Error("writeRaw accepted a Rich table")
	}
	if err := writeRaw(&out, "   "); err == nil {
		t.Error("writeRaw accepted empty output")
	}
	if err := writeRaw(&out, `{"a": 1}`); err != nil {
		t.Errorf("writeRaw rejected valid JSON: %v", err)
	}
}
