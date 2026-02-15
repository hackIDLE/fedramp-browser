package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanolivertroy/fedramp-tui/internal/api"
	"github.com/ethanolivertroy/fedramp-tui/internal/model"
)

// ViewState represents the current view
type ViewState int

const (
	ViewHome ViewState = iota
	ViewRequirements
	ViewDefinitions
	ViewIndicators
	ViewDetail
)

// AffectsOptions defines the cycle order for affects filtering
var AffectsOptions = []string{"", "Providers", "Agencies", "Assessors", "FedRAMP"}

// KeywordOptions defines the cycle order for keyword filtering
var KeywordOptions = []string{"", "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", "MAY"}

// Messages
type DataLoadedMsg struct {
	Documents    []model.Document
	Requirements []model.Requirement
	Definitions  []model.Definition
	Indicators   []model.Indicator
}

type ErrorMsg struct {
	Err error
}

// Model is the main application model
type Model struct {
	// State
	view         ViewState
	previousView ViewState
	loading      bool
	err          error
	width        int
	height       int

	// Data
	documents    []model.Document
	requirements []model.Requirement
	definitions  []model.Definition
	indicators   []model.Indicator

	// Filters
	documentFilter string // Filter requirements by document code
	keywordFilter  string // Filter requirements by keyword (MUST, MUST NOT, SHOULD, SHOULD NOT, MAY)
	affectsFilter  string // Filter requirements by affected party (Providers, Agencies, Assessors, FedRAMP)

	// Selected item for detail view
	selectedItem list.Item

	// Components
	list          list.Model
	spinner       spinner.Model
	viewport      viewport.Model
	viewportReady bool
	apiClient     *api.Client
	keys          KeyMap
}

// ModelOption configures the TUI model
type ModelOption func(*Model)

// WithRefresh forces fresh data fetch
func WithRefresh(refresh bool) ModelOption {
	return func(m *Model) {
		m.apiClient = api.NewClient(api.WithRefresh(refresh))
	}
}

// NewModel creates a new application model
func NewModel(opts ...ModelOption) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(PrimaryColor)

	m := Model{
		spinner:   s,
		loading:   true,
		view:      ViewHome,
		apiClient: api.NewClient(),
		keys:      DefaultKeyMap(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchData(),
	)
}

