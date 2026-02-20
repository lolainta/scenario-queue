package app

import tea "github.com/charmbracelet/bubbletea"

type Page interface {
	Init() tea.Cmd
	Update(tea.Msg) (Page, tea.Cmd)
	View() string
	Title() string
	IsInForm() bool // Returns true if page is in edit/form mode and needs all keyboard input
}
