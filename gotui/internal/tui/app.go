// Package tui is the Bubble Tea v2 front-end for deepagents.
//
// It is an ADDITIONAL way to run the project, not a reimplementation of it. The
// thread browser reads the CLI's own `threads list --json` envelope, deletion
// shells out to `deepagents threads delete`, and a run shells out to
// `deepagents -n <task> -q`. Nothing about agent execution, thread storage or
// tool policy is reproduced here, so a Go user and a Python user get identical
// results from identical inputs.
//
// The model is one message pump. While a huh form is open the form sees every
// message, so its own key handling stays authoritative; outside a form, the
// active screen owns the keys. Each key is therefore interpreted in exactly one
// place, which is what keeps a single keystroke from being acted on twice.
package tui

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"

	"github.com/jasperan/deepagents/gotui/internal/cli"
	"github.com/jasperan/deepagents/gotui/internal/threads"
)

// screen is which pane the TUI is showing.
type screen int

const (
	screenThreads screen = iota
	screenDetail
	screenRunForm
	screenRunning
	screenOutput
	screenDelete
)

// Thread table columns whose width does not depend on the terminal.
const (
	colAgentWidth   = 12
	colMsgsWidth    = 5
	colUpdatedWidth = 16
	colIDWidth      = 22
)

// Messages.
type (
	threadsMsg struct {
		list []threads.ThreadInfo
		err  error
	}
	agentsMsg struct {
		list []threads.AgentInfo
		err  error
	}
	deleteMsg struct {
		result threads.DeleteResult
		err    error
	}
	streamLineMsg Line
	streamDoneMsg struct{ err error }
)

// Options configures a Model.
type Options struct {
	// Launcher is the resolved deepagents CLI.
	Launcher cli.Launcher
	// Runner executes it; the zero value falls back to ShellRunner.
	Runner Runner
	// Limit caps the thread listing; zero lets the CLI apply its own default.
	Limit int
	// Agent filters the listing to one agent; empty shows all.
	Agent string
	// Sort is "created" or "updated"; empty defers to the CLI's config.
	Sort string
	// Branch filters the listing by git branch; empty shows all.
	Branch string
	// Model seeds the run form's model field (flag or environment).
	Model string
	// Width and Height are the initial dimensions, before the first
	// WindowSizeMsg arrives.
	Width, Height int
}

// Model is the root Bubble Tea model.
type Model struct {
	launcher cli.Launcher
	runner   Runner

	// Listing parameters, applied to every fetch.
	filterAgent string
	limit       int
	sortBy      string
	branch      string

	screen        screen
	width, height int

	// Thread browser.
	all    []threads.ThreadInfo
	shown  []threads.ThreadInfo
	table  table.Model
	filter textinput.Model
	// filtering is true while the filter field has focus.
	filtering bool

	// Detail pane.
	detail viewport.Model
	target threads.ThreadInfo

	// Forms. Exactly one is non-nil at a time.
	runForm       *huh.Form
	runAnswers    *runAnswers
	deleteForm    *huh.Form
	deleteAnswers *deleteAnswers

	// Run pane.
	spinner    spinner.Model
	output     viewport.Model
	lines      []string
	invocation string
	runErr     error
	stream     chan Line
	done       chan error

	// Status line.
	status     string
	statusKind string
	loading    bool

	quitting bool
}

// New builds the model and immediately requests the thread listing.
func New(opts Options) Model {
	runner := opts.Runner
	if runner == nil {
		runner = NewShellRunner()
	}
	width, height := opts.Width, opts.Height
	if width <= 0 {
		width = 100
	}
	if height <= 0 {
		height = 30
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = infoStyle

	filter := textinput.New()
	filter.Placeholder = "filter threads"
	filter.Prompt = "/ "

	m := Model{
		launcher:    opts.Launcher,
		runner:      runner,
		filterAgent: strings.TrimSpace(opts.Agent),
		limit:       opts.Limit,
		sortBy:      opts.Sort,
		branch:      opts.Branch,
		screen:      screenThreads,
		width:       width,
		height:      height,
		filter:      filter,
		spinner:     sp,
		output:      viewport.New(),
		detail:      viewport.New(),
		loading:     true,
		status:      "loading threads…",
		statusKind:  "busy",
		runAnswers:  &runAnswers{Model: opts.Model},
	}
	m.table = table.New(
		table.WithColumns(m.columns()),
		table.WithHeight(m.tableHeight()),
		table.WithFocused(true),
	)
	m.applySize()
	return m
}

// Init starts the thread listing and the spinner.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadThreads(), m.spinner.Tick)
}

