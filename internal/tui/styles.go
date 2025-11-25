package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	primaryColor   = lipgloss.Color("#FF6B35") // MLB Orange
	secondaryColor = lipgloss.Color("#004687") // MLB Blue
	accentColor    = lipgloss.Color("#FFFFFF")
	mutedColor     = lipgloss.Color("#626262")
	successColor   = lipgloss.Color("#04B575")
	errorColor     = lipgloss.Color("#FF4136")
)

// Styles
var (
	// App frame
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title bar
	TitleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Background(secondaryColor).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	// Header for sections
	HeaderStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	// Selected item in list
	SelectedStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Background(secondaryColor).
			Bold(true).
			Padding(0, 1)

	// Normal item in list
	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	// Muted/secondary text
	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Status bar at bottom
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Breadcrumb
	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			MarginBottom(1)

	// Table header
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(mutedColor)

	// Box/border style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Background(secondaryColor).
			Bold(true).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(lipgloss.Color("#1a1a1a")).
				Padding(0, 2)

	// Filter input
	FilterStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	// Success message
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Error message
	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
)
