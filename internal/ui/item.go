package ui

type FrontendItem struct {
	title, value, desc string
}

func (i FrontendItem) Title() string       { return i.title }
func (i FrontendItem) Description() string { return i.desc }
func (i FrontendItem) FilterValue() string { return i.title }

func NewFrontendItem(title, value, desc string) FrontendItem {
	return FrontendItem{title: title, value: value, desc: desc}
}

func (i FrontendItem) Value() string {
	return i.value
}