func (m Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		data, err := m.apiClient.FetchConsolidatedDocument()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		var allRequirements []model.Requirement
		var definitions []model.Definition
		var indicators []model.Indicator

		// Parse definitions from FRD section
		if defs, err := m.apiClient.ParseDefinitions(data); err == nil {
			definitions = defs
		}

		// Parse indicators from KSI section
		if inds, err := m.apiClient.ParseIndicators(data); err == nil {
			indicators = inds
		}

		// Parse requirements for each FRR process document
		for _, code := range api.DocumentOrder {
			if code == "FRD" || code == "KSI" {
				continue
			}
			if reqs, err := m.apiClient.ParseRequirements(data, code); err == nil {
				allRequirements = append(allRequirements, reqs...)
			}
		}

		// Also parse KSI process requirements (KSI has both themes and process requirements)
		if ksiReqs, err := m.apiClient.ParseRequirements(data, "KSI"); err == nil {
			allRequirements = append(allRequirements, ksiReqs...)
		}

		// Build document list with counts and rich info
		documents := api.GetDocumentMetadata()
		reqCounts := make(map[string]int)
		for _, r := range allRequirements {
			reqCounts[r.DocumentCode]++
		}
		for i := range documents {
			documents[i].RequirementCount = reqCounts[documents[i].Code]
			// Enrich with info from JSON
			if info, err := m.apiClient.ParseDocumentInfo(data, documents[i].Code); err == nil {
				api.EnrichDocument(&documents[i], info)
			}
		}
		// Special counts for FRD and KSI
		for i := range documents {
			if documents[i].Code == "FRD" {
				documents[i].RequirementCount = len(definitions)
			} else if documents[i].Code == "KSI" {
				documents[i].RequirementCount = len(indicators)
			}
		}

		return DataLoadedMsg{
			Documents:    documents,
			Requirements: allRequirements,
			Definitions:  definitions,
			Indicators:   indicators,
		}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keys in detail view
		if m.view == ViewDetail {
			// ESC/Backspace/q to go back
			if msg.Type == tea.KeyEsc || msg.Type == tea.KeyBackspace || msg.Type == tea.KeyEscape || msg.String() == "q" {
				m.view = m.previousView
				m.viewportReady = false
				m.updateListForView()
				return m, nil
			}
			// Enter to navigate from document detail to its requirements
			if msg.Type == tea.KeyEnter {
				if doc, ok := m.selectedItem.(model.DocumentItem); ok {
					switch doc.Code {
					case "FRD":
						m.view = ViewDefinitions
						m.documentFilter = ""
					case "KSI":
						m.view = ViewIndicators
						m.documentFilter = ""
					default:
						m.view = ViewRequirements
						m.documentFilter = doc.Code
					}
					m.viewportReady = false
					m.updateListForView()
					return m, nil
				}
			}
			// Pass other keys to viewport for scrolling
			if m.viewportReady {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}

		// Don't handle keys while filtering
		if m.list.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "enter":
			// Handle Enter in list view to go to detail
			if !m.loading && m.list.FilterState() != list.Filtering {
				if item := m.list.SelectedItem(); item != nil {
					m.selectedItem = item
					m.previousView = m.view
					m.view = ViewDetail
					// Initialize viewport for scrolling
					m.viewport = viewport.New(m.width-4, m.height-8)
					m.viewport.SetContent(m.renderDetailContent())
					m.viewportReady = true
					return m, nil
				}
			}
		case "f":
			// Clear all filters when in requirements view
			if m.view == ViewRequirements && (m.documentFilter != "" || m.keywordFilter != "" || m.affectsFilter != "") {
				m.documentFilter = ""
				m.keywordFilter = ""
				m.affectsFilter = ""
				m.updateListForView()
				return m, nil
			}
		case "x":
			// Cycle through affects filter in requirements view
			if m.view == ViewRequirements {
				currentIdx := 0
				for i, opt := range AffectsOptions {
					if opt == m.affectsFilter {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(AffectsOptions)
				m.affectsFilter = AffectsOptions[nextIdx]
				m.updateListForView()
				return m, nil
			}
		case "k":
			// Cycle through keyword filter in requirements view
			if m.view == ViewRequirements {
				currentIdx := 0
				for i, opt := range KeywordOptions {
					if opt == m.keywordFilter {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(KeywordOptions)
				m.keywordFilter = KeywordOptions[nextIdx]
				m.updateListForView()
				return m, nil
			}
		case "m":
			// Toggle MUST filter in requirements view
			if m.view == ViewRequirements {
				if m.keywordFilter == "MUST" {
					m.keywordFilter = ""
				} else {
					m.keywordFilter = "MUST"
				}
				m.updateListForView()
				return m, nil
			}
		case "s":
			// Toggle SHOULD filter in requirements view (but not when 's' is for other things)
			if m.view == ViewRequirements && m.list.FilterState() != list.Filtering {
				if m.keywordFilter == "SHOULD" {
					m.keywordFilter = ""
				} else {
					m.keywordFilter = "SHOULD"
				}
				m.updateListForView()
				return m, nil
			}
		case "1":
			if m.view != ViewDetail {
				m.view = ViewHome
				m.documentFilter = ""
				m.keywordFilter = ""
				m.updateListForView()
			}
		case "2":
			if m.view != ViewDetail {
				m.view = ViewRequirements
				m.updateListForView()
			}
		case "3":
			if m.view != ViewDetail {
				m.view = ViewDefinitions
				m.updateListForView()
			}
		case "4":
			if m.view != ViewDetail {
				m.view = ViewIndicators
				m.updateListForView()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.loading {
			m.list.SetSize(msg.Width-4, msg.Height-10)
		}
		if m.viewportReady {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 8
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case DataLoadedMsg:
		m.loading = false
		m.documents = msg.Documents
		m.requirements = msg.Requirements
		m.definitions = msg.Definitions
		m.indicators = msg.Indicators

		// Initialize list with documents
		m.initList()
		return m, nil

	case ErrorMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil
	}

	// Pass messages to list component
	if !m.loading && m.err == nil && m.view != ViewDetail {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) initList() {
	delegate := NewItemDelegate()
	m.list = list.New(m.getDocumentItems(), delegate, m.width-4, m.height-10)
	m.list.Title = "FedRAMP Documentation"
	m.list.SetShowStatusBar(true)
	m.list.SetFilteringEnabled(true)
	m.list.Styles.Title = TitleStyle
	m.list.FilterInput.Prompt = "Filter: "
	m.list.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(PrimaryColor)

	// Use exact substring matching
	m.list.Filter = func(term string, targets []string) []list.Rank {
		var ranks []list.Rank
		term = strings.ToLower(term)
		for i, target := range targets {
			if strings.Contains(strings.ToLower(target), term) {
				ranks = append(ranks, list.Rank{Index: i})
			}
		}
		return ranks
	}
}

func (m *Model) updateListForView() {
	var items []list.Item
	var title string

	switch m.view {
	case ViewHome:
		items = m.getDocumentItems()
		title = "FedRAMP Documents"
	case ViewRequirements:
		items = m.getRequirementItems()
		var filterHints []string
		if m.documentFilter != "" {
			filterHints = append(filterHints, m.documentFilter)
		}
		if m.keywordFilter != "" {
			filterHints = append(filterHints, m.keywordFilter+" only")
		}
		if m.affectsFilter != "" {
			filterHints = append(filterHints, m.affectsFilter)
		}
		if len(filterHints) > 0 {
			title = fmt.Sprintf("Requirements (%d) [%s] - x: affects, k: keyword, f: clear", len(items), strings.Join(filterHints, ", "))
		} else {
			title = fmt.Sprintf("FedRAMP Requirements (%d) - x: affects, m: MUST, s: SHOULD, k: cycle", len(items))
		}
	case ViewDefinitions:
		items = m.getDefinitionItems()
		title = fmt.Sprintf("FedRAMP Definitions (%d)", len(m.definitions))
	case ViewIndicators:
		items = m.getIndicatorItems()
		title = fmt.Sprintf("Key Security Indicators (%d)", len(m.indicators))
	}

	m.list.SetItems(items)
	m.list.Title = title
	m.list.ResetSelected()
	m.list.ResetFilter()
}

func (m Model) getDocumentItems() []list.Item {
	items := make([]list.Item, len(m.documents))
	for i, d := range m.documents {
		items[i] = model.DocumentItem{Document: d}
	}
	return items
}

func (m Model) getRequirementItems() []list.Item {
	var items []list.Item
	for _, r := range m.requirements {
		// Filter by document if set
		if m.documentFilter != "" && r.DocumentCode != m.documentFilter {
			continue
		}
		// Filter by keyword if set
		if m.keywordFilter != "" && r.PrimaryKeyWord != m.keywordFilter {
			continue
		}
		// Filter by affects if set
		if m.affectsFilter != "" {
			found := false
			for _, a := range r.Affects {
				if a == m.affectsFilter {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		items = append(items, model.RequirementItem{Requirement: r})
	}
	return items
}

func (m Model) getDefinitionItems() []list.Item {
	items := make([]list.Item, len(m.definitions))
	for i, d := range m.definitions {
		items[i] = model.DefinitionItem{Definition: d}
	}
	return items
}

func (m Model) getIndicatorItems() []list.Item {
	items := make([]list.Item, len(m.indicators))
	for i, ind := range m.indicators {
		items[i] = model.IndicatorItem{Indicator: ind}
	}
	return items
}

// View renders the model
func (m Model) View() string {
	if m.loading {
		return AppStyle.Render(
			fmt.Sprintf("\n\n   %s Loading FedRAMP documents...\n\n", m.spinner.View()),
		)
	}

	if m.err != nil {
		return AppStyle.Render(
			lipgloss.NewStyle().Foreground(ErrorColor).Render(
				fmt.Sprintf("\n\n   Error: %v\n\n   Press q to quit.", m.err),
			),
		)
	}

	switch m.view {
	case ViewDetail:
		return m.renderDetailView()
	default:
		// Constrain list height to leave room for header
		headerHeight := 4
		listHeight := m.height - headerHeight - 4
		if listHeight < 10 {
			listHeight = 10 // minimum height
		}
		listStyle := lipgloss.NewStyle().MaxHeight(listHeight)
		return AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			m.renderHeader(),
			listStyle.Render(m.list.View())))
	}
}
