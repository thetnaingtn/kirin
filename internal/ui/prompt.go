package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thetnaingtn/kirin/internal/kirin"
)

const (
	stepAppName = iota
	stepModuleName
	stepFrontendLibrary
	stepStartScaffolding
)

type Prompt struct {
	latency        time.Duration
	step           int
	appNameInput   textinput.Model
	moduleInput    textinput.Model
	frontendList   list.Model
	appName        string
	moduleName     string
	frontendChoice string
	quitting       bool
	spinner        spinner.Model
	err            error
}

func NewPrompt() *Prompt {
	// Initialize spinner
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.Spinner = spinner.MiniDot
	// Create app name input
	appInput := NewAppNameInput()

	// Create module input
	moduleInput := NewModuleInput()

	// Create frontend library list

	frontendList := NewFrontendList()

	return &Prompt{
		step:         stepAppName,
		appNameInput: appInput,
		moduleInput:  moduleInput,
		frontendList: frontendList,
		spinner:      s,
	}
}

func (p *Prompt) Next() {
	if p.step < stepStartScaffolding {
		p.step++
	}
}

func (p *Prompt) Init() tea.Cmd {
	return textinput.Blink
}

func (p *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			p.quitting = true
			return p, tea.Quit

		case "enter":
			switch p.step {
			case stepAppName:
				if p.appNameInput.Value() != "" {
					p.appName = p.appNameInput.Value()
					p.moduleInput.Focus()
					p.Next()
				}
			case stepModuleName:
				if p.moduleInput.Value() != "" {
					p.moduleName = p.moduleInput.Value()
					p.Next()
				}
			case stepFrontendLibrary:
				if selectedItem, ok := p.frontendList.SelectedItem().(FrontendItem); ok {
					p.frontendChoice = selectedItem.Value()
					cmds = append(cmds, startScaffolding, p.createProject())
				}
			}
		}
	case scaffoldStart:
		p.Next()
		cmds = append(cmds, p.spinner.Tick)

	case scaffoldFinish:
		p.quitting = true
	case scaffoldError:
		p.err = msg.err

	case tea.WindowSizeMsg:
		h, v := InputStyle.GetFrameSize()
		p.appNameInput.Width = msg.Width - h
		p.moduleInput.Width = msg.Width - h
		p.frontendList.SetWidth(msg.Width - h)
		p.frontendList.SetHeight(msg.Height - v)
	}

	// Update the appropriate model based on current step
	switch p.step {
	case stepAppName:
		p.appNameInput, cmd = p.appNameInput.Update(msg)
	case stepModuleName:
		p.moduleInput, cmd = p.moduleInput.Update(msg)
	case stepFrontendLibrary:
		p.frontendList, cmd = p.frontendList.Update(msg)
	case stepStartScaffolding:
		p.spinner, cmd = p.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *Prompt) View() string {
	var b strings.Builder

	if p.quitting {
		b.WriteString(fmt.Sprintf("Create new application project named: %s (module %s)\n\n", p.appName, p.moduleName))
		b.WriteString(fmt.Sprintf("✨ Done in %s. Press q to quit.\n", p.latency))

		return b.String()
	}

	if p.err != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("Can't create project at the moment: %s", p.err.Error())))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Press q to quit"))
		return b.String()
	}

	switch p.step {
	case stepAppName:
		b.WriteString(TitleStyle.Render("Step 1/3: App Name"))
		b.WriteString("\n\n")
		b.WriteString("What would you like to name your app?\n\n")
		b.WriteString(InputStyle.Render(p.appNameInput.View()))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Press Enter to continue • Press q to quit"))

	case stepModuleName:
		b.WriteString(TitleStyle.Render("Step 2/3: Module Name"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("App Name: %s\n", p.appName))
		b.WriteString("What's your Go module name?\n\n")
		b.WriteString(InputStyle.Render(p.moduleInput.View()))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Press Enter to continue • Press q to quit"))

	case stepFrontendLibrary:
		b.WriteString(TitleStyle.Render("Step 3/3: Frontend Library"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("App Name: %s\n", p.appName))
		b.WriteString(fmt.Sprintf("Module Name: %s\n", p.moduleName))
		b.WriteString("Choose your frontend library:\n\n")
		b.WriteString(p.frontendList.View())
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("Press Enter to select • Press q to quit"))

	case stepStartScaffolding:
		b.WriteString("Here's your configuration:\n\n")
		b.WriteString(fmt.Sprintf("📱 App Name: %s\n", p.appName))
		b.WriteString(fmt.Sprintf("📦 Module Name: %s\n", p.moduleName))
		b.WriteString(fmt.Sprintf("⚡ Frontend Library: %s\n", p.frontendChoice))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%s Scaffolding your project...\n\n", p.spinner.View()))
	}

	return b.String()
}

// GetResults returns the collected form data
func (p *Prompt) GetResults() (appName, moduleName, frontendChoice string) {
	return p.appName, p.moduleName, p.frontendChoice
}

func (p *Prompt) createProject() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		appName := p.appName
		modName := appName
		if p.moduleName != "" {
			modName = p.moduleName
		}

		if err := kirin.CreateProject(appName, modName, p.frontendChoice); err != nil {
			return scaffoldError{err: err}
		}

		p.latency = kirin.FormatTime(time.Since(start))

		return scaffoldFinish{}
	}
}