// loadThreads asks the CLI for the thread listing.
func (m *Model) loadThreads() tea.Cmd {
	runner := m.runner
	inv := cli.ThreadsList(m.launcher, cli.ThreadsListOptions{
		Agent:  m.filterAgent,
		Limit:  m.limit,
		Sort:   m.sortBy,
		Branch: m.branch,
	})
	return func() tea.Msg {
		out, err := runner.Capture(context.Background(), inv)
		if err != nil {
			return threadsMsg{err: err}
		}
		list, err := threads.ParseThreadsList([]byte(out))
		if err != nil {
			return threadsMsg{err: err}
		}
		return threadsMsg{list: list}
	}
}

// loadAgents asks the CLI which agents exist, so the run form offers real names.
//
// A failure is not fatal: the CLI's own default agent still works, so the form
// opens with free text rather than blocking the user.
func (m *Model) loadAgents() tea.Cmd {
	runner := m.runner
	inv := cli.AgentsList(m.launcher)
	return func() tea.Msg {
		out, err := runner.Capture(context.Background(), inv)
		if err != nil {
			return agentsMsg{err: err}
		}
		agents, err := threads.ParseAgents([]byte(out))
		if err != nil {
			return agentsMsg{err: err}
		}
		return agentsMsg{list: agents}
	}
}

// deleteThread asks the CLI to delete one thread.
func (m *Model) deleteThread(id string) tea.Cmd {
	runner := m.runner
	inv := cli.ThreadsDelete(m.launcher, id)
	return func() tea.Msg {
		out, err := runner.Capture(context.Background(), inv)
		if err != nil {
			return deleteMsg{err: err}
		}
		result, err := threads.ParseThreadsDelete([]byte(out))
		if err != nil {
			return deleteMsg{err: err}
		}
		return deleteMsg{result: result}
	}
}

// startRun streams an agent run into the output pane.
func (m *Model) startRun() tea.Cmd {
	inv, err := cli.Run(m.launcher, cli.RunOptions{
		Task:           m.runAnswers.Task,
		Agent:          m.runAnswers.Agent,
		Model:          m.runAnswers.Model,
		Quiet:          true,
		AutoApprove:    m.runAnswers.AutoApprove,
		ShellAllowList: m.runAnswers.ShellAllow,
	})
	if err != nil {
		m.runErr = err
		m.lines = []string{errorLine(err.Error())}
		m.output.SetContentLines(m.lines)
		m.screen = screenOutput
		return nil
	}

	m.invocation = inv.CommandLine()
	m.lines = []string{invocationLine(m.invocation)}
	m.output.SetContentLines(m.lines)
	m.output.GotoBottom()
	m.runErr = nil
	m.screen = screenRunning
	m.status = "running…"
	m.statusKind = "busy"

	runner := m.runner
	lines := make(chan Line, 128)
	done := make(chan error, 1)
	m.stream, m.done = lines, done

	go func() {
		err := runner.Stream(context.Background(), inv, func(l Line) { lines <- l })
		close(lines)
		done <- err
	}()

	return tea.Batch(waitForStream(lines), waitForDone(done))
}

func waitForStream(ch <-chan Line) tea.Cmd {
	return func() tea.Msg {
		l, ok := <-ch
		if !ok {
			return nil
		}
		return streamLineMsg(l)
	}
}

func waitForDone(ch <-chan error) tea.Cmd {
	return func() tea.Msg { return streamDoneMsg{err: <-ch} }
}

