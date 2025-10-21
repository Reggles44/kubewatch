package display

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/reggles44/kubewatch/pkg/printer"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

type Model struct {
	items   []*unstructured.Unstructured
	matches []int

	dataCh chan watch.Event

	// window
	windowWidth     int
	windowHeight    int
	windowYPosition int

	input textinput.Model
}

func New(dataCh chan watch.Event) *Model {
	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = lipgloss.NewStyle()
	// input.Placeholder = ""
	input.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	input.TextStyle = lipgloss.NewStyle()
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	return &Model{
		dataCh: dataCh,
		input:  input,
	}
}

func (m *Model) addObj(obj *unstructured.Unstructured) {
	if obj == nil {
		return
	}
	m.items = append([]*unstructured.Unstructured{obj}, m.items...)
}

func (m *Model) removeObj(obj *unstructured.Unstructured) {
	if obj == nil {
		return
	}
	for i, o := range m.items {
		if o.GetName() == obj.GetName() {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return
		}
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
			event := <-m.dataCh
			obj := event.Object.(*unstructured.Unstructured)
			obj.GroupVersionKind()
			switch event.Type {
			case watch.Added:
				m.addObj(obj)

			case watch.Modified:
				m.removeObj(obj)
				m.addObj(obj)

			case watch.Deleted:
				m.removeObj(obj)

				// case watch.Bookmark:
				// case watch.Error:
				// default:

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
	builder.WriteString(strconv.Itoa(len(m.matches)))
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(len(m.items)))
	builder.WriteRune(' ')
	borderWidth := max(m.windowWidth-builder.Len(), 0)
	builder.WriteString(strings.Repeat("─", borderWidth))
	rows = append(rows, builder.String())

	// Write Data
	if len(m.matches) > 0 {
		itemParams := [][]string{}
		for _, match := range m.matches {
			obj := m.items[match]
			params := printer.GetParams(obj.GroupVersionKind(), obj, now)
			itemParams = append(itemParams, params)
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

			row := buffer.String()
			if len(row) >= m.windowWidth {
				row = row[:m.windowWidth-3] + "..."
			}
			itemRows = append(itemRows, row)

			if len(itemRows) >= m.windowHeight-len(rows) {
				break
			}
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

	var matches []int
	if search == "" {
		for i := range len(m.items) {
			matches = append(matches, i)
		}
	} else {
		for i := range len(m.items) {
			matches = append(matches, i)
		}

		// for i, obj := range m.items {
		// 	if fuzzy.Match(search, m.resource.Key(obj)) {
		// 		matches = append(matches, i)
		// 	}
		// }
	}

	m.matches = matches
}

func (m *Model) reload() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return reloadMsg{}
	})
}
