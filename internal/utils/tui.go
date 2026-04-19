package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorPrimary = lipgloss.Color("33")
	colorSuccess = lipgloss.Color("34")
	colorWarning = lipgloss.Color("208")
	colorError   = lipgloss.Color("196")
	colorInfo    = lipgloss.Color("39")
	colorMuted   = lipgloss.Color("243")

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

func PrintTitle(text string) {
	fmt.Println(styleTitle.Render(text))
}

func PrintSection(text string) {
	fmt.Println(styleSection.Render(text))
}

func PrintPhase(phaseName string) {
	fmt.Printf("%s %s\n", stylePhase.Render("→"), phaseName)
}

func PrintSuccess(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleSuccess.Render("✓"), msg)
}

func PrintError(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleError.Render("✗"), msg)
}

func PrintWarning(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleWarning.Render("⚠"), msg)
}

func PrintInfo(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s %s\n", styleInfo.Render("ℹ"), msg)
}

func PrintStats(label string, args ...interface{}) {
	msg := fmt.Sprintf(label, args...)
	fmt.Printf("%s %s\n", styleMuted.Render("●"), msg)
}

func PrintDebug(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	fmt.Printf("%s\n", styleMuted.Render(msg))
}

func PrintRaw(text string, args ...interface{}) {
	fmt.Printf(text, args...)
}

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

func PrintTableSeparator() {
	fmt.Println(styleMuted.Render("══════════════════════════════════════════════"))
}

var styleInfo = lipgloss.NewStyle().Foreground(colorInfo)