// Update is the single message pump.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.applySize()
		return m, nil

	case tea.KeyPressMsg:
		// Ctrl+C is the one binding that works from every screen, including
		// inside a form, so an accidental modal can never trap the user.
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case threadsMsg:
		m.loading = false
		if msg.err != nil {
			m.setStatus("error", msg.err.Error())
			// Keep whatever was already listed: a failed refresh should not wipe
			// the rows the user was looking at.
			return m, nil
		}
		m.all = msg.list
		m.applyFilter()
		if len(m.all) == 0 {
			m.setStatus("", "no threads recorded yet — press n to run an agent")
		} else {
			m.setStatus("", fmt.Sprintf("%d thread(s)", len(m.all)))
		}
		return m, nil

	case agentsMsg:
		if msg.err != nil {
			m.setStatus("warning", "could not list agents; using the CLI default: "+msg.err.Error())
		}
		return m, m.openRunForm(msg.list)

	case deleteMsg:
		if msg.err != nil {
			m.setStatus("error", "delete failed: "+msg.err.Error())
			return m, nil
		}
		if !msg.result.Deleted {
			m.setStatus("warning", "thread "+msg.result.ThreadID+" was not found")
			return m, nil
		}
		m.setStatus("success", "deleted thread "+msg.result.ThreadID)
		m.loading = true
		return m, m.loadThreads()

	case streamLineMsg:
		m.appendLine(Line(msg))
		return m, waitForStream(m.stream)

	case streamDoneMsg:
		m.runErr = msg.err
		m.screen = screenOutput
		if msg.err != nil {
			m.setStatus("error", msg.err.Error())
		} else {
			m.setStatus("success", "run finished")
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// A form is authoritative while it is open.
	if m.runForm != nil {
		return m.updateRunForm(msg)
	}
	if m.deleteForm != nil {
		return m.updateDeleteForm(msg)
	}

	switch m.screen {
	case screenRunning:
		return m.updateRunning(msg)
	case screenThreads:
		return m.updateThreads(msg)
	case screenDetail:
		return m.updateDetail(msg)
	case screenOutput:
		return m.updateOutput(msg)
	}
	return m, nil
}

// updateThreads handles the browser's keys.
func (m Model) updateThreads(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.filtering {
		return m.updateFilter(msg)
	}

	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "r":
		m.loading = true
		m.setStatus("busy", "refreshing…")
		return m, m.loadThreads()
	case "n":
		m.setStatus("busy", "loading agents…")
		return m, m.loadAgents()
	case "/":
		m.filtering = true
		m.filter.Focus()
		return m, nil
	case "d":
		thread, ok := m.selected()
		if !ok {
			return m, nil
		}
		// Deletion is destructive and the CLI does not ask, so it always goes
		// through a confirmation form rather than a bare keystroke.
		form, answers := deleteForm(thread)
		m.deleteForm, m.deleteAnswers, m.target = form, answers, thread
		m.screen = screenDelete
		return m, tea.Batch(form.Init(), m.spinner.Tick)
	case "enter":
		thread, ok := m.selected()
		if !ok {
			return m, nil
		}
		m.target = thread
		m.detail.SetContent(m.detailText(thread))
		m.detail.GotoTop()
		m.screen = screenDetail
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// updateFilter drives the filter field. Enter accepts the filter, esc clears it.
func (m Model) updateFilter(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "enter":
			m.filtering = false
			m.filter.Blur()
			m.applyFilter()
			return m, nil
		case "esc":
			m.filtering = false
			m.filter.Blur()
			m.filter.SetValue("")
			m.applyFilter()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	m.applyFilter()
	return m, cmd
}

// updateDetail handles the detail pane.
func (m Model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "q", "esc", "enter":
			m.screen = screenThreads
			return m, nil
		case "d":
			form, answers := deleteForm(m.target)
			m.deleteForm, m.deleteAnswers = form, answers
			m.screen = screenDelete
			return m, tea.Batch(form.Init(), m.spinner.Tick)
		}
	}
	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}

// updateRunning collects streamed output and ends the run.
func (m Model) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case streamLineMsg:
		m.appendLine(Line(msg))
		return m, waitForStream(m.stream)
	case streamDoneMsg:
		m.runErr = msg.err
		m.screen = screenOutput
		if msg.err != nil {
			m.setStatus("error", msg.err.Error())
		} else {
			m.setStatus("success", "run finished")
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	var cmd tea.Cmd
	m.output, cmd = m.output.Update(msg)
	return m, cmd
}

