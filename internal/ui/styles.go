package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	ColorBg      = lipgloss.Color("#0D1117")
	ColorSurface = lipgloss.Color("#161B22")
	ColorAccent  = lipgloss.Color("#58A6FF")
	ColorSuccess = lipgloss.Color("#2EA043")
	ColorWarning = lipgloss.Color("#D29922")
	ColorText    = lipgloss.Color("#C9D1D9")
	ColorError   = lipgloss.Color("#FF6B6B")
)

var (
	AppStyle    = lipgloss.NewStyle().Padding(1, 2)
	TitleStyle  = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	SubtitleStyle = lipgloss.NewStyle().Foreground(ColorText).Faint(true)
	SuccessStyle  = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorWarning)
	TextStyle     = lipgloss.NewStyle().Foreground(ColorText)
	SurfaceStyle  = lipgloss.NewStyle().Background(ColorSurface).Padding(0, 1)
	AccentStyle   = lipgloss.NewStyle().Foreground(ColorAccent)
	DimStyle      = lipgloss.NewStyle().Foreground(ColorText).Faint(true)

	LogoStyle = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Margin(0, 0, 1, 0)
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	ReportCaller:    false,
})

func init() {
	Logger.SetStyles(
		&log.Styles{
			Levels: map[log.Level]lipgloss.Style{
				log.DebugLevel: lipgloss.NewStyle().Foreground(ColorText),
				log.InfoLevel:  lipgloss.NewStyle().Foreground(ColorAccent),
				log.WarnLevel:  lipgloss.NewStyle().Foreground(ColorWarning),
				log.ErrorLevel: lipgloss.NewStyle().Foreground(ColorError),
			},
			Key:     lipgloss.NewStyle().Foreground(ColorAccent),
			Value:   lipgloss.NewStyle().Foreground(ColorText),
			Message: lipgloss.NewStyle().Foreground(ColorText),
		},
	)
}
