package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"

	"github.com/jasperan/deepagents/gotui/internal/huhstyle"
	"github.com/jasperan/deepagents/gotui/internal/threads"
)

// runAnswers holds every answer the run form collects.
//
// It is heap-allocated and shared by pointer because huh writes through the
// pointer it is given. Binding a field to a member of a value-typed model would
// silently persist nothing, since Update receives a copy of the model.
type runAnswers struct {
	Task        string
	Agent       string
	Model       string
	AutoApprove bool
	ShellAllow  string
}

// deleteAnswers holds the deletion confirmation.
type deleteAnswers struct {
	Confirmed bool
}

// ShellAllowDisabled is the value meaning "leave shell commands disabled", which
// is the CLI's own default when -S is not passed.
const ShellAllowDisabled = ""

// themed attaches the project theme, accessibility flag and keymap.
//
// The keymap matters: a bare field has no keymap, so its keys would be silently
// ignored and the form would appear frozen.
func themed(f *huh.Form) *huh.Form {
	return f.
		WithTheme(huh.ThemeFunc(huhstyle.Theme)).
		WithAccessible(huhstyle.Accessible()).
		WithKeyMap(huh.NewDefaultKeyMap())
}

// ValidateDefaulted accepts an empty answer as "keep the value already in the
// field".
//
// It exists for accessible mode. huh's screen-reader path runs a field's
// validator on the raw line and only afterwards substitutes the field's default,
// and it never prints that default. A pre-filled field whose validator rejects
// "" therefore re-prompts on every bare Enter, so a screen-reader user cannot
// accept a value they cannot see.
func ValidateDefaulted(inner func(string) error) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		return inner(s)
	}
}

// ValidateDefaultedValue is ValidateDefaulted for a field whose pre-filled value
// may itself be empty, which is the normal case when the seed comes from a flag
// or the environment. Blank stays invalid when there is nothing to keep.
func ValidateDefaultedValue(prefilled string, inner func(string) error) func(string) error {
	if strings.TrimSpace(prefilled) == "" {
		return inner
	}
	return ValidateDefaulted(inner)
}

// runForm builds the form that describes one agent run.
//
// It is two groups rather than one because a multi-line task field plus four
// option fields do not fit a 24-row terminal, and huh clips whatever overflows.
// Groups are pages, which also gives the form a natural task-then-options shape.
func runForm(agents []string, defaultAgent, defaultModel string) (*huh.Form, *runAnswers) {
	answers := &runAnswers{
		Agent: defaultAgent,
		Model: defaultModel,
	}
	if answers.Agent == "" && len(agents) > 0 {
		answers.Agent = agents[0]
	}

	options := make([]huh.Option[string], 0, len(agents))
	for _, name := range agents {
		// Static options only: OptionsFunc yields an empty list in accessible mode.
		options = append(options, huh.NewOption(name, name))
	}

	taskField := huh.NewText().
		Title("Task").
		Description("what the agent should do; submitted as the value of -n").
		Placeholder("Summarize the repository layout and save it to /memory/layout.md").
		Lines(4).
		CharLimit(4000).
		Value(&answers.Task).
		// The task is never pre-filled, so a blank answer must still be refused.
		Validate(huh.ValidateNotEmpty())

	agentField := huh.NewSelect[string]().
		Title("Agent").
		Description("passed as -a/--agent").
		Options(options...).
		Value(&answers.Agent)

	modelField := huh.NewInput().
		Title("Model").
		Description("optional; passed as -M/--model. Blank uses the CLI's own default.").
		Placeholder("claude-sonnet-4-6").
		Value(&answers.Model).
		// Pre-filled from a flag or the environment, so a bare Enter must be able
		// to keep it under a screen reader; see ValidateDefaultedValue.
		Validate(ValidateDefaultedValue(defaultModel, validateModelName))

	approveField := huh.NewConfirm().
		Title("Approve tool calls automatically?").
		Description("passes -y. The agent then acts without asking, so leave it off unless you mean it.").
		Affirmative("Yes, auto-approve").
		Negative("No, ask me").
		Value(&answers.AutoApprove)

	shellField := huh.NewSelect[string]().
		Title("Shell commands").
		Description("passed as -S/--shell-allow-list. Disabled is the CLI's default.").
		Options(
			huh.NewOption("Disabled (CLI default)", ShellAllowDisabled),
			huh.NewOption("Recommended safe set", "recommended"),
			huh.NewOption("All commands", "all"),
		).
		Value(&answers.ShellAllow)

	form := huh.NewForm(
		huh.NewGroup(taskField).
			Title("Run an agent").
			Description("runs `deepagents -n <task> -q`, so the answer is the CLI's own"),
		huh.NewGroup(agentField, modelField, approveField, shellField).
			Title("Options"),
	)
	return themed(form), answers
}

// validateModelName rejects a model name that cannot be one.
//
// The value is passed to the CLI as a single argv element, so an embedded space
// is always a mistake -- a copy-pasted "claude sonnet" would reach the model
// resolver as an unknown model rather than as the typo it is. A blank value is
// valid and means "use the CLI's own default", which is why the empty string is
// accepted here and the whole field is not required.
func validateModelName(value string) error {
	if value == "" {
		return nil
	}
	if strings.ContainsAny(value, " \t") {
		return fmt.Errorf("a model name cannot contain spaces: %q", value)
	}
	return nil
}

// deleteForm builds the confirmation for a destructive thread deletion.
//
// The CLI deletes without asking, so the confirmation lives here rather than in
// argv. The thread's own identity is shown, because "delete thread ..." with no
// visible id is not something a user can verify.
func deleteForm(thread threads.ThreadInfo) (*huh.Form, *deleteAnswers) {
	answers := &deleteAnswers{}
	body := fmt.Sprintf("thread   %s\nagent    %s\nmessages %s\nupdated  %s\nprompt   %s",
		thread.ThreadID, thread.AgentLabel(), thread.Messages(), thread.Updated(),
		Truncate(thread.Prompt(), 160))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Delete this thread?").Description(body),
			huh.NewConfirm().
				Title("Run `deepagents threads delete` now?").
				Description("This removes the conversation and its checkpoints.").
				Affirmative("Delete it").
				Negative("Keep it").
				Value(&answers.Confirmed),
		).Title("Confirm deletion"),
	)
	return themed(form), answers
}
