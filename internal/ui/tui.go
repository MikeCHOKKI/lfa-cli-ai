package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekimart/lfa-cli-ai/internal/detect"
	"github.com/ekimart/lfa-cli-ai/internal/installer"
)

type state int

const (
	stateInit state = iota
	stateWelcome
	stateModeSelect
	stateDetect
	stateInstallPrompt
	stateInstall
	stateDeploy
	stateDone
	stateError
)

type detectDoneMsg struct {
	ocInstalled     bool
	ocPath          string
	ollamaInstalled bool
	ollamaReachable bool
}

type installDoneMsg struct {
	err error
}

type deployDoneMsg struct {
	err error
}

type model struct {
	state    state
	width    int
	height   int
	autoMode bool
	cursor   int
	err      error

	detectResult detectDoneMsg
	ollamaEnable bool

	spinner spinner.Model
	loading bool
}

func RunDashboard() error {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func initialModel() model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(ColorAccent)
	return model{
		state:    stateInit,
		autoMode: false,
		spinner:  s,
		loading:  false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "up", "k":
			return m.handleUp()
		case "down", "j":
			return m.handleDown()
		}

	case detectDoneMsg:
		m.detectResult = msg
		m.loading = false
		if m.autoMode {
			return m.startInstall()
		}
		m.state = stateInstallPrompt
		return m, nil

	case installDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		if m.autoMode {
			return m.startDeploy()
		}
		return m, startDeployCmd()

	case deployDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		m.state = stateDone
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	return m, nil
}

func (m *model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateInit, stateWelcome:
		m.state = stateModeSelect
		return m, nil
	case stateModeSelect:
		m.autoMode = false
		m.state = stateDetect
		m.loading = true
		return m, startDetectCmd()
	case stateInstallPrompt:
		return m.startInstall()
	case stateDeploy:
		return m.startDeploy()
	case stateDone:
		return m, tea.Quit
	case stateError:
		return m, tea.Quit
	}
	return m, nil
}

func (m *model) handleUp() (tea.Model, tea.Cmd) {
	if m.state == stateModeSelect && m.cursor > 0 {
		m.cursor--
	}
	return m, nil
}

func (m *model) handleDown() (tea.Model, tea.Cmd) {
	if m.state == stateModeSelect && m.cursor < 1 {
		m.cursor++
	}
	return m, nil
}

func (m *model) startInstall() (tea.Model, tea.Cmd) {
	ocInstalled, _ := detect.DetectOpenCode()
	if ocInstalled {
		return m, startDeployCmd()
	}
	m.state = stateInstall
	m.loading = true
	return m, startInstallCmd()
}

func (m *model) startDeploy() (tea.Model, tea.Cmd) {
	m.state = stateDeploy
	m.loading = true
	return m, startDeployCmd()
}

func (m model) View() string {
	switch m.state {
	case stateInit, stateWelcome:
		return m.welcomeView()
	case stateModeSelect:
		return m.modeSelectView()
	case stateDetect:
		return m.detectView()
	case stateInstallPrompt:
		return m.installPromptView()
	case stateInstall:
		return m.installView()
	case stateDeploy:
		return m.deployView()
	case stateDone:
		return m.doneView()
	case stateError:
		return m.errorView()
	}
	return ""
}

func (m model) welcomeView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("  ⚡ LFA CLI"))
	b.WriteString("\n\n")
	b.WriteString(SubtitleStyle.Render("  OpenCode AI Configuration Tool"))
	b.WriteString("\n\n\n")
	b.WriteString(DimStyle.Render("  Press Enter to continue..."))
	return AppStyle.Render(b.String())
}

func (m model) modeSelectView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" Select Mode"))
	b.WriteString("\n\n")

	cursor := "  "
	selected := "▸ "
	choices := []string{"🚀  Auto (full automatic)", "🎮  Guided (step by step)"}

	for i, choice := range choices {
		if m.cursor == i {
			b.WriteString(AccentStyle.Render(selected + choice))
		} else {
			b.WriteString(DimStyle.Render(cursor + choice))
		}
		b.WriteString("\n\n")
	}

	b.WriteString(DimStyle.Render("\n  ↑/↓ navigate • Enter select • q quit"))
	return AppStyle.Render(b.String())
}

