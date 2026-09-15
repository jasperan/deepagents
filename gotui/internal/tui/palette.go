package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// The approved design tokens, shared with the huh forms through huhstyle.
//
// These 14 values are the only colours any source file in this TUI may contain.
// Keeping them in one place is what makes that rule checkable rather than
// aspirational, and it keeps the chrome (headers, status lines, error banners)
// in the same palette as the themed forms.
//
// huh's built-in ThemeCatppuccin() is NOT used: it paints with auxiliary
// Catppuccin shades that are not approved tokens.
const (
	hexBG       = "#1e1e2e" // bg
	hexSurface  = "#181825" // surface
	hexElevated = "#313244" // elevated
	hexText     = "#cdd6f4" // text
	hexSubtext  = "#a6adc8" // subtext
	hexMuted    = "#6c7086" // muted
	hexDim      = "#585b70" // dim
	hexPrimary  = "#89b4fa" // primary
	hexInfo     = "#89dceb" // info
	hexSuccess  = "#a6e3a1" // success
	hexWarning  = "#f9e2af" // warning
	hexError    = "#f38ba8" // error

	// Declared for completeness of the token set; the browser has no use for
	// either yet, and referencing them keeps the palette reviewable as a whole.
	hexHighest   = "#45475a" // highest
	hexSecondary = "#cba6f7" // secondary
)

var (
	_ = hexHighest
	_ = hexSecondary
)

// MinWidth is the narrowest width this UI renders at.
//
// A terminal can be resized to a handful of columns, and every width-derived
// computation below subtracts a border or a padding from the terminal width. At
// one or two columns those subtractions go negative and a slice operation into
// the result panics, so all width maths is funnelled through ClampWidth first.
const MinWidth = 6

// ClampWidth floors a terminal width at MinWidth.
func ClampWidth(width int) int {
	if width < MinWidth {
		return MinWidth
	}
	return width
}

// Chrome styles. Titles are bold primary; metadata is subtext; errors are the
// only place the error token appears, because it is semantic.
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(hexPrimary))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(hexMuted))
	subtextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(hexSubtext))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(hexSuccess))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(hexWarning))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(hexError))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(hexInfo))
	panelStyle   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(hexDim)).
			Padding(0, 1)
)

// errorLine prefixes a message with the error token so a failure is visually
// unmistakable, without leaking a colour literal into the call sites.
func errorLine(msg string) string { return errorStyle.Render("error: ") + msg }

// invocationLine shows the exact command that ran. No secret is ever part of an
// argv in this front-end, so this is safe to display verbatim.
func invocationLine(invocation string) string { return mutedStyle.Render("$ " + invocation) }

// Pane draws a titled panel, clamped so it always has interior room.
func Pane(title, body string, width int) string {
	width = ClampWidth(width)
	if width < 12 {
		// Too narrow for a border plus padding to leave room for content, so drop
		// the chrome and render the text plainly rather than a 0-column box.
		if title == "" {
			return body
		}
		return titleStyle.Render(title) + "\n" + body
	}
	heading := ""
	if title != "" {
		heading = titleStyle.Render(title) + "\n"
	}
	// Width sets the whole block (border and padding included), so passing the
	// terminal width here keeps the panel inside it.
	return panelStyle.Width(width).Render(heading + body)
}

// Header is the persistent identity bar: name on the left, context on the right.
func Header(left, right string, width int) string {
	width = ClampWidth(width)
	title := titleStyle.Render(left)
	context := subtextStyle.Render(right)
	gap := width - lipgloss.Width(title) - lipgloss.Width(context) - 2
	if gap < 1 {
		// No room to align: stack the context, or drop it when even that is too
		// narrow to be readable.
		if width < 20 || context == "" {
			return title
		}
		return title + "\n" + context
	}
	return title + strings.Repeat(" ", gap) + context
}

// Footer is the key-hint bar.
func Footer(hint string, width int) string {
	width = ClampWidth(width)
	return subtextStyle.Render(Truncate(hint, width-1))
}

// Status is a one-line status message rendered in a semantic colour.
func Status(kind, message string, width int) string {
	width = ClampWidth(width)
	style := subtextStyle
	switch kind {
	case "error":
		style = errorStyle
	case "success":
		style = successStyle
	case "warning":
		style = warningStyle
	case "busy":
		style = infoStyle
	}
	return style.Render(Truncate(message, width-1))
}

// Pad right-pads s to n columns.
func Pad(s string, n int) string {
	if n <= 0 {
		return s
	}
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

// Truncate shortens s to n runes, adding a tilde when it cuts.
//
// n <= 0 returns the empty string rather than panicking on a negative slice
// bound, which is the shape the narrow-terminal bug took.
func Truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n == 1 {
		return "~"
	}
	return string(runes[:n-1]) + "~"
}

// Wrap soft-wraps text to width columns, preserving existing line breaks.
func Wrap(s string, width int) string {
	width = ClampWidth(width)
	var out []string
	for _, line := range strings.Split(s, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			out = append(out, "")
			continue
		}
		current := ""
		for _, word := range fields {
			switch {
			case current == "":
				current = word
			case len(current)+1+len(word) <= width:
				current += " " + word
			default:
				out = append(out, current)
				current = word
			}
		}
		out = append(out, current)
	}
	return strings.Join(out, "\n")
}
