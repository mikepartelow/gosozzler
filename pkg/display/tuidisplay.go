package display

import (
	"fmt"
	"io"
	"mp/sozzler/pkg/sozzler"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TuiDisplay struct{}

func (d *TuiDisplay) Error(e string) {
	fmt.Println(e)
}

// -----------------------------------------------------------------------------
// Screens
// -----------------------------------------------------------------------------

type screen int

const (
	screenList screen = iota
	screenDetail
)

// -----------------------------------------------------------------------------
// Styles for list view (global)
// -----------------------------------------------------------------------------

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

// -----------------------------------------------------------------------------
// List item + delegate
// -----------------------------------------------------------------------------

type item struct{ recipe *sozzler.Recipe }

func (i item) FilterValue() string { return i.recipe.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("• %s", i.recipe.Name+" "+i.recipe.FancyRating())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

// -----------------------------------------------------------------------------
// Model
// -----------------------------------------------------------------------------

type model struct {
	list     list.Model
	choice   *sozzler.Recipe
	screen   screen
	width    int
	height   int
	quitting bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		// Detail screen navigation
		if m.screen == screenDetail {
			switch msg.String() {
			case "esc", "backspace", "left":
				m.screen = screenList
				return m, nil
			}
			// Ignore other keys while viewing details
			return m, nil
		}

		// List screen navigation
		switch msg.String() {
		case "enter":
			if it, ok := m.list.SelectedItem().(item); ok {
				m.choice = it.recipe
				m.screen = screenDetail
			}
			return m, nil
		}
	}

	// Only update the list when on the list screen
	if m.screen == screenDetail {
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.screen == screenDetail && m.choice != nil {
		card := renderRecipeCard(m.choice)
		hint := lipgloss.NewStyle().Faint(true).MarginTop(1).PaddingLeft(2).
			Render("esc/backspace to go back · q to quit")
		return lipgloss.JoinVertical(lipgloss.Left, card, hint)
	}

	// Optional hint under the list
	hint := lipgloss.NewStyle().Faint(true).MarginTop(1).PaddingLeft(2).
		Render("enter to open · / to filter · q to quit")
	return "\n" + m.list.View() + "\n" + hint
}

// -----------------------------------------------------------------------------
// Public API
// -----------------------------------------------------------------------------

func (d *TuiDisplay) List(recipes []*sozzler.Recipe) {
	var items []list.Item
	for _, recipe := range recipes {
		items = append(items, item{recipe: recipe})
	}

	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "🍸🍹 Recipes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	km := l.KeyMap

	km.PrevPage.SetKeys(append(km.PrevPage.Keys(), "ctrl+b")...)
	km.NextPage.SetKeys(append(km.NextPage.Keys(), "ctrl+f")...)

	km.PrevPage.SetHelp("PgUp/^B", "page up")
	km.NextPage.SetHelp("PgDn/^F", "page down")

	l.KeyMap = km

	m := model{list: l, screen: screenList}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// Show remains useful for non-interactive output (e.g., piping)
// It reuses the same renderer as the interactive detail screen.
func (d *TuiDisplay) Show(recipe *sozzler.Recipe) {
	fmt.Println(renderRecipeCard(recipe))
}
