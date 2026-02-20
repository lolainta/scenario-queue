package pages

import (
	"context"

	"scenarioctl/app"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Loader func(context.Context) ([]table.Row, error)

// FieldType defines the type of form field
type FieldType string

const (
	FieldTypeText   FieldType = "text"
	FieldTypeSelect FieldType = "select"
)

type SelectOption struct {
	Label string
	Value string
}

// FieldDef defines a form field
type FieldDef struct {
	Label   string
	Type    FieldType
	Options []SelectOption
}

// CRUDCallbacks defines optional callbacks for CRUD operations
type CRUDCallbacks struct {
	OnCreate func() error
	OnUpdate func(rowIndex int) error
	OnDelete func(rowIndex int) error
}

type TablePage struct {
	title       string
	loader      Loader
	table       table.Model
	crud        *CRUDCallbacks
	mode        string // "view", "create", "edit"
	form        *CRUDForm
	err         string
	currentRows []table.Row
}

type CRUDForm struct {
	fields     []textinput.Model
	fieldDefs  []FieldDef
	labels     []string
	focusIndex int
	rowIndex   int
	onSubmit   func(values []string) error
	// For select fields
	selectIndex map[int]int // field index -> selected option index
}

type loadedMsg struct {
	rows []table.Row
	err  error
}

type crudSubmitMsg struct {
	err error
}

func NewTablePage(title string, cols []table.Column, loader Loader) app.Page {
	t := table.New(table.WithColumns(cols))
	t.Focus() // Focus table by default so it can handle keyboard input
	return &TablePage{
		title:  title,
		loader: loader,
		table:  t,
		mode:   "view",
	}
}

// WithCRUD adds CRUD capability to the table page
func (m *TablePage) WithCRUD(callbacks *CRUDCallbacks) *TablePage {
	m.crud = callbacks
	return m
}

// StartForm initiates a form for create/edit operations (backward compatible)
func (m *TablePage) StartForm(numFields int, labels []string, rowIndex int, onSubmit func([]string) error) {
	defs := make([]FieldDef, numFields)
	for i, label := range labels {
		defs[i] = FieldDef{Label: label, Type: FieldTypeText}
	}
	m.StartFormWithDefs(defs, rowIndex, onSubmit)
}

// StartFormWithDefs initiates a form with field definitions (supports select fields)
func (m *TablePage) StartFormWithDefs(fieldDefs []FieldDef, rowIndex int, onSubmit func([]string) error) {
	fields := make([]textinput.Model, 0)
	labels := make([]string, 0)
	selectIndex := make(map[int]int)
	fieldIndex := 0

	for _, def := range fieldDefs {
		if def.Type == FieldTypeText {
			ti := textinput.New()
			ti.CharLimit = 256
			ti.Placeholder = def.Label
			fields = append(fields, ti)
			labels = append(labels, def.Label)
		} else if def.Type == FieldTypeSelect {
			// For select fields, we still need a textinput placeholder
			ti := textinput.New()
			ti.CharLimit = 256
			ti.Placeholder = def.Label
			if len(def.Options) > 0 {
				ti.SetValue(def.Options[0].Label)
			}
			fields = append(fields, ti)
			labels = append(labels, def.Label)
			selectIndex[fieldIndex] = 0 // Start at first option
		}
		fieldIndex++
	}

	if len(fields) > 0 {
		fields[0].Focus()
	}

	m.form = &CRUDForm{
		fields:      fields,
		fieldDefs:   fieldDefs,
		labels:      labels,
		focusIndex:  0,
		rowIndex:    rowIndex,
		onSubmit:    onSubmit,
		selectIndex: selectIndex,
	}
	m.mode = "edit"
}

func (m *TablePage) Init() tea.Cmd {
	return m.load()
}

func (m *TablePage) load() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.loader(context.Background())
		return loadedMsg{rows: rows, err: err}
	}
}

