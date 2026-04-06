package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hackIDLE/fedramp-browser/internal/model"
)

// renderHeader renders the view navigation bar
func (m Model) renderHeader() string {
	views := []struct {
		key   string
		name  string
		state ViewState
	}{
		{"1", "Documents", ViewHome},
		{"2", "Requirements", ViewRequirements},
		{"3", "Definitions", ViewDefinitions},
		{"4", "Indicators", ViewIndicators},
	}

	var tabs []string
	for _, v := range views {
		active := m.view == v.state
		tab := fmt.Sprintf("[%s] %s", v.key, v.name)
		tabs = append(tabs, ViewBadge(tab, active))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...) + "\n\n"
}

// renderDetailContent returns the content for the detail view (used by viewport)
func (m Model) renderDetailContent() string {
	switch item := m.selectedItem.(type) {
	case model.RequirementItem:
		return m.renderRequirementDetail(item)
	case model.DefinitionItem:
		return m.renderDefinitionDetail(item)
	case model.IndicatorItem:
		return m.renderIndicatorDetail(item)
	case model.DocumentItem:
		return m.renderDocumentDetail(item)
	}
	return ""
}

// renderDetailView renders the detail view for any item type
func (m Model) renderDetailView() string {
	var b strings.Builder

	if m.viewportReady {
		b.WriteString(m.viewport.View())
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/↓/j/k scroll • q/ESC back"))

	return AppStyle.Render(b.String())
}

func (m Model) renderRequirementDetail(r model.RequirementItem) string {
	var b strings.Builder

	// Header with badges
	b.WriteString(DetailTitleStyle.Render(r.Name))
	b.WriteString("\n")
	b.WriteString(DocumentBadge(r.DocumentCode))
	b.WriteString(" ")
	if r.PrimaryKeyWord != "" {
		b.WriteString(KeywordBadge(r.PrimaryKeyWord))
		b.WriteString(" ")
	}
	b.WriteString(ImpactBadge(r.Impact.Low, r.Impact.Moderate, r.Impact.High))
	b.WriteString("\n\n")

	// ID
	b.WriteString(DetailLabelStyle.Render("ID:"))
	b.WriteString(DetailValueStyle.Render(r.ID))
	b.WriteString("\n")

	// FKA (formerly known as)
	if r.FKA != "" {
		b.WriteString(DetailLabelStyle.Render("Formerly:"))
		b.WriteString(DimStyle.Render(r.FKA))
		b.WriteString("\n")
	}

	// Category
	if r.Category != "" {
		b.WriteString(DetailLabelStyle.Render("Category:"))
		b.WriteString(DetailValueStyle.Render(r.Category))
		b.WriteString("\n")
	}

	// Applicability
	if r.Applicability != "" {
		b.WriteString(DetailLabelStyle.Render("Applicability:"))
		b.WriteString(DetailValueStyle.Render(r.Applicability))
		b.WriteString("\n")
	}

	// Statement
	b.WriteString(DetailLabelStyle.Render("Statement:"))
	b.WriteString("\n")
	b.WriteString(StatementStyle.Render(wrapText(r.Statement, m.width-10)))
	b.WriteString("\n")

	// Affects
	if len(r.Affects) > 0 {
		b.WriteString(DetailLabelStyle.Render("Affects:"))
		b.WriteString(DetailValueStyle.Render(strings.Join(r.Affects, ", ")))
		b.WriteString("\n")
	}

	// Note
	if r.Note != "" {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Note:"))
		b.WriteString("\n")
		b.WriteString(NoteStyle.Render(wrapText(r.Note, m.width-10)))
	}

	return b.String()
}

func (m Model) renderDefinitionDetail(d model.DefinitionItem) string {
	var b strings.Builder

	// Header
	b.WriteString(DetailTitleStyle.Render(d.Term))
	b.WriteString("\n\n")

	// ID
	b.WriteString(DetailLabelStyle.Render("ID:"))
	b.WriteString(DetailValueStyle.Render(d.ID))
	b.WriteString("\n")

	// FKA (formerly known as)
	if d.FKA != "" {
		b.WriteString(DetailLabelStyle.Render("Formerly:"))
		b.WriteString(DimStyle.Render(d.FKA))
		b.WriteString("\n")
	}

	// Alternatives
	if len(d.Alts) > 0 {
		b.WriteString(DetailLabelStyle.Render("Also known as:"))
		b.WriteString(DetailValueStyle.Render(strings.Join(d.Alts, ", ")))
		b.WriteString("\n")
	}

	// Definition
	b.WriteString("\n")
	b.WriteString(DetailLabelStyle.Render("Definition:"))
	b.WriteString("\n")
	b.WriteString(StatementStyle.Render(wrapText(d.Text, m.width-10)))

	// Note
	if d.Note != "" {
		b.WriteString("\n\n")
		b.WriteString(DetailLabelStyle.Render("Note:"))
		b.WriteString("\n")
		b.WriteString(NoteStyle.Render(wrapText(d.Note, m.width-10)))
	}

	// Reference
	if d.Reference != "" || d.ReferenceURL != "" {
		b.WriteString("\n\n")
		b.WriteString(DetailLabelStyle.Render("Reference:"))
		if d.Reference != "" {
			b.WriteString(DetailValueStyle.Render(d.Reference))
		}
		if d.ReferenceURL != "" {
			b.WriteString("\n")
			b.WriteString(DetailURLStyle.Render(d.ReferenceURL))
		}
	}

	return b.String()
}

func (m Model) renderIndicatorDetail(ind model.IndicatorItem) string {
	var b strings.Builder

	// Header with badges
	b.WriteString(DetailTitleStyle.Render(ind.Name))
	b.WriteString("\n")
	if ind.Retired {
		b.WriteString(RetiredBadge())
		b.WriteString(" ")
	}
	b.WriteString(ImpactBadge(ind.Impact.Low, ind.Impact.Moderate, ind.Impact.High))
	b.WriteString("\n\n")

	// ID
	b.WriteString(DetailLabelStyle.Render("ID:"))
	b.WriteString(DetailValueStyle.Render(ind.ID))
	b.WriteString("\n")

	// FKA (formerly known as)
	if ind.FKA != "" {
		b.WriteString(DetailLabelStyle.Render("Formerly:"))
		b.WriteString(DimStyle.Render(ind.FKA))
		b.WriteString("\n")
	}

	// Theme
	b.WriteString(DetailLabelStyle.Render("Theme:"))
	b.WriteString(DetailValueStyle.Render(ind.ThemeName))
	b.WriteString("\n")

	if ind.ThemeDesc != "" {
		b.WriteString(DimStyle.Render(wrapText(ind.ThemeDesc, m.width-10)))
		b.WriteString("\n")
	}

	// Statement
	b.WriteString("\n")
	b.WriteString(DetailLabelStyle.Render("Statement:"))
	b.WriteString("\n")
	b.WriteString(StatementStyle.Render(wrapText(ind.Statement, m.width-10)))

	// Controls
	if len(ind.Controls) > 0 {
		b.WriteString("\n\n")
		b.WriteString(DetailLabelStyle.Render(fmt.Sprintf("SP 800-53 Controls (%d):", len(ind.Controls))))
		b.WriteString("\n")
		for _, ctrl := range ind.Controls {
			b.WriteString(ControlStyle.Render(fmt.Sprintf("  %s: %s", ctrl.ControlID, ctrl.Title)))
			b.WriteString("\n")
		}
	}

	// Note
	if ind.Note != "" {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Note:"))
		b.WriteString("\n")
		b.WriteString(NoteStyle.Render(wrapText(ind.Note, m.width-10)))
	}

	// Reference
	if ind.Reference != "" || ind.ReferenceURL != "" {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Reference:"))
		if ind.Reference != "" {
			b.WriteString(DetailValueStyle.Render(ind.Reference))
		}
		if ind.ReferenceURL != "" {
			b.WriteString("\n")
			b.WriteString(DetailURLStyle.Render(ind.ReferenceURL))
		}
	}

	return b.String()
}

func (m Model) renderDocumentDetail(d model.DocumentItem) string {
	var b strings.Builder

	b.WriteString(DetailTitleStyle.Render(d.Name))
	b.WriteString("\n")
	b.WriteString(DocumentBadge(d.Code))
	b.WriteString("\n\n")

	// Basic info
	b.WriteString(DetailLabelStyle.Render("Description:"))
	b.WriteString(DetailValueStyle.Render(d.Document.Description))
	b.WriteString("\n")

	if d.RequirementCount > 0 {
		b.WriteString(DetailLabelStyle.Render("Requirements:"))
		b.WriteString(DetailValueStyle.Render(fmt.Sprintf("%d", d.RequirementCount)))
		b.WriteString("\n")
	}

	// Purpose
	if d.Document.Purpose != "" {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Purpose:"))
		b.WriteString("\n")
		b.WriteString(StatementStyle.Render(wrapText(d.Document.Purpose, m.width-10)))
		b.WriteString("\n")
	}

	// Expected Outcomes
	if len(d.Document.ExpectedOutcomes) > 0 {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Expected Outcomes:"))
		b.WriteString("\n")
		for _, outcome := range d.Document.ExpectedOutcomes {
			b.WriteString(DimStyle.Render("  • "))
			b.WriteString(DetailValueStyle.Render(wrapText(outcome, m.width-14)))
			b.WriteString("\n")
		}
	}

	// Program Status (Effective Info)
	if len(d.Document.EffectiveInfo) > 0 {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Program Status:"))
		b.WriteString("\n")
		for version, eff := range d.Document.EffectiveInfo {
			b.WriteString(ControlStyle.Render(fmt.Sprintf("  %s: ", version)))
			b.WriteString(DetailValueStyle.Render(eff.Is))
			if eff.CurrentStatus != "" {
				b.WriteString(DimStyle.Render(fmt.Sprintf(" (%s)", eff.CurrentStatus)))
			}
			b.WriteString("\n")
			if eff.StartDate != "" {
				b.WriteString(DimStyle.Render(fmt.Sprintf("    Start: %s", eff.StartDate)))
				if eff.EndDate != "" {
					b.WriteString(DimStyle.Render(fmt.Sprintf(" - End: %s", eff.EndDate)))
				}
				b.WriteString("\n")
			}
		}
	}

	// Authority
	if len(d.Document.Authority) > 0 {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render("Authority:"))
		b.WriteString("\n")
		for _, auth := range d.Document.Authority {
			b.WriteString(DimStyle.Render("  • "))
			b.WriteString(DetailValueStyle.Render(auth.Reference))
			b.WriteString("\n")
			if auth.Description != "" {
				b.WriteString(DimStyle.Render(fmt.Sprintf("    %s", wrapText(auth.Description, m.width-14))))
				b.WriteString("\n")
			}
			if auth.ReferenceURL != "" {
				b.WriteString(DetailURLStyle.Render(fmt.Sprintf("    %s", auth.ReferenceURL)))
				b.WriteString("\n")
			}
		}
	}

	// Releases
	if len(d.Document.Releases) > 0 {
		b.WriteString("\n")
		b.WriteString(DetailLabelStyle.Render(fmt.Sprintf("Releases (%d):", len(d.Document.Releases))))
		b.WriteString("\n")
		// Show latest 3 releases
		maxReleases := 3
		if len(d.Document.Releases) < maxReleases {
			maxReleases = len(d.Document.Releases)
		}
		for i := 0; i < maxReleases; i++ {
			rel := d.Document.Releases[i]
			b.WriteString(ControlStyle.Render(fmt.Sprintf("  %s", rel.ID)))
			if rel.PublishedDate != "" {
				b.WriteString(DimStyle.Render(fmt.Sprintf(" (%s)", rel.PublishedDate)))
			}
			b.WriteString("\n")
			if rel.Description != "" {
				b.WriteString(DimStyle.Render(fmt.Sprintf("    %s", wrapText(rel.Description, m.width-14))))
				b.WriteString("\n")
			}
		}
		if len(d.Document.Releases) > maxReleases {
			b.WriteString(DimStyle.Render(fmt.Sprintf("  ... and %d more releases", len(d.Document.Releases)-maxReleases)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Press Enter to view requirements for this document"))

	return b.String()
}

// wrapText wraps text to the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)

		if lineLen+wordLen+1 > width && lineLen > 0 {
			result.WriteString("\n")
			lineLen = 0
		}

		if i > 0 && lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}

		result.WriteString(word)
		lineLen += wordLen
	}

	return result.String()
}