func (m model) detectView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" Detection"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(fmt.Sprintf("  %s %s", m.spinner.View(), DimStyle.Render("Scanning system...")))
	} else {
		check := SuccessStyle.Render("✓")
		cross := ErrorStyle.Render("✗")

		b.WriteString(fmt.Sprintf("  OS:        %s\n", AccentStyle.Render(detect.DetectOS().String())))

		if m.detectResult.ocInstalled {
			b.WriteString(fmt.Sprintf("  OpenCode:  %s %s\n", check, DimStyle.Render(m.detectResult.ocPath)))
		} else {
			b.WriteString(fmt.Sprintf("  OpenCode:  %s %s\n", cross, WarningStyle.Render("not installed")))
		}

		if m.detectResult.ollamaInstalled {
			b.WriteString(fmt.Sprintf("  Ollama:    %s %s\n", check, DimStyle.Render("installed")))
		} else {
			b.WriteString(fmt.Sprintf("  Ollama:    %s %s\n", WarningStyle.Render("~"), DimStyle.Render("not found")))
		}

		if m.detectResult.ollamaReachable {
			b.WriteString(fmt.Sprintf("  Ollama API: %s %s\n", check, DimStyle.Render("reachable")))
		} else {
			b.WriteString(fmt.Sprintf("  Ollama API: %s %s\n", cross, DimStyle.Render("unreachable")))
		}

		b.WriteString(DimStyle.Render("\n  Press Enter to continue..."))
	}

	return AppStyle.Render(b.String())
}

func (m model) installPromptView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" OpenCode"))
	b.WriteString("\n\n")

	if m.detectResult.ocInstalled {
		b.WriteString(SuccessStyle.Render("  ✓ OpenCode is already installed"))
		b.WriteString(DimStyle.Render("\n  Press Enter to deploy configuration..."))
	} else {
		b.WriteString(WarningStyle.Render("  OpenCode is not installed"))
		b.WriteString("\n\n")
		b.WriteString(TextStyle.Render("  The installer will download and install OpenCode"))
		b.WriteString(DimStyle.Render("\n  Press Enter to continue..."))
	}

	return AppStyle.Render(b.String())
}

func (m model) installView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" Installing OpenCode"))
	b.WriteString("\n\n")
	if m.loading {
		b.WriteString(fmt.Sprintf("  %s %s", m.spinner.View(), DimStyle.Render("Downloading and installing...")))
	} else {
		b.WriteString(SuccessStyle.Render("  ✓ Installation complete"))
	}
	return AppStyle.Render(b.String())
}

func (m model) deployView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" Deployment"))
	b.WriteString("\n\n")
	if m.loading {
		b.WriteString(fmt.Sprintf("  %s %s", m.spinner.View(), DimStyle.Render("Deploying configuration...")))
	} else {
		b.WriteString(SuccessStyle.Render("  ✓ Deployment complete"))
	}
	return AppStyle.Render(b.String())
}

func (m model) doneView() string {
	var b strings.Builder
	b.WriteString(SuccessStyle.Render("  ✓ Setup Complete!"))
	b.WriteString("\n\n")
	b.WriteString(TextStyle.Render("  OpenCode is configured and ready."))
	b.WriteString("\n\n")
	b.WriteString(DimStyle.Render("  Press Enter to quit..."))
	return AppStyle.Render(b.String())
}

func (m model) errorView() string {
	var b strings.Builder
	b.WriteString(ErrorStyle.Render("  ✗ Error"))
	b.WriteString("\n\n")
	b.WriteString(TextStyle.Render(fmt.Sprintf("  %v", m.err)))
	b.WriteString(DimStyle.Render("\n\n  Press Enter to quit..."))
	return AppStyle.Render(b.String())
}

func startDetectCmd() tea.Cmd {
	return func() tea.Msg {
		ocInstalled, ocPath := detect.DetectOpenCode()
		ollamaInstalled, ollamaReachable := detect.DetectOllama()
		return detectDoneMsg{
			ocInstalled:     ocInstalled,
			ocPath:          ocPath,
			ollamaInstalled: ollamaInstalled,
			ollamaReachable: ollamaReachable,
		}
	}
}

func startInstallCmd() tea.Cmd {
	return func() tea.Msg {
		o := detect.DetectOS()
		err := installer.InstallOpenCode(o, installer.DefaultVersion)
		return installDoneMsg{err: err}
	}
}

func startDeployCmd() tea.Cmd {
	return func() tea.Msg {
		o := detect.DetectOS()
		err := installer.DeployConfig(o, "data", true)
		return deployDoneMsg{err: err}
	}
}