// updateOutput handles the finished-run pane.
func (m Model) updateOutput(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "q", "esc", "enter":
			m.screen = screenThreads
			m.loading = true
			return m, m.loadThreads()
		case "r":
			if strings.TrimSpace(m.runAnswers.Task) != "" {
				return m, m.startRun()
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.output, cmd = m.output.Update(msg)
	return m, cmd
}

// updateRunForm drives the run form and reacts to its terminal states.
func (m Model) updateRunForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
		m.runForm = nil
		m.screen = screenThreads
		return m, nil
	}
	updated, cmd := m.runForm.Update(msg)
	if f, ok := updated.(*huh.Form); ok {
		m.runForm = f
	}
	switch m.runForm.State {
	case huh.StateCompleted:
		m.runForm = nil
		return m, m.startRun()
	case huh.StateAborted:
		m.runForm = nil
		m.screen = screenThreads
		return m, nil
	}
	return m, cmd
}

// updateDeleteForm drives the confirmation and reacts to its terminal states.
func (m Model) updateDeleteForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
		m.deleteForm = nil
		m.screen = screenThreads
		return m, nil
	}
	updated, cmd := m.deleteForm.Update(msg)
	if f, ok := updated.(*huh.Form); ok {
		m.deleteForm = f
	}
	switch m.deleteForm.State {
	case huh.StateCompleted:
		confirmed := m.deleteAnswers != nil && m.deleteAnswers.Confirmed
		id := m.target.ThreadID
		m.deleteForm = nil
		m.screen = screenThreads
		if !confirmed || id == "" {
			m.setStatus("", "deletion cancelled; nothing was removed")
			return m, nil
		}
		m.setStatus("busy", "deleting "+id+"…")
		return m, m.deleteThread(id)
	case huh.StateAborted:
		m.deleteForm = nil
		m.screen = screenThreads
		return m, nil
	}
	return m, cmd
}

// openRunForm builds the run form once the agent list is known.
func (m *Model) openRunForm(agents []threads.AgentInfo) tea.Cmd {
	names := threads.AgentNames(agents)
	defaultAgent := threads.DefaultAgentName(agents)
	if len(names) == 0 {
		// No agent list (a fresh install, or a failed query) is not a dead end:
		// the CLI uses its own default agent when -a is omitted, so the form still
		// works with a single blank selector entry, which omits the flag.
		names = []string{""}
	}
	if defaultAgent == "" {
		defaultAgent = names[0]
	}
	form, answers := runForm(names, defaultAgent, m.runAnswers.Model)
	m.runForm, m.runAnswers = form, answers
	m.screen = screenRunForm
	return form.Init()
}

// applyFilter recomputes the rows shown for the current filter text.
func (m *Model) applyFilter() {
	m.shown = threads.FilterThreads(m.all, m.filter.Value())
	rows := make([]table.Row, 0, len(m.shown))
	for _, t := range m.shown {
		rows = append(rows, table.Row{
			Truncate(t.ThreadID, colIDWidth),
			Truncate(t.AgentLabel(), colAgentWidth),
			t.Messages(),
			t.Updated(),
			Truncate(t.Prompt(), m.promptColumnWidth()),
		})
	}
	m.table.SetRows(rows)
	if cursor := m.table.Cursor(); cursor >= len(rows) && len(rows) > 0 {
		m.table.SetCursor(len(rows) - 1)
	}
	if len(rows) == 0 {
		m.table.SetCursor(0)
	}
}

// selected returns the thread under the table cursor.
func (m Model) selected() (threads.ThreadInfo, bool) {
	index := m.table.Cursor()
	if index < 0 || index >= len(m.shown) {
		return threads.ThreadInfo{}, false
	}
	return m.shown[index], true
}

