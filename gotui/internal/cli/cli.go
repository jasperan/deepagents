// Package cli locates the real deepagents CLI and builds the argv/env for
// invoking it.
//
// Every subcommand and flag spelling here is taken from
// libs/cli/deepagents_cli/main.py and cross-checked against the JSON writers in
// libs/cli/deepagents_cli/output.py and sessions.py. Nothing is guessed: the
// `threads` subparsers and their short forms are declared there, and the
// envelopes this module's callers decode are produced by write_json, which is
// the CLI's documented machine-readable contract.
//
// No secret ever rides in argv. The Oracle password belongs in the child's
// environment (the DEEPAGENTS_*-prefixed variables OracleConfig reads), so an
// invocation stays safe to print in a pane, a log, or a process list.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// BinaryName is the shorter console script declared in libs/cli/pyproject.toml.
//
// `deepagents` and `deepagents-cli` are the same entry point
// (deepagents_cli:cli_main); the shorter one is the documented name.
const BinaryName = "deepagents"

// EntryPoint is the console-script target declared in pyproject.toml. It is what
// identifies the CLI's own package when this front-end has to find the checkout
// on its own, because the repository root has no pyproject.toml of its own.
const EntryPoint = "deepagents_cli:cli_main"

// EnvBinary overrides CLI discovery, for a checkout that is not on PATH.
const EnvBinary = "DEEPAGENTS_BIN"

// ErrBinaryNotFound is returned when the deepagents CLI cannot be located.
//
// This front-end cannot work without it, so the error names the fix rather than
// failing later with an opaque exec error.
var ErrBinaryNotFound = errors.New(
	"the deepagents CLI was not found; this front-end delegates to it rather than " +
		"reimplementing it, so install it (uv tool install deepagents-cli or " +
		"pip install deepagents-cli) or point " + EnvBinary + " at its entry point",
)

// Launcher is how the CLI is reached: a program, an optional argument prefix,
// and the working directory it should run in.
//
// Prefix exists because a source checkout is normally driven through uv
// (`uv run deepagents`), and uv resolves the project from its working
// directory. Bin+Prefix together are the command.
type Launcher struct {
	Bin    string
	Prefix []string
	// Dir is the child's working directory; empty inherits the parent's.
	Dir string
}

// Argv prefixes args with the program and its prefix.
func (l Launcher) Argv(args ...string) []string {
	out := make([]string, 0, 1+len(l.Prefix)+len(args))
	out = append(out, l.Bin)
	out = append(out, l.Prefix...)
	return append(out, args...)
}

// Describe renders a command for display.
func (l Launcher) Describe(args ...string) string { return strings.Join(l.Argv(args...), " ") }

// Resolve locates the deepagents CLI.
//
// The search order is: an explicit override, the DEEPAGENTS_BIN environment
// variable, `deepagents` on PATH, the checkout's own .venv entry point found
// above the working directory or above this executable, and finally
// `uv run deepagents` scoped to a directory whose pyproject.toml declares the
// entry point.
//
// The two starting points matter. This binary is run from wherever a shell
// happens to be, so looking only above the working directory would miss the
// checkout it was built from and fall through to a `uv run` that cannot resolve
// the project either.
func Resolve(override string) (Launcher, error) {
	if explicit := strings.TrimSpace(override); explicit != "" {
		return resolveExplicit(explicit)
	}
	if env := strings.TrimSpace(os.Getenv(EnvBinary)); env != "" {
		return resolveExplicit(env)
	}
	if path, err := exec.LookPath(BinaryName); err == nil {
		return Launcher{Bin: path}, nil
	}

	roots := searchRoots()
	for _, start := range roots {
		if venv := findVenvCLI(start); venv != "" {
			return Launcher{Bin: venv}, nil
		}
	}
	if uv, err := exec.LookPath("uv"); err == nil {
		for _, start := range roots {
			if project := findProjectDir(start); project != "" {
				// uv resolves the project from its working directory, so pin it;
				// otherwise `uv run` cannot find the entry point.
				return Launcher{Bin: uv, Prefix: []string{"run", BinaryName}, Dir: project}, nil
			}
		}
	}
	return Launcher{}, ErrBinaryNotFound
}

