package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lfa-cli/lfa-cli-ai/internal/config"
	"github.com/lfa-cli/lfa-cli-ai/internal/detect"
	"github.com/lfa-cli/lfa-cli-ai/internal/installer"
)

type state int

const (
	stateInit state = iota
	stateWelcome
	stateModeSelect
	stateDetect
	stateInstallPrompt
	stateInstall
	statePGSetup
	stateDeploy
	stateDone
	stateError
)

type detectDoneMsg struct {
	ocInstalled     bool
	ocPath          string
	ollamaInstalled bool
	ollamaReachable bool
	pgInstalled     bool
	pgRunning       bool
	pgVersion       string
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
	pgHost       string
	pgPort       string
	pgUser       string
	pgPassword   string
	pgDBName     string
	pgInputField int // 0=host, 1=port, 2=user, 3=password, 4=dbname

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
		case "backspace":
			if m.state == statePGSetup {
				return m.handlePGDelete()
			}
		default:
			if m.state == statePGSetup {
				return m.handlePGInput(msg.String())
			}
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
		return m.startPGSetup()

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
	case statePGSetup:
		return m.startDeploy()
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
	switch m.state {
	case stateModeSelect:
		if m.cursor > 0 {
			m.cursor--
		}
	case statePGSetup:
		if m.pgInputField > 0 {
			m.pgInputField--
		}
	}
	return m, nil
}

func (m *model) handleDown() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateModeSelect:
		if m.cursor < 1 {
			m.cursor++
		}
	case statePGSetup:
		if m.pgInputField < 4 {
			m.pgInputField++
		}
	}
	return m, nil
}

func (m *model) handlePGDelete() (tea.Model, tea.Cmd) {
	switch m.pgInputField {
	case 0:
		if len(m.pgHost) > 0 { m.pgHost = m.pgHost[:len(m.pgHost)-1] }
	case 1:
		if len(m.pgPort) > 0 { m.pgPort = m.pgPort[:len(m.pgPort)-1] }
	case 2:
		if len(m.pgUser) > 0 { m.pgUser = m.pgUser[:len(m.pgUser)-1] }
	case 3:
		if len(m.pgPassword) > 0 { m.pgPassword = m.pgPassword[:len(m.pgPassword)-1] }
	case 4:
		if len(m.pgDBName) > 0 { m.pgDBName = m.pgDBName[:len(m.pgDBName)-1] }
	}
	return m, nil
}

func (m *model) handlePGInput(key string) (tea.Model, tea.Cmd) {
	if len(key) != 1 || key[0] < 32 || key[0] > 126 {
		return m, nil
	}
	switch m.pgInputField {
	case 0:
		m.pgHost += key
	case 1:
		m.pgPort += key
	case 2:
		m.pgUser += key
	case 3:
		m.pgPassword += key
	case 4:
		m.pgDBName += key
	}
	return m, nil
}

func (m *model) startInstall() (tea.Model, tea.Cmd) {
	ocInstalled, _ := detect.DetectOpenCode()
	if ocInstalled {
		return m.startPGSetup()
	}
	m.state = stateInstall
	m.loading = true
	return m, startInstallCmd()
}

func (m *model) startPGSetup() (tea.Model, tea.Cmd) {
	if m.autoMode {
		m.pgHost = "localhost"
		m.pgPort = "5432"
		m.pgUser = "postgres"
		m.pgPassword = "postgres"
		m.pgDBName = "opencode_db"
		return m, startDeployCmdWithPG(m.pgHost, m.pgPort, m.pgUser, m.pgPassword, m.pgDBName)
	}
	m.state = statePGSetup
	m.pgInputField = 0
	m.pgHost = "localhost"
	m.pgPort = "5432"
	m.pgUser = "postgres"
	m.pgPassword = "postgres"
	m.pgDBName = "opencode_db"
	return m, nil
}

func (m *model) startDeploy() (tea.Model, tea.Cmd) {
	m.state = stateDeploy
	m.loading = true
	return m, startDeployCmdWithPG(m.pgHost, m.pgPort, m.pgUser, m.pgPassword, m.pgDBName)
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
	case statePGSetup:
		return m.pgSetupView()
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

		if m.detectResult.pgInstalled {
			b.WriteString(fmt.Sprintf("  PostgreSQL: %s %s", check, DimStyle.Render("installed")))
			if m.detectResult.pgRunning {
				b.WriteString(fmt.Sprintf(" %s\n", SuccessStyle.Render("(running)")))
			} else {
				b.WriteString(fmt.Sprintf(" %s\n", WarningStyle.Render("(stopped)")))
			}
		} else {
			b.WriteString(fmt.Sprintf("  PostgreSQL: %s %s\n", WarningStyle.Render("~"), DimStyle.Render("not installed")))
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

func (m model) pgSetupView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" PostgreSQL Configuration"))
	b.WriteString("\n\n")
	b.WriteString(TextStyle.Render("  Configure your PostgreSQL connection:"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		value string
	}{
		{"Host", m.pgHost},
		{"Port", m.pgPort},
		{"User", m.pgUser},
		{"Password", m.pgPassword},
		{"Database", m.pgDBName},
	}

	for i, f := range fields {
		cursor := "  "
		if m.pgInputField == i {
			cursor = "▸ "
			b.WriteString(AccentStyle.Render(cursor + f.label + ": "))
			b.WriteString(AccentStyle.Render(f.value + "█"))
		} else {
			b.WriteString(DimStyle.Render(cursor + f.label + ": "))
			b.WriteString(TextStyle.Render(f.value))
		}
		b.WriteString("\n\n")
	}

	b.WriteString(DimStyle.Render("\n  ↑/↓ navigate • Enter confirm • Type to edit • q quit"))
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
		pgStatus := detect.DetectPostgreSQL()
		return detectDoneMsg{
			ocInstalled:     ocInstalled,
			ocPath:          ocPath,
			ollamaInstalled: ollamaInstalled,
			ollamaReachable: ollamaReachable,
			pgInstalled:     pgStatus.Installed,
			pgRunning:       pgStatus.Running,
			pgVersion:       pgStatus.Version,
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

func startDeployCmdWithPG(host, port, user, password, dbname string) tea.Cmd {
	return func() tea.Msg {
		o := detect.DetectOS()
		configDir := detect.GetOpenCodeConfigDir(o)

		if err := installer.DeployConfig(o, true); err != nil {
			return deployDoneMsg{err: err}
		}

		if host == "" {
			return deployDoneMsg{err: nil}
		}

		conn := config.PGConnection{Host: host, Port: port, User: user, Password: password, DBName: dbname}

		if dplyErr := installer.DeployPGMcpServer(configDir, conn); dplyErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: PG MCP deploy: %v\n", dplyErr)
		}

		// Initialize PostgreSQL schema
		if err := installer.InitPGSchema(conn); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: PG schema init: %v\n", err)
		}

		cfgPath := config.GetConfigPath(o)
		cfg, cfgErr := config.LoadExistingConfig(cfgPath)
		if cfgErr == nil {
			config.LinkPostgreSQL(cfg, o, conn.Host, conn.Port, conn.User, conn.Password, conn.DBName)

			// Also pick up tokens from env if present
			if val := os.Getenv("GITHUB_TOKEN"); val != "" {
				config.LinkToken(cfg, "github", "GITHUB_PERSONAL_ACCESS_TOKEN", val)
			}

			if wErr := config.WriteConfig(cfgPath, cfg); wErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: write PG config: %v\n", wErr)
			}
		}

		return deployDoneMsg{err: nil}
	}
}
