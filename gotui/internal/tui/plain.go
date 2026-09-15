package tui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jasperan/deepagents/gotui/internal/cli"
	"github.com/jasperan/deepagents/gotui/internal/threads"
)

// Actions a scripted (non-interactive) invocation can request.
const (
	ActionThreads = "threads"
	ActionAgents  = "agents"
	ActionDelete  = "delete"
	ActionRun     = "run"
)

// ActionRequest is a scripted, non-interactive request.
//
// Every field is either a filter for a read or an option forwarded to the CLI,
// so a scripted run produces exactly the invocation the TUI would have built.
type ActionRequest struct {
	Action string
	// ThreadID is required by ActionDelete.
	ThreadID string
	// Task is required by ActionRun.
	Task  string
	Agent string
	Model string
	// Limit, Sort and Branch filter ActionThreads.
	Limit   int
	Sort    string
	Branch  string
	Verbose bool
	// Confirmed gates ActionDelete. Only --yes sets it: a prompt is never the
	// only route through a delete, and a pipe is never mistaken for consent.
	Confirmed bool
	// AutoApprove maps to -y for ActionRun.
	AutoApprove bool
	// ShellAllowList maps to -S for ActionRun.
	ShellAllowList string
	// JSON prints the CLI's own envelope verbatim instead of a rendered table.
	JSON bool
}

// RunAction executes a scripted request.
//
// It never reimplements the work: reads and deletes are delegated to the CLI and
// their own JSON envelopes are decoded, and a run is delegated and streamed
// through. A --json request forwards the CLI's bytes untouched so the output
// contract stays exactly the CLI's.
func RunAction(ctx context.Context, runner Runner, launcher cli.Launcher, req ActionRequest, out io.Writer) error {
	switch req.Action {
	case ActionThreads:
		return runThreadsAction(ctx, runner, launcher, req, out)

	case ActionAgents:
		return runAgentsAction(ctx, runner, launcher, req, out)

	case ActionDelete:
		if !req.Confirmed {
			return errors.New("refusing to delete without --yes: " +
				"deletion removes the conversation and its checkpoints")
		}
		if strings.TrimSpace(req.ThreadID) == "" {
			return errors.New("delete needs a thread id: --delete <thread_id> --yes")
		}
		raw, err := runner.Capture(ctx, cli.ThreadsDelete(launcher, req.ThreadID))
		if err != nil {
			return err
		}
		if req.JSON {
			return writeRaw(out, raw)
		}
		result, err := threads.ParseThreadsDelete([]byte(raw))
		if err != nil {
			return err
		}
		if !result.Deleted {
			fmt.Fprintf(out, "thread %s was not found\n", result.ThreadID)
			return nil
		}
		fmt.Fprintf(out, "deleted thread %s\n", result.ThreadID)
		return nil

	case ActionRun:
		if strings.TrimSpace(req.Task) == "" {
			return errors.New("run needs a task: --run \"<task>\"")
		}
		inv, err := cli.Run(launcher, cli.RunOptions{
			Task:           req.Task,
			Agent:          req.Agent,
			Model:          req.Model,
			Quiet:          true,
			AutoApprove:    req.AutoApprove,
			ShellAllowList: req.ShellAllowList,
		})
		if err != nil {
			return err
		}
		// Streaming straight to the writer keeps a non-interactive run's output
		// identical to running the CLI directly.
		return runner.Stream(ctx, inv, func(l Line) {
			fmt.Fprintln(out, l.Text)
		})
	}
	return fmt.Errorf("unknown action %q", req.Action)
}

