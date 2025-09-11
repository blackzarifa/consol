package tui

import "github.com/charmbracelet/lipgloss"

const (
	ColorDarkGray    = "246"
	ColorLightGray   = "240"
	ColorYellow      = "136"
	ColorDarkGreen   = "28"
	ColorDarkBlue    = "19"
	ColorGreen       = "22"
	ColorBlue        = "18"
	ColorWhite       = "15"
	ColorBlack       = "0"
	ColorBrightWhite = "255"
)

var (
	HeaderTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Border(lipgloss.NormalBorder()).
				Padding(0, 1)

	ConflictMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Dark: ColorDarkGray, Light: ColorLightGray})

	ResolvedLineNumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorYellow)).
				BorderRight(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(ColorYellow))

	CursorStyle = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Dark: ColorBrightWhite, Light: ColorBlack}).
			Foreground(lipgloss.AdaptiveColor{Dark: ColorBlack, Light: ColorBrightWhite}).
			Blink(true)

	OursBranchStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorDarkGreen)).
			Foreground(lipgloss.Color(ColorWhite)).
			Bold(true)

	TheirsBranchStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(ColorDarkBlue)).
				Foreground(lipgloss.Color(ColorWhite)).
				Bold(true)

	OursContentStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(ColorGreen)).
				Foreground(lipgloss.Color(ColorWhite))

	TheirsContentStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(ColorBlue)).
				Foreground(lipgloss.Color(ColorWhite))
)
