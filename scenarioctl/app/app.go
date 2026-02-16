package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NavItem struct {
	Name string
	New  func() Page
}

func (i NavItem) Title() string       { return i.Name }
func (i NavItem) Description() string { return "" }
func (i NavItem) FilterValue() string { return i.Name }

type App struct {
	nav      list.Model
	page     Page
	width    int
	height   int
	focusNav bool // true = nav has focus, false = page has focus
}

func New(items []NavItem) App {
	listItems := make([]list.Item, len(items))
	for i := range items {
		listItems[i] = items[i]
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 24, 14)
	l.Title = "Tables"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return App{
		nav:      l,
		page:     items[0].New(),
		focusNav: true, // Start with nav focused
	}
}

func (m App) Init() tea.Cmd {
	return m.page.Init()
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.nav.SetSize(24, msg.Height-2)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Enter only works in nav mode to select a table
			if m.focusNav && !m.page.IsInForm() {
				if it, ok := m.nav.SelectedItem().(NavItem); ok {
					m.page = it.New()
					m.focusNav = false // Move focus to the new page
					return m, m.page.Init()
				}
			}
		case "esc":
			// Esc always goes back one level
			if m.page.IsInForm() {
				// In form - let page handle it (should close form and go back to table view)
				var pageCmd tea.Cmd
				m.page, pageCmd = m.page.Update(msg)
				return m, pageCmd
			} else if !m.focusNav {
				// In table view - go back to nav
				m.focusNav = true
				return m, nil
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Route messages based on mode
	if m.focusNav && !m.page.IsInForm() {
		// Stage 1: Navigation - only nav receives input
		var cmd tea.Cmd
		m.nav, cmd = m.nav.Update(msg)
		return m, cmd
	} else {
		// Stage 2 & 3: Table view or form - only page receives input
		var pageCmd tea.Cmd
		m.page, pageCmd = m.page.Update(msg)
		return m, pageCmd
	}
}

func (m App) View() string {
	sideStyle := lipgloss.NewStyle().
		Width(24).
		Height(m.height).
		Padding(1)
	// Make nav faint when not focused
	if !m.focusNav {
		sideStyle = sideStyle.Faint(true)
	}
	side := sideStyle.Render(m.nav.View())

	bodyStyle := lipgloss.NewStyle().
		Width(max(0, m.width-26)).
		Height(m.height).
		Padding(1, 2)
	// Make page faint when in nav mode
	if m.focusNav {
		bodyStyle = bodyStyle.Faint(true)
	}

	pageView := m.page.View()

	// Add help text based on stage
	if m.page.IsInForm() {
		// Stage 3: In form
		pageView = pageView + "\n\n[Enter] Save | [Esc] Cancel"
	} else if !m.focusNav {
		// Stage 2: Table view
		pageView = pageView + "\n\n[↑↓] Navigate | [n] New | [e] Edit | [d] Delete | [Esc] Back to Tables"
	}

	body := bodyStyle.Render(m.page.Title() + "\n\n" + pageView)

	return lipgloss.JoinHorizontal(lipgloss.Top, side, body)
}
