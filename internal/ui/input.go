package ui

import "github.com/charmbracelet/bubbles/textinput"

func NewAppNameInput() textinput.Model {
	appInput := textinput.New()
	appInput.Placeholder = "Enter your app name..."
	appInput.Focus()
	appInput.CharLimit = 50
	appInput.Width = 40
	return appInput
}
func NewModuleInput() textinput.Model {
	moduleInput := textinput.New()
	moduleInput.Placeholder = "Enter your module name (e.g., github.com/user/app)..."
	moduleInput.CharLimit = 100
	moduleInput.Width = 50
	return moduleInput
}
