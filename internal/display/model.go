package display

import (
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Model[T any] struct {
	Items map[string]*T

	ChangeChannel chan [2]*T

	getParams func(obj *T, now time.Time) []string
	getKey    func(obj *T) string

	// window
	windowWidth     int
	windowHeight    int
	windowYPosition int

	input textinput.Model
}

func New[T any](
	initial map[string]*T,
	getParams func(obj *T, now time.Time) []string,
	getKey func(obj *T) string,
) *Model[T] {
	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = lipgloss.NewStyle()
	input.Placeholder = "Filter..."
	input.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	input.TextStyle = lipgloss.NewStyle().Faint(true)
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	return &Model[T]{
		Items:         initial,
		ChangeChannel: make(chan [2]*T),
		getParams:     getParams,
		getKey:        getKey,
		input:         input,
	}
}

func (m *Model[T]) Init() tea.Cmd {
	cmds := []tea.Cmd{
		textinput.Blink,
		m.reload(),
	}

	return tea.Batch(cmds...)
}

func (m *Model[T]) View() string {
	now := time.Now()
	windowStyle := lipgloss.NewStyle()

	rows := []string{}

	// Write Filter Row
	rows = append(rows, m.input.View())

	// Write Separator Row
	var builder strings.Builder
	builder.WriteString(strings.Repeat("─", m.windowWidth))
	rows = append(rows, builder.String())

	// Write Data
	if len(m.Items) > 0 {
		itemParams := [][]string{}
		for _, obj := range m.Items {
			itemParams = append(itemParams, m.getParams(obj, now))
		}

		columnLengths := make([]int, len(itemParams[0]))
		columnBuffer := 3
		for _, line := range itemParams {
			for i, val := range line {
				columnLengths[i] = max(columnLengths[i], len(val))
			}
		}

		lineLength := 0
		for _, l := range columnLengths {
			lineLength += l + columnBuffer
		}

		itemRows := []string{}
		for _, params := range itemParams {
			var buffer strings.Builder
			for i, v := range params {
				buffer.WriteString(v)
				buffer.WriteString(strings.Repeat(" ", columnLengths[i]+columnBuffer-len(v)))
			}
			itemRows = append(itemRows, buffer.String())
		}

		sort.Strings(itemRows)
		rows = append(rows, lipgloss.JoinVertical(lipgloss.Left, itemRows...))
	}

	return windowStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m *Model[T]) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// window
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// m.fixYPosition()
		// m.fixCursor()
		// m.fixWidth()
	case reloadMsg:
		return m, m.reload()
	}
	
	var cmds []tea.Cmd
	// beforeValue := m.input.Value()
	//
	{
		input, cmd := m.input.Update(message)
		m.input= input
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// func (m *Model[T]) filter() {
// 	s := m.input.Value()
//
// }

func (m *Model[T]) reload() tea.Cmd {
	return tea.Tick(1 * time.Second, func(t time.Time) tea.Msg {
		return reloadMsg{}
	})
}
