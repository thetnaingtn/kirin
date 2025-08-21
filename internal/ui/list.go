package ui

import (
	"github.com/charmbracelet/bubbles/list"
)

func NewFrontendList() list.Model {

	items := []list.Item{
		NewFrontendItem("Vue", "vue", "Progressive JavaScript framework"),
		NewFrontendItem("React", "react", "JavaScript library for building user interfaces"),
		NewFrontendItem("Svelte", "svelte", "Cybernetically enhanced web apps"),
	}

	frontendList := list.New(items, list.NewDefaultDelegate(), 60, 10)
	frontendList.Title = "Choose Frontend Library"
	frontendList.SetShowStatusBar(false)
	frontendList.SetFilteringEnabled(false)
	frontendList.Styles.Title = TitleStyle

	return frontendList
}