// columns lays out the table for the current width.
//
// Every width is clamped positive: a table given a negative column width panics,
// and a terminal can be resized to a few columns.
func (m Model) columns() []table.Column {
	width := ClampWidth(m.width)
	idWidth := minInt(colIDWidth, maxInt(6, width/5))
	agentWidth := minInt(colAgentWidth, maxInt(5, width/7))
	updatedWidth := minInt(colUpdatedWidth, maxInt(6, width/6))
	fixed := idWidth + agentWidth + colMsgsWidth + updatedWidth + 4
	promptWidth := maxInt(6, width-fixed)
	return []table.Column{
		{Title: "Thread", Width: idWidth},
		{Title: "Agent", Width: agentWidth},
		{Title: "Msgs", Width: colMsgsWidth},
		{Title: "Updated", Width: updatedWidth},
		{Title: "Prompt", Width: promptWidth},
	}
}

// promptColumnWidth is the width budget of the prompt column.
func (m Model) promptColumnWidth() int {
	columns := m.columns()
	return columns[len(columns)-1].Width
}

// tableHeight sizes the table to the rows available under the header and hints.
func (m Model) tableHeight() int {
	height := maxInt(3, m.height-8)
	// Do not render more blank rows than there are threads.
	if len(m.shown) > 0 && height > len(m.shown)+1 {
		height = len(m.shown) + 1
	}
	return height
}

// applySize re-applies the terminal dimensions to every pane.
func (m *Model) applySize() {
	width := ClampWidth(m.width)
	m.table.SetColumns(m.columns())
	m.table.SetWidth(width)
	m.table.SetHeight(m.tableHeight())

	m.filter.SetWidth(maxInt(4, width-4))

	m.output.SetWidth(maxInt(4, width-2))
	m.output.SetHeight(maxInt(2, m.height-6))

	m.detail.SetWidth(maxInt(4, width-2))
	m.detail.SetHeight(maxInt(2, m.height-4))

	// Forms must be resized too, or a narrow terminal wraps their borders and
	// huh clips the fields that overflow.
	m.resizeForm(m.runForm)
	m.resizeForm(m.deleteForm)
}

// resizeForm pins a form to the model's dimensions.
func (m *Model) resizeForm(f *huh.Form) {
	if f == nil {
		return
	}
	f.WithWidth(maxInt(20, m.width-4))
	f.WithHeight(maxInt(6, m.height-6))
}

// setStatus records a status line message.
func (m *Model) setStatus(kind, message string) {
	m.status = message
	m.statusKind = kind
}

// appendLine adds one line of subprocess output, keeping the view at the bottom
// where progress appears.
func (m *Model) appendLine(l Line) {
	text := l.Text
	if l.IsStderr() {
		text = warningStyle.Render(text)
	}
	m.lines = append(m.lines, text)
	m.output.SetContentLines(m.lines)
	m.output.GotoBottom()
}

// detailText renders one thread's metadata for the detail pane.
func (m Model) detailText(t threads.ThreadInfo) string {
	lines := []string{
		titleStyle.Render("Thread " + t.ThreadID),
		"",
		subtextStyle.Render("agent     ") + t.AgentLabel(),
		subtextStyle.Render("messages  ") + t.Messages(),
		subtextStyle.Render("branch    ") + t.Branch(),
		subtextStyle.Render("cwd       ") + t.WorkingDir(),
		subtextStyle.Render("created   ") + t.Created(),
		subtextStyle.Render("updated   ") + t.Updated(),
	}
	if t.LatestCheckpointID != nil && strings.TrimSpace(*t.LatestCheckpointID) != "" {
		lines = append(lines, subtextStyle.Render("checkpoint ")+Truncate(*t.LatestCheckpointID, 40))
	}
	lines = append(lines, "", titleStyle.Render("Opening prompt"), Wrap(t.Prompt(), maxInt(20, m.detail.Width()-2)))
	return strings.Join(lines, "\n")
}

// View renders the active screen.
func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

// render builds the active screen's text.
//
// Every width used here is clamped, so rendering at any terminal size — including
// the handful of columns a user can resize down to — is safe.
func (m Model) render() string {
	width := ClampWidth(m.width)
	var b strings.Builder

	b.WriteString(Header("deepagents", m.headerContext(), width))
	b.WriteString("\n")

	switch {
	case m.runForm != nil:
		b.WriteString(m.formPane("Run an agent", m.runForm, width))
	case m.deleteForm != nil:
		b.WriteString(m.formPane("Delete thread", m.deleteForm, width))
	case m.screen == screenRunning:
		b.WriteString(m.runningPane(width))
	case m.screen == screenOutput:
		b.WriteString(m.outputPane(width))
	case m.screen == screenDetail:
		b.WriteString(Pane("Thread detail", m.detail.View(), width))
	default:
		b.WriteString(m.threadsPane(width))
	}

	b.WriteString("\n")
	b.WriteString(Status(m.statusKind, m.status, width))
	b.WriteString("\n")
	b.WriteString(Footer(m.hint(), width))
	return b.String()
}