func (m *TablePage) Update(msg tea.Msg) (app.Page, tea.Cmd) {
	switch msg := msg.(type) {

	case loadedMsg:
		if msg.err == nil {
			m.table.SetRows(msg.rows)
			m.currentRows = msg.rows
			m.table.Focus() // Make sure table has focus to handle keyboard input
		}
		return m, nil

	case crudSubmitMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			// Stay in edit mode so user can see error and try again
			return m, nil
		}
		m.mode = "view"
		m.form = nil
		m.err = ""
		m.table.Focus() // Re-focus table when returning to view mode
		return m, m.load()

	case tea.KeyMsg:
		if m.mode == "view" {
			switch msg.String() {
			case "r":
				return m, m.load()
			case "n":
				if m.crud != nil {
					m.mode = "create"
					m.table.Blur() // Blur table when entering form
					if m.crud.OnCreate != nil {
						m.crud.OnCreate()
					}
					// Return Blink command from first focused field
					if m.form != nil && len(m.form.fields) > 0 {
						return m, textinput.Blink
					}
				}
				return m, nil
			case "e":
				if m.table.Cursor() < 0 {
					m.err = "No row selected for editing"
					return m, nil
				}
				if m.crud != nil && m.table.Cursor() < len(m.currentRows) {
					m.table.Blur() // Blur table when entering form
					m.crud.OnUpdate(m.table.Cursor())
					// Pre-populate form fields with current row values (skip ID column at index 0)
					if m.form != nil {
						currentRow := m.currentRows[m.table.Cursor()]
						for i := 1; i < len(currentRow) && (i-1) < len(m.form.fields); i++ {
							fieldIdx := i - 1
							if m.form.fieldDefs[fieldIdx].Type == FieldTypeSelect {
								options := m.form.fieldDefs[fieldIdx].Options
								selectedIndex := 0
								for j := 0; j < len(options); j++ {
									if options[j].Label == currentRow[i] || options[j].Value == currentRow[i] {
										selectedIndex = j
										break
									}
								}
								m.form.selectIndex[fieldIdx] = selectedIndex
								if selectedIndex < len(options) {
									m.form.fields[fieldIdx].SetValue(options[selectedIndex].Label)
								}
							} else {
								m.form.fields[fieldIdx].SetValue(currentRow[i])
							}
						}
					}
					// Return Blink command from first focused field
					if m.form != nil && len(m.form.fields) > 0 {
						return m, textinput.Blink
					}
				}
				return m, nil
			case "d":
				if m.table.Cursor() < 0 {
					m.err = "No row selected for deletion"
					return m, nil
				}
				if m.crud != nil && m.table.Cursor() < len(m.currentRows) {
					cursor := m.table.Cursor()
					if m.crud.OnDelete != nil {
						if err := m.crud.OnDelete(cursor); err != nil {
							m.err = err.Error()
						} else {
							return m, m.load()
						}
					}
				}
				return m, nil
			default:
				// Let table handle arrow keys and other navigation
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		} else if m.mode == "edit" && m.form != nil {
			switch msg.String() {
			case "esc":
				m.mode = "view"
				m.form = nil
				m.err = ""
				m.table.Focus() // Re-focus table when exiting form

				return m, nil
			case "tab":
				m.form.focusIndex = (m.form.focusIndex + 1) % len(m.form.fields)
				m.updateFormFocus()
				return m, nil
			case "shift+tab":
				m.form.focusIndex = (m.form.focusIndex - 1 + len(m.form.fields)) % len(m.form.fields)
				m.updateFormFocus()
				return m, nil
			case "up":
				// Handle Up arrow for select fields
				if m.form.focusIndex < len(m.form.fieldDefs) && m.form.fieldDefs[m.form.focusIndex].Type == FieldTypeSelect {
					options := m.form.fieldDefs[m.form.focusIndex].Options
					if len(options) > 0 {
						if idx, ok := m.form.selectIndex[m.form.focusIndex]; ok {
							idx = (idx - 1 + len(options)) % len(options)
							m.form.selectIndex[m.form.focusIndex] = idx
							m.form.fields[m.form.focusIndex].SetValue(options[idx].Label)
						}
					}
				}
				return m, nil
			case "down":
				// Handle Down arrow for select fields
				if m.form.focusIndex < len(m.form.fieldDefs) && m.form.fieldDefs[m.form.focusIndex].Type == FieldTypeSelect {
					options := m.form.fieldDefs[m.form.focusIndex].Options
					if len(options) > 0 {
						if idx, ok := m.form.selectIndex[m.form.focusIndex]; ok {
							idx = (idx + 1) % len(options)
							m.form.selectIndex[m.form.focusIndex] = idx
							m.form.fields[m.form.focusIndex].SetValue(options[idx].Label)
						}
					}
				}
				return m, nil
			case "enter", "return":
				values := make([]string, len(m.form.fields))
				for i, f := range m.form.fields {
					// For select fields, use option value instead of display label
					if i < len(m.form.fieldDefs) && m.form.fieldDefs[i].Type == FieldTypeSelect {
						options := m.form.fieldDefs[i].Options
						if idx, ok := m.form.selectIndex[i]; ok && idx < len(options) {
							values[i] = options[idx].Value
						}
					} else {
						values[i] = f.Value()
					}
				}
				return m, func() tea.Msg {
					err := m.form.onSubmit(values)
					return crudSubmitMsg{err: err}
				}
			default:
				// For select fields, don't pass input to textinput
				if m.form.focusIndex < len(m.form.fieldDefs) && m.form.fieldDefs[m.form.focusIndex].Type == FieldTypeSelect {
					return m, nil
				}
				// Pass other keys to the focused text input
				var cmd tea.Cmd
				m.form.fields[m.form.focusIndex], cmd = m.form.fields[m.form.focusIndex].Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m *TablePage) updateFormFocus() {
	if m.form == nil {
		return
	}
	for i := 0; i < len(m.form.fields); i++ {
		if i == m.form.focusIndex {
			m.form.fields[i].Focus()
		} else {
			m.form.fields[i].Blur()
		}
	}
}

func (m *TablePage) View() string {
	if m.mode == "edit" && m.form != nil {
		var view string
		title := "Edit Record"
		if m.form.rowIndex < 0 {
			title = "Create Record"
		}
		view += title + "\n\n"
		for i := 0; i < len(m.form.fields); i++ {
			// Show label
			if i < len(m.form.labels) {
				view += m.form.labels[i] + ":\n"
			}
			// Render select fields differently
			if i < len(m.form.fieldDefs) && m.form.fieldDefs[i].Type == FieldTypeSelect {
				// Use → indicator for select fields (different from textinput's cursor)
				focusIndicator := "  "
				if i == m.form.focusIndex {
					focusIndicator = "→ "
				}
				options := m.form.fieldDefs[i].Options
				if len(options) > 0 {
					if idx, ok := m.form.selectIndex[i]; ok {
						// Show current selection with up/down arrows
						view += focusIndicator + "[↑↓] "
						for j := 0; j < len(options); j++ {
							if j == idx {
								view += "◉ " + options[j].Label
							} else {
								view += "○ " + options[j].Label
							}
							if j < len(options)-1 {
								view += " | "
							}
						}
						view += "\n"
					}
				}
			} else {
				// Never show > for text fields, textinput has its own cursor
				view += m.form.fields[i].View() + "\n"
			}
		}
		view += "\n[Tab] Next | [Shift+Tab] Prev | [Enter] Save | [Esc] Cancel"
		if m.err != "" {
			view += "\n\nError: " + m.err
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Render(view)
	}

	view := m.table.View()
	if m.err != "" {
		view = view + "\n\nError: " + m.err
	}
	return view
}

func (m *TablePage) Title() string {
	return m.title
}

func (m *TablePage) IsInForm() bool {
	return m.mode == "edit" && m.form != nil
}
