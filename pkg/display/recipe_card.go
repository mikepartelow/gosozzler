package display

import (
	"fmt"
	"mp/sozzler/pkg/sozzler"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type recipeCard struct {
	recipe *sozzler.Recipe
	style  lipgloss.Style
}

func (rc *recipeCard) Title() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#0F172A", Dark: "#FDE68A"}).
		Render(rc.recipe.Name + " " + rc.recipe.FancyRating())
}

func (rc *recipeCard) NotesTitle() string {
	return lipgloss.NewStyle().
		Bold(true).
		Faint(true).
		MarginTop(1).
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}).
		Render("Notes")
}

func (rc *recipeCard) Notes() string {
	// No width measurement or wrapping: let content determine the box size.
	notesBody := lipgloss.NewStyle().Faint(true).Render(strings.TrimSpace(rc.recipe.Notes))
	return lipgloss.JoinVertical(lipgloss.Left, rc.NotesTitle(), notesBody)
}

func (rc *recipeCard) Ingredients() string {
	bullet := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#10B981", Dark: "#34D399"}).
		Render("• ")

	qtyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#0369A1", Dark: "#22D3EE"})
	unitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6D28D9", Dark: "#A78BFA"})
	ingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#111827", Dark: "#E5E7EB"})

	// Align columns based on actual content widths (still fine for auto-fit).
	qtyW, unitW := 0, 0
	for _, c := range rc.recipe.Components {
		if x := lipgloss.Width(fmt.Sprint(c.Quantity)); x > qtyW {
			qtyW = x
		}
		if x := lipgloss.Width(strings.TrimSpace(fmt.Sprint(c.Unit))); x > unitW {
			unitW = x
		}
	}
	colQty := qtyStyle.Width(qtyW)
	colUnit := unitStyle.Width(unitW)

	var rows []string
	for _, c := range sozzler.FancyOrder(rc.recipe.Components) {
		q := colQty.Render(fmt.Sprint(c.Quantity))
		u := colUnit.Render(strings.TrimSpace(fmt.Sprint(c.Unit)))
		ing := ingStyle.Render(fmt.Sprint(c.Ingredient))
		row := lipgloss.JoinHorizontal(lipgloss.Left, bullet, q, " ", u, "  ", ing)
		rows = append(rows, row)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderRecipeCard(recipe *sozzler.Recipe) string {
	rc := &recipeCard{
		recipe: recipe,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#C084FC", Dark: "#7C3AED"}).
			Padding(1, 4), // no Width set
	}

	header := lipgloss.JoinHorizontal(lipgloss.Left, "🍸 ", rc.Title())
	ingredientsBlock := rc.Ingredients()
	notesBlock := rc.Notes()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		ingredientsBlock,
		"",
		notesBlock,
	)

	// No Width(...) call: the border will hug the widest line in `content`.
	return rc.style.Render(content)
}

func (d *TuiDisplay) String(s string) {
	fmt.Print(s)
}
