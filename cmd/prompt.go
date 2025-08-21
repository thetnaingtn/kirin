package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	stepAppName = iota
	stepModuleName
	stepFrontendLibrary
	stepStartScaffolding
)

type frontendItem struct {
	title, value, desc string
}

func (i frontendItem) Title() string       { return i.title }
func (i frontendItem) Description() string { return i.desc }
func (i frontendItem) FilterValue() string { return i.title }

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
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

func NewPrompt() *Prompt {
	// Initialize spinner
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.Spinner = spinner.MiniDot
	// Create app name input
	appInput := textinput.New()
	appInput.Placeholder = "Enter your app name..."
	appInput.Focus()
	appInput.CharLimit = 50
	appInput.Width = 40

	// Create module input
	moduleInput := textinput.New()
	moduleInput.Placeholder = "Enter your module name (e.g., github.com/user/app)..."
	moduleInput.CharLimit = 100
	moduleInput.Width = 50

	// Create frontend library list
	items := []list.Item{
		frontendItem{title: "Vue", value: "vue", desc: "Progressive JavaScript framework"},
		frontendItem{title: "React", value: "react", desc: "JavaScript library for building user interfaces"},
		frontendItem{title: "Svelte", value: "svelte", desc: "Cybernetically enhanced web apps"},
	}

	frontendList := list.New(items, list.NewDefaultDelegate(), 60, 10)
	frontendList.Title = "Choose Frontend Library"
	frontendList.SetShowStatusBar(false)
	frontendList.SetFilteringEnabled(false)
	frontendList.Styles.Title = titleStyle

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
				if selectedItem, ok := p.frontendList.SelectedItem().(frontendItem); ok {
					p.frontendChoice = selectedItem.value
					cmds = append(cmds, startScaffolding, p.createProject())
				}
			}
		}
	case scaffoldStart:
		p.Next()
		cmds = append(cmds, p.spinner.Tick)

	case scaffoldFinish:
		p.quitting = true

	case tea.WindowSizeMsg:
		h, v := inputStyle.GetFrameSize()
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
		b.WriteString(fmt.Sprintf("Create new application project in %s (module %s)\n\n", p.appName, p.moduleName))
		b.WriteString(fmt.Sprintf("✨ Done in %s. Press q to quit.\n", p.latency))

		return b.String()
	}

	switch p.step {
	case stepAppName:
		b.WriteString(titleStyle.Render("Step 1/3: App Name"))
		b.WriteString("\n\n")
		b.WriteString("What would you like to name your app?\n\n")
		b.WriteString(inputStyle.Render(p.appNameInput.View()))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Enter to continue • Press q to quit"))

	case stepModuleName:
		b.WriteString(titleStyle.Render("Step 2/3: Module Name"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("App Name: %s\n", p.appName))
		b.WriteString("What's your Go module name?\n\n")
		b.WriteString(inputStyle.Render(p.moduleInput.View()))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Enter to continue • Press q to quit"))

	case stepFrontendLibrary:
		b.WriteString(titleStyle.Render("Step 3/3: Frontend Library"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("App Name: %s\n", p.appName))
		b.WriteString(fmt.Sprintf("Module Name: %s\n", p.moduleName))
		b.WriteString("Choose your frontend library:\n\n")
		b.WriteString(p.frontendList.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press Enter to select • Press q to quit"))

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

type scaffoldFinish struct{}
type scaffoldStart struct{}

func (p *Prompt) createProject() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		appName := p.appName
		modName := appName
		if p.moduleName != "" {
			modName = p.moduleName
		}

		wd, _ := os.Getwd()
		projectPath := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, appName)

		_ = os.Mkdir(projectPath, 0750)

		git, _ := exec.LookPath("git")

		c := exec.Command(git, "clone", "-b", fmt.Sprintf("frontend/%s", p.frontendChoice), cloneUrl, projectPath)

		_ = c.Run()

		_ = replace(projectPath, "go.mod", "bolierplate", modName)

		_ = replace(projectPath, "*.go", "bolierplate", modName)

		p.latency = time.Since(start)

		return scaffoldFinish{}
	}
}

func startScaffolding() tea.Msg {
	return scaffoldStart{}
}
