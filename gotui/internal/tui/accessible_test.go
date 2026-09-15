package tui

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

// This file covers the accessible (screen-reader) path.
//
// Why it exists: huh's accessible mode is implemented only in Form.Run -- the
// flag is never read by an embedded form -- and its prompt path runs a field's
// validator on the raw line and only afterwards substitutes the field's default,
// without printing it. A pre-filled field whose validator rejects "" therefore
// re-prompts on every bare Enter, so a screen-reader user cannot accept a value
// they cannot see.
//
// The unit tests below pin the validation contract. The end-to-end test drives
// the real form through huh's accessible path with scripted input, which is the
// only way a form can be driven in a test.

var errBlank = errors.New("blank answer rejected")

func TestValidateDefaultedAcceptsBlank(t *testing.T) {
	inner := func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errBlank
		}
		if s == "bad" {
			return errors.New("not usable")
		}
		return nil
	}
	wrapped := ValidateDefaulted(inner)

	for _, in := range []string{"", "   ", "\t"} {
		if err := wrapped(in); err != nil {
			t.Errorf("ValidateDefaulted(inner)(%q) = %v, want nil", in, err)
		}
	}
	if err := wrapped("bad"); err == nil {
		t.Error("ValidateDefaulted(inner)(\"bad\") = nil, want the inner validator's error")
	}
	if err := wrapped("claude-sonnet-4-6"); err != nil {
		t.Errorf("ValidateDefaulted(inner)(valid) = %v, want nil", err)
	}
}

// The guard: a field is usually seeded from a flag or the environment, and
// either may be absent. Blank stays invalid when there is nothing to keep.
func TestValidateDefaultedValueOnlyRelaxesWhenSomethingIsKept(t *testing.T) {
	inner := func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errBlank
		}
		return nil
	}

	if err := ValidateDefaultedValue("", inner)(""); err == nil {
		t.Error("empty seed accepted a blank answer; the required field was weakened")
	}
	if err := ValidateDefaultedValue("   ", inner)(""); err == nil {
		t.Error("whitespace-only seed accepted a blank answer")
	}
	if err := ValidateDefaultedValue("claude-sonnet-4-6", inner)(""); err != nil {
		t.Errorf("seeded field rejected blank: %v", err)
	}
}

func TestValidateModelNameRejectsSpacesButAllowsBlank(t *testing.T) {
	// Blank means "use the CLI's own default", which is a real choice.
	if err := validateModelName(""); err != nil {
		t.Errorf("validateModelName(\"\") = %v, want nil", err)
	}
	if err := validateModelName("claude-sonnet-4-6"); err != nil {
		t.Errorf("validateModelName(valid) = %v, want nil", err)
	}
	// A copy-pasted "claude sonnet" would otherwise reach the model resolver.
	for _, bad := range []string{"claude sonnet", "gpt 5.2", "model\tname"} {
		if err := validateModelName(bad); err == nil {
			t.Errorf("validateModelName(%q) = nil, want a rejection", bad)
		}
	}
}

// runFormAccessible drives the real run form through huh's accessible path with
// scripted input, returning everything it wrote.
func runFormAccessible(t *testing.T, agents []string, seededModel, input string) string {
	t.Helper()
	var out bytes.Buffer
	form, _ := runForm(agents, "coder", seededModel)
	form = form.
		WithAccessible(true).
		WithInput(strings.NewReader(input)).
		WithOutput(&out)

	done := make(chan error, 1)
	go func() { done <- form.Run() }()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Run returned %v (a canceled form is not a failure here)", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("the accessible form did not finish within 10s; output so far:\n%s", out.String())
	}
	return out.String()
}

// The accessible path must render as plain prompts and must not reject a blank
// answer on a field that has a seeded value to keep: a screen-reader user is
// never shown that value, so they cannot retype it.
func TestAccessibleFormDoesNotRejectABlankSeededModel(t *testing.T) {
	// One line for the task, then blank answers for the remaining fields.
	out := runFormAccessible(t, []string{"coder", "researcher"}, "claude-sonnet-4-6",
		"summarize the repo\n\n\n\n\n")

	if strings.Contains(out, "cannot contain spaces") {
		t.Errorf("a blank answer on the seeded model field was rejected, so a screen-reader "+
			"user cannot keep it.\noutput:\n%s", out)
	}
	if !strings.Contains(out, "Task") {
		t.Errorf("the accessible form never asked for the task.\noutput:\n%s", out)
	}
}

// The negative half, on the real form: the required task has no default to keep,
// so a blank answer must be refused rather than silently accepted.
//
// huh's accessible path substitutes a field's default only after validation, and
// this field has none, so the inner validator is the thing standing between a
// blank task and a confusing "task is required" failure from the CLI.
func TestAccessibleFormKeepsRequiredFieldsRequired(t *testing.T) {
	out := runFormAccessible(t, []string{"coder"}, "", "\n\n\n")

	if !strings.Contains(out, "Task") {
		t.Fatalf("the form never reached the task field.\noutput:\n%s", out)
	}
	if !strings.Contains(out, "input cannot be empty") {
		t.Errorf("a blank required task was accepted; the validator never ran.\noutput:\n%s", out)
	}
}

// A checkout with no agents stored, or an agent query that failed, must still
// produce a usable form: omitting -a is exactly what makes the CLI use its own
// default agent.
func TestRunFormOpensWithNoAgentList(t *testing.T) {
	form, answers := runForm(nil, "", "")
	if form == nil || answers == nil {
		t.Fatal("runForm returned nil for an empty agent list")
	}
	if answers.Agent != "" {
		t.Errorf("Agent = %q, want empty so -a is omitted and the CLI uses its default", answers.Agent)
	}
}