// runThreadsAction lists threads, forwarding the CLI's JSON when asked.
func runThreadsAction(ctx context.Context, runner Runner, launcher cli.Launcher, req ActionRequest, out io.Writer) error {
	raw, err := runner.Capture(ctx, cli.ThreadsList(launcher, cli.ThreadsListOptions{
		Agent:   req.Agent,
		Limit:   req.Limit,
		Sort:    req.Sort,
		Branch:  req.Branch,
		Verbose: req.Verbose,
	}))
	if err != nil {
		return err
	}
	if req.JSON {
		return writeRaw(out, raw)
	}
	list, err := threads.ParseThreadsList([]byte(raw))
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(out, "No threads are stored yet.")
		return nil
	}
	fmt.Fprintf(out, "%d thread(s)\n", len(list))
	fmt.Fprintf(out, "  %-*s  %-12s  %5s  %-16s  %s\n",
		colIDWidth, "THREAD", "AGENT", "MSGS", "UPDATED", "PROMPT")
	for _, t := range list {
		fmt.Fprintf(out, "  %-*s  %-12s  %5s  %-16s  %s\n",
			colIDWidth, Truncate(t.ThreadID, colIDWidth),
			Truncate(t.AgentLabel(), 12),
			t.Messages(),
			t.Updated(),
			Truncate(t.Prompt(), 60))
	}
	return nil
}

// runAgentsAction lists the available agents.
func runAgentsAction(ctx context.Context, runner Runner, launcher cli.Launcher, req ActionRequest, out io.Writer) error {
	raw, err := runner.Capture(ctx, cli.AgentsList(launcher))
	if err != nil {
		return err
	}
	if req.JSON {
		return writeRaw(out, raw)
	}
	list, err := threads.ParseAgents([]byte(raw))
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(out, "No agents are stored yet.")
		return nil
	}
	fmt.Fprintf(out, "%d agent(s)\n", len(list))
	for _, a := range list {
		marker := "  "
		if a.IsDefault {
			marker = "* "
		}
		fmt.Fprintf(out, "  %s%s\n", marker, a.Name)
	}
	return nil
}

// writeRaw forwards the CLI's own output untouched.
//
// Re-encoding a decoded value would risk dropping a field the CLI added, so the
// envelope is passed through byte for byte.
func writeRaw(out io.Writer, raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return errors.New("the deepagents CLI produced no output")
	}
	if !json.Valid([]byte(trimmed)) {
		return errors.New("the CLI did not produce JSON; is --json supported by this version?")
	}
	_, err := fmt.Fprintln(out, trimmed)
	return err
}

// PlainOptions drive the screen-reader (accessible) prompt path.
type PlainOptions struct {
	Launcher cli.Launcher
	Agents   []threads.AgentInfo
	Model    string
	Input    io.Reader
	Output   io.Writer
}

// RunPlainRunPrompts runs the run form as standalone plain prompts.
//
// This path exists because huh's accessible rendering is only implemented in
// Form.Run: when a form is embedded in a Bubble Tea program the flag is never
// consulted, so an accessible user would get the full-screen UI regardless. All
// fields therefore live in ONE form, read once: huh wraps the reader in its own
// scanner, which buffers ahead, so a second form over the same stdin would start
// at EOF.
func RunPlainRunPrompts(opts PlainOptions) (RunAnswers, error) {
	in, out := opts.Input, opts.Output
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}

	names := threads.AgentNames(opts.Agents)
	defaultAgent := threads.DefaultAgentName(opts.Agents)
	if len(names) == 0 {
		names = []string{""}
	}
	if defaultAgent == "" {
		defaultAgent = names[0]
	}

	form, answers := runForm(names, defaultAgent, opts.Model)
	form = form.WithInput(in).WithOutput(out)
	if err := form.Run(); err != nil {
		return RunAnswers{}, err
	}
	return RunAnswers{
		Task:        answers.Task,
		Agent:       answers.Agent,
		Model:       answers.Model,
		AutoApprove: answers.AutoApprove,
		ShellAllow:  answers.ShellAllow,
	}, nil
}

// RunAnswers are the plain-path answers, exported so main can act on them.
type RunAnswers struct {
	Task        string
	Agent       string
	Model       string
	AutoApprove bool
	ShellAllow  string
}

// ParseActionFlags decides which scripted action a flag set requests.
//
// Order matters: the most specific and most destructive flag wins, so
// "--threads --delete x" cannot silently degrade into a read.
func ParseActionFlags(task, threadID string, threadsFlag, agentsFlag bool) (string, bool) {
	switch {
	case strings.TrimSpace(task) != "":
		return ActionRun, true
	case strings.TrimSpace(threadID) != "":
		return ActionDelete, true
	case threadsFlag:
		return ActionThreads, true
	case agentsFlag:
		return ActionAgents, true
	}
	return "", false
}
