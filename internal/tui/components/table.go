package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// Table is an interactive table component with selection and keyboard navigation
type Table struct {
	table    table.Model
	columns  []table.Column
	title    string
	width    int
	height   int
	focused  bool
}

// NewTable creates a new table with columns and rows
func NewTable(columns []table.Column, rows []table.Row, width, height int) Table {
	theme := styles.GetTheme()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	// Set table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(theme.Primary))

	s.Selected = s.Selected.
		Foreground(lipgloss.Color(theme.Primary)).
		Background(lipgloss.Color(theme.Highlight)).
		Bold(true)

	t.SetStyles(s)

	return Table{
		table:   t,
		columns: columns,
		width:   width,
		height:  height,
		focused: true,
	}
}

// NewTableWithTitle creates a new table with a title
func NewTableWithTitle(title string, columns []table.Column, rows []table.Row, width, height int) Table {
	t := NewTable(columns, rows, width, height)
	t.title = title
	return t
}

// Init initializes the table
func (t Table) Init() tea.Cmd {
	return nil
}

// Update updates the table
func (t Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			t.table.MoveUp(1)
		case "down", "j":
			t.table.MoveDown(1)
		case "pgup":
			t.table.MoveUp(t.height)
		case "pgdown":
			t.table.MoveDown(t.height)
		case "home", "g":
			t.table.GotoTop()
		case "end", "G":
			t.table.GotoBottom()
		default:
			t.table, cmd = t.table.Update(msg)
		}
	default:
		t.table, cmd = t.table.Update(msg)
	}

	return t, cmd
}

// View renders the table
func (t Table) View() string {
	if t.title != "" {
		titleStyle := styles.TitleStyle
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(t.title),
			t.table.View(),
		)
	}
	return t.table.View()
}

// SetRows updates the table rows
func (t *Table) SetRows(rows []table.Row) {
	t.table.SetRows(rows)
}

// SetColumns updates the table columns
func (t *Table) SetColumns(columns []table.Column) {
	t.table.SetColumns(columns)
}

// SelectedRow returns the currently selected row
func (t Table) SelectedRow() table.Row {
	return t.table.SelectedRow()
}

// Cursor returns the current cursor position
func (t Table) Cursor() int {
	return t.table.Cursor()
}

// SetCursor sets the cursor position
func (t *Table) SetCursor(n int) {
	t.table.SetCursor(n)
}

// Focus sets the table focus state
func (t *Table) Focus() {
	t.focused = true
	t.table.Focus()
}

// Blur removes focus from the table
func (t *Table) Blur() {
	t.focused = false
	t.table.Blur()
}

// SetWidth sets the table width
func (t *Table) SetWidth(width int) {
	t.width = width
	t.table.SetWidth(width)
}

// SetHeight sets the table height
func (t *Table) SetHeight(height int) {
	t.height = height
	t.table.SetHeight(height)
}

// SetTitle sets the table title
func (t *Table) SetTitle(title string) {
	t.title = title
}

// GetRows returns all table rows
func (t Table) GetRows() []table.Row {
	return t.table.Rows()
}

// GetColumns returns all table columns
func (t Table) GetColumns() []table.Column {
	return t.columns
}

// Helper function to create a column
func NewColumn(title string, width int) table.Column {
	return table.Column{
		Title: title,
		Width: width,
	}
}

// Helper function to create a row from strings
func NewRow(values ...string) table.Row {
	return table.Row(values)
}