// resolveExplicit accepts either a name on PATH or a path to an entry point.
func resolveExplicit(value string) (Launcher, error) {
	if path, err := exec.LookPath(value); err == nil {
		return Launcher{Bin: path}, nil
	}
	if info, err := os.Stat(value); err == nil && !info.IsDir() {
		return Launcher{Bin: value}, nil
	}
	return Launcher{}, fmt.Errorf("%s %q is neither on PATH nor an existing file", EnvBinary, value)
}

// searchRoots is where discovery looks for the checkout: the working directory
// first, then the directory holding this executable, which during development
// is gotui/cmd/deepagents-tui inside the repository.
func searchRoots() []string {
	roots := make([]string, 0, 2)
	if cwd, err := os.Getwd(); err == nil {
		roots = append(roots, cwd)
	}
	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		roots = append(roots, filepath.Dir(exe))
	}
	return roots
}

// findVenvCLI walks up from start looking for a checkout's installed entry
// point. Both layouts are checked: a package-local .venv, and a repository-level
// one with the CLI under libs/cli.
func findVenvCLI(start string) string {
	dir := start
	for {
		candidates := []string{
			filepath.Join(dir, ".venv", "bin", BinaryName),
			filepath.Join(dir, "libs", "cli", ".venv", "bin", BinaryName),
		}
		for _, candidate := range candidates {
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// findProjectDir walks up from start looking for the directory that declares the
// CLI console script, which is what `uv run deepagents` needs as its working
// directory.
func findProjectDir(start string) string {
	dir := start
	for {
		candidate := filepath.Join(dir, "pyproject.toml")
		if data, err := os.ReadFile(candidate); err == nil &&
			bytes.Contains(data, []byte(EntryPoint)) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// EnvWith returns the child's environment.
//
// The child inherits the user's DEEPAGENTS_* variables (Oracle connection, model
// API keys) untouched, so it behaves exactly as the Python CLI would. Overrides
// replace rather than append: a duplicate entry would let a stale value win by
// position in some shells.
func EnvWith(overrides map[string]string) []string {
	env := os.Environ()
	if len(overrides) == 0 {
		return env
	}
	out := make([]string, 0, len(env)+len(overrides))
	for _, kv := range env {
		name, _, _ := strings.Cut(kv, "=")
		if _, replaced := overrides[name]; replaced {
			continue
		}
		out = append(out, kv)
	}
	for name, value := range overrides {
		out = append(out, name+"="+value)
	}
	return out
}

// Invocation is a ready-to-exec subprocess call.
type Invocation struct {
	Argv []string
	Env  []string
	Dir  string
}

// CommandLine renders argv for display. It cannot contain a secret, because no
// secret is ever placed in argv.
func (i Invocation) CommandLine() string { return strings.Join(i.Argv, " ") }

// Redact masks any DEEPAGENTS_* secret value that appears in a rendered string.
//
// No invocation in this front-end needs it, because secrets travel by
// environment, but a caller that prints the environment alongside a command must
// still not leak the password.
func Redact(s string, env []string) string {
	for _, kv := range env {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || len(value) < 4 || !strings.HasPrefix(name, "DEEPAGENTS_") {
			continue
		}
		if isSecretName(name) {
			s = strings.ReplaceAll(s, value, "********")
		}
	}
	return s
}

// isSecretName reports whether a DEEPAGENTS_* variable holds a credential.
func isSecretName(name string) bool {
	upper := strings.ToUpper(name)
	for _, marker := range []string{"PASSWORD", "TOKEN", "SECRET", "API_KEY"} {
		if strings.Contains(upper, marker) {
			return true
		}
	}
	return false
}

// ThreadsListOptions mirrors `deepagents threads list`.
type ThreadsListOptions struct {
	// Agent filters by agent name; empty shows every agent.
	Agent string
	// Limit caps the rows; zero omits the flag and lets the CLI apply its own
	// default (DA_CLI_RECENT_THREADS, else its built-in limit).
	Limit int
	// Sort is "created" or "updated"; empty omits the flag so the CLI reads its
	// own config file, which is the documented default.
	Sort string
	// Branch filters by git branch; empty shows every branch.
	Branch string
	// Verbose asks for the full column set.
	Verbose bool
}

// ThreadsList builds `deepagents threads list --json [flags]`.
//
// `--json` is always passed: this front-end is a machine consumer of the CLI's
// output and must never scrape the Rich table.
func ThreadsList(l Launcher, opts ThreadsListOptions) Invocation {
	args := []string{"threads", "list", "--json"}
	if a := strings.TrimSpace(opts.Agent); a != "" {
		args = append(args, "--agent", a)
	}
	if opts.Limit > 0 {
		args = append(args, "-n", strconv.Itoa(opts.Limit))
	}
	if opts.Sort != "" {
		args = append(args, "--sort", opts.Sort)
	}
	if b := strings.TrimSpace(opts.Branch); b != "" {
		args = append(args, "--branch", b)
	}
	if opts.Verbose {
		args = append(args, "-v")
	}
	return Invocation{Argv: l.Argv(args...), Env: EnvWith(nil), Dir: l.Dir}
}

// ThreadsDelete builds `deepagents threads delete <id> --json`.
func ThreadsDelete(l Launcher, threadID string) Invocation {
	return Invocation{
		Argv: l.Argv("threads", "delete", threadID, "--json"),
		Env:  EnvWith(nil),
		Dir:  l.Dir,
	}
}

// AgentsList builds `deepagents list --json`.
func AgentsList(l Launcher) Invocation {
	return Invocation{Argv: l.Argv("list", "--json"), Env: EnvWith(nil), Dir: l.Dir}
}

// RunOptions mirrors the run-relevant subset of the top-level CLI.
type RunOptions struct {
	// Task is the single non-interactive task. It is passed as the value of `-n`,
	// never as a bare positional, so a task beginning with "-" is not mistaken
	// for a flag.
	Task string
	// Agent selects the agent to run as (`-a/--agent`).
	Agent string
	// Model overrides the model (`-M/--model`); empty uses the CLI's default.
	Model string
	// Quiet sends the agent's response to stdout and everything else to stderr
	// (`-q/--quiet`). The CLI only permits it alongside -n or piped stdin.
	Quiet bool
	// NoStream passes `--no-stream`, buffering the answer instead of streaming
	// token by token.
	NoStream bool
	// AutoApprove passes `-y`, approving tool calls without prompting. It is
	// opt-in: it lets the agent act unattended.
	AutoApprove bool
	// ShellAllowList maps to `-S/--shell-allow-list`. Empty leaves shell
	// commands disabled, which is the CLI's default.
	ShellAllowList string
}

// Run builds `deepagents -n <task> [flags]`.
//
// The task is one argv element, so it needs no quoting here; the exec layer
// passes it through as a single argument.
func Run(l Launcher, opts RunOptions) (Invocation, error) {
	if strings.TrimSpace(opts.Task) == "" {
		return Invocation{}, errors.New("a task is required to run an agent")
	}
	args := []string{"-n", opts.Task}
	if opts.Quiet {
		args = append(args, "-q")
	}
	if a := strings.TrimSpace(opts.Agent); a != "" {
		args = append(args, "-a", a)
	}
	if m := strings.TrimSpace(opts.Model); m != "" {
		args = append(args, "-M", m)
	}
	if opts.NoStream {
		args = append(args, "--no-stream")
	}
	if opts.AutoApprove {
		args = append(args, "-y")
	}
	if s := strings.TrimSpace(opts.ShellAllowList); s != "" {
		args = append(args, "-S", s)
	}
	return Invocation{Argv: l.Argv(args...), Env: EnvWith(nil), Dir: l.Dir}, nil
}
