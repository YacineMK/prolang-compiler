package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Styles for different output types
var (
	// Colors
	colorPrimary = lipgloss.Color("33")
	colorSuccess = lipgloss.Color("34")
	colorWarning = lipgloss.Color("208")
	colorError   = lipgloss.Color("196")
	colorInfo    = lipgloss.Color("39")
	colorMuted   = lipgloss.Color("243")

	// Styles
	styleTitle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Margin(0, 0, 1, 0)

	styleSection = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	stylePhase = lipgloss.NewStyle().
			Foreground(colorInfo).
			Bold(true)

	styleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	styleWarning = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)

	styleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleStats = lipgloss.NewStyle().
			Foreground(colorInfo)
)

// ─── Output Functions ───────────────────────────────────────────

// PrintTitle prints a main title
func PrintTitle(text string) {
	fmt.Println(styleTitle.Render(text))
}

// PrintSection prints a section header with border
func PrintSection(text string) {
	fmt.Println(styleSection.Render(text))
}

// PrintPhase prints the start of a compilation phase
func PrintPhase(phaseName string) {
	fmt.Printf("%s %s\n", stylePhase.Render("→"), phaseName)
}

// PrintSuccess prints a success message
func PrintSuccess(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleSuccess.Render("✓"), msg)
}

// PrintError prints an error message
func PrintError(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleError.Render("✗"), msg)
}

// PrintWarning prints a warning message
func PrintWarning(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleWarning.Render("⚠"), msg)
}

// PrintInfo prints an info message
func PrintInfo(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleInfo.Render("ℹ"), msg)
}

// PrintStats prints statistics
func PrintStats(label string, args ...interface{}) {
	msg := fmt.Sprintf(label, args...)
	fmt.Printf("%s %s\n", styleMuted.Render("●"), msg)
}

// PrintDebug prints debug information
func PrintDebug(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s\n", styleMuted.Render(msg))
}

// PrintRaw prints raw text without styling
func PrintRaw(text string, args ...interface{}) {
	fmt.Printf(text, args...)
}

// PrintBox prints text in a styled box
func PrintBox(title, content string) {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		Foreground(colorInfo)

	if title != "" {
		fmt.Println(styleSection.Render(title))
	}
	fmt.Println(box.Render(content))
}

// PrintTable prints a simple table separator
func PrintTableSeparator() {
	fmt.Println(styleMuted.Render("══════════════════════════════════════════════"))
}

var styleInfo = lipgloss.NewStyle().Foreground(colorInfo)
