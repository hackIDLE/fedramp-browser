package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hackIDLE/fedramp-browser/internal/model"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.view != ViewHome {
		t.Errorf("Expected initial view to be ViewHome, got %v", m.view)
	}

	if !m.loading {
		t.Error("Expected loading to be true initially")
	}

	if m.apiClient == nil {
		t.Error("Expected apiClient to be initialized")
	}
}

func TestViewStateTransitions(t *testing.T) {
	m := NewModel()
	m.loading = false

	// Simulate data loaded
	m.documents = []model.Document{
		{Code: "FRD", Name: "FedRAMP Definitions"},
		{Code: "VDR", Name: "Vulnerability Detection"},
	}
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR", Name: "Test Req"},
	}
	m.definitions = []model.Definition{
		{ID: "FRD-1", Term: "Test Term"},
	}
	m.indicators = []model.Indicator{
		{ID: "KSI-1", Name: "Test Indicator"},
	}
	m.initList()

	tests := []struct {
		name     string
		key      string
		expected ViewState
	}{
		{"Press 2 for Requirements", "2", ViewRequirements},
		{"Press 3 for Definitions", "3", ViewDefinitions},
		{"Press 4 for Indicators", "4", ViewIndicators},
		{"Press 1 for Home", "1", ViewHome},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newM, _ := m.Update(msg)
			updated := newM.(Model)

			if updated.view != tt.expected {
				t.Errorf("Expected view %v, got %v", tt.expected, updated.view)
			}
		})
	}
}

func TestEnterKeyInDetailView(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.width = 100
	m.height = 40
	m.definitions = []model.Definition{
		{ID: "FRD-1", Term: "Test"},
	}
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR"},
	}
	m.indicators = []model.Indicator{}
	m.documents = []model.Document{}
	m.initList()

	// Test with DocumentItem (FRD)
	m.view = ViewDetail
	m.previousView = ViewHome
	m.selectedItem = model.DocumentItem{
		Document: model.Document{Code: "FRD", Name: "Definitions"},
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newM, _ := m.Update(msg)
	updated := newM.(Model)

	if updated.view != ViewDefinitions {
		t.Errorf("Expected view ViewDefinitions for FRD, got %v", updated.view)
	}

	// Test with DocumentItem (VDR)
	updated.view = ViewDetail
	updated.selectedItem = model.DocumentItem{
		Document: model.Document{Code: "VDR", Name: "Vulnerability"},
	}

	newM, _ = updated.Update(msg)
	updated = newM.(Model)

	if updated.view != ViewRequirements {
		t.Errorf("Expected view ViewRequirements for VDR, got %v", updated.view)
	}

	if updated.documentFilter != "VDR" {
		t.Errorf("Expected documentFilter 'VDR', got '%s'", updated.documentFilter)
	}
}

func TestDocumentFilter(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR", Name: "VDR Req 1"},
		{ID: "VDR-2", DocumentCode: "VDR", Name: "VDR Req 2"},
		{ID: "UCM-1", DocumentCode: "UCM", Name: "UCM Req 1"},
	}

	// No filter - should return all
	items := m.getRequirementItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 items without filter, got %d", len(items))
	}

	// With filter
	m.documentFilter = "VDR"
	items = m.getRequirementItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 items with VDR filter, got %d", len(items))
	}

	// Verify all items are VDR
	for _, item := range items {
		req := item.(model.RequirementItem)
		if req.DocumentCode != "VDR" {
			t.Errorf("Expected document code VDR, got %s", req.DocumentCode)
		}
	}
}

func TestClearFilter(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.width = 100
	m.height = 40
	m.view = ViewRequirements
	m.documentFilter = "VDR"
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR"},
	}
	m.definitions = []model.Definition{}
	m.indicators = []model.Indicator{}
	m.documents = []model.Document{}
	m.initList()

	// Press 'f' to clear filter
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
	newM, _ := m.Update(msg)
	updated := newM.(Model)

	if updated.documentFilter != "" {
		t.Errorf("Expected documentFilter to be cleared, got '%s'", updated.documentFilter)
	}
}

func TestKeywordFilter(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR", PrimaryKeyWord: "MUST", Name: "Must Req 1"},
		{ID: "VDR-2", DocumentCode: "VDR", PrimaryKeyWord: "SHOULD", Name: "Should Req 1"},
		{ID: "VDR-3", DocumentCode: "VDR", PrimaryKeyWord: "MUST", Name: "Must Req 2"},
		{ID: "UCM-1", DocumentCode: "UCM", PrimaryKeyWord: "SHOULD", Name: "UCM Should"},
	}

	// No filter - should return all
	items := m.getRequirementItems()
	if len(items) != 4 {
		t.Errorf("Expected 4 items without filter, got %d", len(items))
	}

	// Filter by MUST
	m.keywordFilter = "MUST"
	items = m.getRequirementItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 MUST items, got %d", len(items))
	}

	// Verify all items are MUST
	for _, item := range items {
		req := item.(model.RequirementItem)
		if req.PrimaryKeyWord != "MUST" {
			t.Errorf("Expected keyword MUST, got %s", req.PrimaryKeyWord)
		}
	}

	// Filter by SHOULD
	m.keywordFilter = "SHOULD"
	items = m.getRequirementItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 SHOULD items, got %d", len(items))
	}

	// Combined filter: document + keyword
	m.documentFilter = "VDR"
	m.keywordFilter = "MUST"
	items = m.getRequirementItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 VDR+MUST items, got %d", len(items))
	}
}

func TestKeywordFilterToggle(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.width = 100
	m.height = 40
	m.view = ViewRequirements
	m.requirements = []model.Requirement{
		{ID: "VDR-1", DocumentCode: "VDR", PrimaryKeyWord: "MUST"},
	}
	m.definitions = []model.Definition{}
	m.indicators = []model.Indicator{}
	m.documents = []model.Document{}
	m.initList()

	// Press 'm' to enable MUST filter
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")}
	newM, _ := m.Update(msg)
	updated := newM.(Model)

	if updated.keywordFilter != "MUST" {
		t.Errorf("Expected keywordFilter 'MUST', got '%s'", updated.keywordFilter)
	}

	// Press 'm' again to toggle off
	newM, _ = updated.Update(msg)
	updated = newM.(Model)

	if updated.keywordFilter != "" {
		t.Errorf("Expected keywordFilter to be cleared, got '%s'", updated.keywordFilter)
	}

	// Press 's' to enable SHOULD filter
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	newM, _ = updated.Update(msg)
	updated = newM.(Model)

	if updated.keywordFilter != "SHOULD" {
		t.Errorf("Expected keywordFilter 'SHOULD', got '%s'", updated.keywordFilter)
	}
}

func TestBackNavigation(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.width = 100
	m.height = 40
	m.view = ViewDetail
	m.previousView = ViewRequirements
	m.requirements = []model.Requirement{}
	m.definitions = []model.Definition{}
	m.indicators = []model.Indicator{}
	m.documents = []model.Document{}
	m.initList()

	// Press ESC to go back
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(msg)
	updated := newM.(Model)

	if updated.view != ViewRequirements {
		t.Errorf("Expected view ViewRequirements after ESC, got %v", updated.view)
	}

	// Test with backspace
	updated.view = ViewDetail
	updated.previousView = ViewDefinitions

	msg = tea.KeyMsg{Type: tea.KeyBackspace}
	newM, _ = updated.Update(msg)
	updated = newM.(Model)

	if updated.view != ViewDefinitions {
		t.Errorf("Expected view ViewDefinitions after Backspace, got %v", updated.view)
	}
}
