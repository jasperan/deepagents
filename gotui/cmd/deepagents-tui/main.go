// Command deepagents-tui is an ADDITIONAL way to run deepagents: a Go front-end
// in the charm v2 + huh stack that drives the same `deepagents` CLI the Python
// user runs.
//
// It never reimplements agent execution, thread storage or tool policy. The
// thread browser decodes the CLI's own `threads list --json` envelope, deletion
// shells out to `deepagents threads delete`, and a run shells out to
// `deepagents -n <task> -q`. A Go user and a Python user therefore get identical
// results from identical inputs.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/jasperan/deepagents/gotui/internal/cli"
	"github.com/jasperan/deepagents/gotui/internal/huhstyle"
	"github.com/jasperan/deepagents/gotui/internal/threads"
	"github.com/jasperan/deepagents/gotui/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "deepagents-tui: "+err.Error())
		os.Exit(1)
	}
}

func run() error {
	var (
		binary = flag.String("binary", "", "path to the deepagents CLI (default: "+cli.EnvBinary+", PATH, then a checkout)")

		threadsFlag = flag.Bool("threads", false, "list threads, then exit")
		agentsFlag  = flag.Bool("agents", false, "list available agents, then exit")
		deleteFlag  = flag.String("delete", "", "delete a thread by id (requires --yes)")
		yesFlag     = flag.Bool("yes", false, "confirm a destructive --delete without prompting")
		runFlag     = flag.String("run", "", "run one agent task, then exit")

		limitFlag   = flag.Int("limit", 0, "cap the thread listing (default: the CLI's own)")
		agentFlag   = flag.String("agent", "", "filter threads, or select the agent to run as")
		sortFlag    = flag.String("sort", "", "sort threads by created or updated (default: the CLI's config)")
		branchFlag  = flag.String("branch", "", "filter threads by git branch")
		verboseFlag = flag.Bool("verbose", false, "show every thread column")

		modelFlag   = flag.String("model", "", "model to run (passed as -M/--model)")
		approveFlag = flag.Bool("auto-approve", false,
			"pass -y, letting the agent run tools without asking")
		shellFlag = flag.String("shell-allow-list", "",
			"passed as -S: \"recommended\", \"all\", or a comma-separated command list")

		jsonFlag    = flag.Bool("json", false, "emit the CLI's own JSON for scripted actions")
		noInputFlag = flag.Bool("no-input", false, "never prompt; fail instead if input is required")
	)
	flag.Parse()

	launcher, err := cli.Resolve(*binary)
	if err != nil {
		return err
	}
	runner := tui.NewShellRunner()
	ctx := context.Background()

	// --- scripted path ---------------------------------------------------------
	// This runs before any prompt is considered, so a pipeline never blocks on a
	// question.
	action, scripted := tui.ParseActionFlags(*runFlag, *deleteFlag, *threadsFlag, *agentsFlag)
	if scripted {
		return tui.RunAction(ctx, runner, launcher, tui.ActionRequest{
			Action:         action,
			ThreadID:       *deleteFlag,
			Task:           *runFlag,
			Agent:          *agentFlag,
			Model:          *modelFlag,
			Limit:          *limitFlag,
			Sort:           *sortFlag,
			Branch:         *branchFlag,
			Verbose:        *verboseFlag,
			Confirmed:      *yesFlag,
			AutoApprove:    *approveFlag,
			ShellAllowList: *shellFlag,
			JSON:           *jsonFlag,
		}, os.Stdout)
	}

	// --- screen-reader / piped path --------------------------------------------
	// huh's accessible rendering only exists in its standalone Run path, so the
	// embedded full-screen UI is skipped entirely here.
	if huhstyle.Accessible() {
		fmt.Fprintln(os.Stdout, accessibleNotice)
		return runPlain(ctx, runner, launcher, *modelFlag, *noInputFlag)
	}
	if !huhstyle.Interactive() || *noInputFlag {
		return errors.New("no terminal on stdin: pass an action flag such as --threads, " +
			"--agents, --run \"<task>\" or --delete <id> --yes (or set ACCESSIBLE for plain prompts)")
	}

	// --- full-screen path -------------------------------------------------------
	model := tui.New(tui.Options{
		Launcher: launcher,
		Runner:   runner,
		Limit:    *limitFlag,
		Agent:    *agentFlag,
		Sort:     *sortFlag,
		Branch:   *branchFlag,
		Model:    *modelFlag,
	})
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

// accessibleNotice explains what the plain path is about to ask for.
const accessibleNotice = `deepagents-tui: screen-reader mode.

This will ask for the task and the run options, then run the deepagents CLI
exactly as the Python CLI would. Cancel with ctrl+c.`

// runPlain drives the same run form as plain prompts.
func runPlain(ctx context.Context, runner tui.Runner, launcher cli.Launcher, model string, noInput bool) error {
	if noInput {
		return errors.New("-no-input cannot be combined with ACCESSIBLE plain prompts")
	}

	// The agent list is a convenience: if it cannot be fetched, the form still
	// works because omitting -a uses the CLI's own default agent.
	var agents []threads.AgentInfo
	if raw, err := runner.Capture(ctx, cli.AgentsList(launcher)); err == nil {
		if parsed, parseErr := threads.ParseAgents([]byte(raw)); parseErr == nil {
			agents = parsed
		}
	}

	answers, err := tui.RunPlainRunPrompts(tui.PlainOptions{
		Launcher: launcher,
		Agents:   agents,
		Model:    model,
		Input:    os.Stdin,
		Output:   os.Stdout,
	})
	if err != nil {
		return err
	}
	if strings.TrimSpace(answers.Task) == "" {
		return errors.New("no task was given")
	}

	inv, err := cli.Run(launcher, cli.RunOptions{
		Task:           answers.Task,
		Agent:          answers.Agent,
		Model:          answers.Model,
		Quiet:          true,
		AutoApprove:    answers.AutoApprove,
		ShellAllowList: answers.ShellAllow,
	})
	if err != nil {
		return err
	}
	return runner.Stream(ctx, inv, func(l tui.Line) {
		fmt.Fprintln(os.Stdout, l.Text)
	})
}
