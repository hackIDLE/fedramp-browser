package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hackIDLE/fedramp-browser/internal/model"
)

// ItemDelegate handles rendering of list items
type ItemDelegate struct {
	ShowDescription bool
}

func NewItemDelegate() ItemDelegate {
	return ItemDelegate{ShowDescription: true}
}

func (d ItemDelegate) Height() int {
	if d.ShowDescription {
		return 2
	}
	return 1
}

func (d ItemDelegate) Spacing() int {
	return 0
}

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var title, desc string
	var badges []string

	switch i := item.(type) {
	case model.DocumentItem:
		title = i.Title()
		desc = i.Description()
		badges = append(badges, DocumentBadge(i.Code))

	case model.RequirementItem:
		title = i.Name
		if title == "" {
			title = i.ID
		}
		// Build description line with keyword indicator and impact as colored text
		var descParts []string
		if i.PrimaryKeyWord != "" {
			descParts = append(descParts, KeywordIndicator(i.PrimaryKeyWord))
		}
		descParts = append(descParts, ImpactText(i.Impact.Low, i.Impact.Moderate, i.Impact.High))
		descParts = append(descParts, truncate(i.Statement, 60))
		desc = strings.Join(descParts, " | ")
		badges = append(badges, DocumentBadge(i.DocumentCode))

	case model.DefinitionItem:
		title = i.Term
		desc = truncate(i.Text, 80)

	case model.IndicatorItem:
		title = i.Name
		desc = truncate(i.Statement, 80)
		if i.Retired {
			badges = append(badges, RetiredBadge())
		}
		badges = append(badges, ImpactBadge(i.Impact.Low, i.Impact.Moderate, i.Impact.High))
		if len(i.Controls) > 0 {
			badges = append(badges, ControlStyle.Render(fmt.Sprintf("%d controls", len(i.Controls))))
		}

	default:
		title = item.FilterValue()
	}

	isSelected := index == m.Index()
	isFiltered := m.FilterState() == list.Filtering

	// Style based on selection state
	titleStyle := NormalStyle
	descStyle := DimStyle

	if isSelected && !isFiltered {
		titleStyle = SelectedStyle
		descStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	}

	// Render the item
	var b strings.Builder

	// First line: badges + title
	if len(badges) > 0 {
		b.WriteString(strings.Join(badges, " "))
		b.WriteString(" ")
	}
	b.WriteString(titleStyle.Render(title))
	_, _ = fmt.Fprintln(w, b.String())

	// Second line: description (if enabled)
	if d.ShowDescription && desc != "" {
		_, _ = fmt.Fprintln(w, "  "+descStyle.Render(desc))
	}
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