// headerContext is the right-hand side of the header bar.
func (m Model) headerContext() string {
	switch m.screen {
	case screenRunning:
		return "running"
	case screenDetail:
		return "thread detail"
	case screenOutput:
		return "last run"
	}
	if m.loading {
		return "loading"
	}
	return fmt.Sprintf("%d thread(s)", len(m.shown))
}

// hint is the key legend for the active screen.
func (m Model) hint() string {
	switch {
	case m.runForm != nil:
		return "tab next  •  shift+tab back  •  enter submit  •  esc cancel"
	case m.deleteForm != nil:
		return "enter confirm  •  esc cancel"
	}
	switch m.screen {
	case screenRunning:
		return m.spinner.View() + " streaming output  •  ctrl+c quit"
	case screenOutput:
		return "q back  •  r rerun  •  ctrl+c quit"
	case screenDetail:
		return "d delete  •  q back  •  ctrl+c quit"
	}
	if m.filtering {
		return "type to filter  •  enter accept  •  esc clear"
	}
	return "enter detail  •  d delete  •  n new run  •  / filter  •  r refresh  •  q quit"
}

// formPane renders a form with a title, indented to sit inside the header.
func (m Model) formPane(title string, form *huh.Form, width int) string {
	body := strings.ReplaceAll(form.View(), "\n", "\n  ")
	return Pane(title, "  "+body, width)
}

// threadsPane renders the browser: an optional filter field, the table, and an
// empty state that says what to do next.
func (m Model) threadsPane(width int) string {
	var b strings.Builder
	if m.filtering || m.filter.Value() != "" {
		b.WriteString(m.filter.View())
		b.WriteString("\n")
	}
	if len(m.shown) == 0 {
		if len(m.all) == 0 {
			b.WriteString(subtextStyle.Render("No threads are stored yet."))
			b.WriteString("\n\n")
			b.WriteString(subtextStyle.Render("Press n to run an agent; the conversation is then listed here."))
		} else {
			b.WriteString(subtextStyle.Render("No thread matches the filter."))
			b.WriteString("\n\n")
			b.WriteString(subtextStyle.Render("Press esc to clear it."))
		}
		return Pane("Threads", b.String(), width)
	}
	b.WriteString(m.table.View())
	return Pane("Threads", b.String(), width)
}

// runningPane renders the spinner above the live output pane.
func (m Model) runningPane(width int) string {
	header := m.spinner.View() + " " + subtextStyle.Render(Truncate(m.invocation, maxInt(10, width-6)))
	return Pane("Running", header+"\n\n"+m.output.View(), width)
}

// outputPane renders the finished (or failed) run.
func (m Model) outputPane(width int) string {
	title := "Run finished"
	if m.runErr != nil {
		title = "Run failed"
	}
	return Pane(title, m.output.View(), width)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Drain runs a command tree to completion, discarding timer ticks.
//
// huh re-arms the text input's blink on every update, and the spinner re-arms
// its own tick, so feeding either message back into Update re-arms it forever
// and an undiscriminating drain never terminates: it would walk `depth` levels
// of a 100ms timer and take seconds per call. Both are dropped here. Exported
// because the tests use it.
func Drain(m tea.Model, cmd tea.Cmd, depth int) tea.Model {
	if cmd == nil || depth > 64 {
		return m
	}
	msg := cmd()
	if msg == nil {
		return m
	}
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			m = Drain(m, c, depth+1)
		}
		return m
	}
	switch msg.(type) {
	case cursor.BlinkMsg, spinner.TickMsg:
		// A timer tick only asks for another frame; dropping it is what makes
		// this terminate.
		return m
	}
	next, nextCmd := m.Update(msg)
	return Drain(next, nextCmd, depth+1)
}
