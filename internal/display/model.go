package display

import (
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fzf "github.com/koki-develop/go-fzf"
	"github.com/muesli/termenv"
)

type Model struct {
	items   *items
	matches fzf.Matches

	dataChan chan [2]interface{}

	getParams func(obj interface{}, now time.Time) []string
	getKey    func(obj interface{}) string

	// window
	windowWidth     int
	windowHeight    int
	windowYPosition int

	input textinput.Model
}

func New(
	dataChan chan [2]interface{},
	getParams func(obj interface{}, now time.Time) []string,
	getKey func(obj interface{}) string,
) *Model {
	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = lipgloss.NewStyle()
	input.Placeholder = "Filter..."
	input.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	input.TextStyle = lipgloss.NewStyle().Faint(true)
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	return &Model{
		items:     &items{getKey: getKey},
		dataChan:  dataChan,
		getParams: getParams,
		getKey:    getKey,
		input:     input,
	}
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		textinput.Blink,
		m.reload(),
	}

	// Sync Values
	go func() {
		for {
			update, ok := <-m.dataChan
			if !ok {
				time.Sleep(100 * time.Millisecond)
			}

			oldObj := update[0]
			newObj := update[1]

			// Remove Old
			if oldObj != nil {
				for i, o := range m.items.values {
					if m.getKey(o) == m.getKey(oldObj) {
						m.items.values = append(m.items.values[:i], m.items.values[i+1:]...)
						break
					}
				}
			}

			// Add New
			if newObj != nil {
				m.items.values = append([]interface{}{newObj}, m.items.values...)
			}

			m.filter()
		}
	}()

	return tea.Batch(cmds...)
}

func (m *Model) View() string {
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
	if len(m.matches) > 0 {

		itemParams := [][]string{}
		for _, match := range m.matches {
			// log.Printf("%v", match)
			// if match.Index > len(m.Items.values){
			// 	m.filter()
			// }
			itemParams = append(itemParams, m.getParams(m.items.values[match.Index], now))
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

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update tea msgs
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

	// Update filter
	beforeValue := m.input.Value()

	{
		input, cmd := m.input.Update(message)
		m.input = input
		cmds = append(cmds, cmd)
	}

	if beforeValue != m.input.Value() {
		m.filter()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) filter() {
	search := m.input.Value()

	if search == "" {
		var matches fzf.Matches
		for i := range m.items.Len() {
			matches = append(matches, fzf.Match{
				Str:   m.items.ItemString(i),
				Index: i,
			})
		}
		m.matches = matches
		return
	}

	// TODO: Search opts
	m.matches = fzf.Search(m.items, search)
}

func (m *Model) reload() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return reloadMsg{}
	})
}
