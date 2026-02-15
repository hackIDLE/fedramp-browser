package api

import "encoding/json"

// TopLevelInfo represents the global info section of the consolidated document
type TopLevelInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	LastUpdated string `json:"last_updated"`
}

// ConsolidatedDocument represents the entire FRMR.documentation.json structure
type ConsolidatedDocument struct {
	Info TopLevelInfo               `json:"info"`
	FRD  FRDSection                 `json:"FRD"`
	FRR  map[string]json.RawMessage `json:"FRR"`
	KSI  map[string]ThemeJSON       `json:"KSI"`
}

// FRDSection represents the FedRAMP Definitions section
type FRDSection struct {
	Info DocumentInfo `json:"info"`
	Data struct {
		Both map[string]DefinitionJSON `json:"both"`
	} `json:"data"`
}

// FRRProcess represents a single FRR process (e.g., VDR, UCM, ADS)
type FRRProcess struct {
	Info DocumentInfo                                        `json:"info"`
	Data map[string]map[string]map[string]json.RawMessage    `json:"data"`
}
// Data structure: applicability("both"/"20x"/"rev5") → label("CSO"/"TRC"/...) → id → requirement

// DocumentInfo represents the common info structure in FedRAMP documents
type DocumentInfo struct {
	Name        string                    `json:"name"`
	ShortName   string                    `json:"short_name"`
	WebName     string                    `json:"web_name"`
	Effective   map[string]EffectiveInfo  `json:"effective"`
	Releases    []Release                 `json:"releases"`
	FrontMatter FrontMatter               `json:"front_matter"`
	Labels      map[string]LabelInfo      `json:"labels"`
}

// LabelInfo describes a label used in an FRR process
type LabelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// EffectiveInfo represents version-specific applicability
type EffectiveInfo struct {
	Is            string   `json:"is"`
	SignupURL     string   `json:"signup_url"`
	CurrentStatus string   `json:"current_status"`
	StartDate     string   `json:"start_date"`
	EndDate       string   `json:"end_date"`
	Comments      []string `json:"comments"`
	Warnings      []string `json:"warnings"`
}

// RelatedRFC represents a related RFC reference
type RelatedRFC struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	DiscussionURL string `json:"discussion_url"`
	ShortName     string `json:"short_name"`
	FullName      string `json:"full_name"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
}

// Release represents a document release version
type Release struct {
	ID            string       `json:"id"`
	PublishedDate string       `json:"published_date"`
	Description   string       `json:"description"`
	PublicComment bool         `json:"public_comment"`
	RelatedRFCs   []RelatedRFC `json:"related_rfcs"`
}

// FrontMatter contains authority and purpose information
type FrontMatter struct {
	Authority        []Authority `json:"authority"`
	Purpose          string      `json:"purpose"`
	ExpectedOutcomes []string    `json:"expected_outcomes"`
}

// Authority represents a legal authority reference
type Authority struct {
	Reference     string `json:"reference"`
	ReferenceURL  string `json:"reference_url"`
	Description   string `json:"description"`
	Delegation    string `json:"delegation"`
	DelegationURL string `json:"delegation_url"`
}

// UpdateEntry represents a changelog entry
type UpdateEntry struct {
	Date    string `json:"date"`
	Comment string `json:"comment"`
}

// VariesByLevel holds level-specific requirement/indicator variations
type VariesByLevel struct {
	Low      *LevelVariation `json:"low"`
	Moderate *LevelVariation `json:"moderate"`
	High     *LevelVariation `json:"high"`
}

// LevelVariation represents a level-specific variation of a requirement
type LevelVariation struct {
	Statement     string `json:"statement"`
	PrimaryKeyWord string `json:"primary_key_word"`
	TimeframeType string `json:"timeframe_type"`
	TimeframeNum  int    `json:"timeframe_num"`
}

// RequirementJSON represents a requirement from FRR sections
type RequirementJSON struct {
	ID                   string             `json:"id"`
	FKA                  string             `json:"fka"`
	Statement            string             `json:"statement"`
	Name                 string             `json:"name"`
	Impact               ImpactJSON         `json:"impact"`
	Affects              []string           `json:"affects"`
	PrimaryKeyWord       string             `json:"primary_key_word"`
	Note                 string             `json:"note"`
	Notes                []string           `json:"notes"`
	Terms                []string           `json:"terms"`
	Examples             json.RawMessage    `json:"examples"`
	Danger               string             `json:"danger"`
	Notification         json.RawMessage    `json:"notification"`
	Updated              []UpdateEntry      `json:"updated"`
	VariesByLevel        *VariesByLevel     `json:"varies_by_level"`
	FollowingInformation FollowingInfoField `json:"-"` // Custom unmarshaling
	RawFollowingInfo     json.RawMessage    `json:"following_information"`
}

// FollowingInfoField handles following_information which can be string or []RequirementJSON
type FollowingInfoField []RequirementJSON

// UnmarshalFollowingInfo processes the raw following_information field after initial unmarshal
func (r *RequirementJSON) UnmarshalFollowingInfo() {
	if r.RawFollowingInfo == nil || len(r.RawFollowingInfo) == 0 {
		return
	}
	// Try as array of requirements first
	var reqs []RequirementJSON
	if err := json.Unmarshal(r.RawFollowingInfo, &reqs); err == nil {
		r.FollowingInformation = reqs
		return
	}
	// Otherwise it's a string or other type, ignore
}

// ImpactJSON represents impact levels
type ImpactJSON struct {
	Low      bool `json:"low"`
	Moderate bool `json:"moderate"`
	High     bool `json:"high"`
}

// DefinitionJSON represents a FedRAMP definition
type DefinitionJSON struct {
	ID           string        `json:"id"`
	FKA          string        `json:"fka"`
	Term         string        `json:"term"`
	Alts         []string      `json:"alts"`
	Definition   string        `json:"definition"`
	Note         string        `json:"note"`
	Notes        []string      `json:"notes"`
	Reference    string        `json:"reference"`
	ReferenceURL string        `json:"reference_url"`
	Updated      []UpdateEntry `json:"updated"`
}

// IndicatorJSON represents a KSI indicator
type IndicatorJSON struct {
	ID            string          `json:"id"`
	FKA           string          `json:"fka"`
	Name          string          `json:"name"`
	Statement     string          `json:"statement"`
	Impact        ImpactJSON      `json:"impact"`
	RawControls   json.RawMessage `json:"controls"`
	Reference     string          `json:"reference"`
	ReferenceURL  string          `json:"reference_url"`
	Note          string          `json:"note"`
	Retired       bool            `json:"retired"`
	Terms         []string        `json:"terms"`
	Updated       []UpdateEntry   `json:"updated"`
	VariesByLevel *VariesByLevel  `json:"varies_by_level"`
}

// ParseControls handles controls that can be []ControlJSON or []string
func (ind *IndicatorJSON) ParseControls() []ControlJSON {
	if ind.RawControls == nil {
		return nil
	}
	// Try as array of objects first
	var controls []ControlJSON
	if err := json.Unmarshal(ind.RawControls, &controls); err == nil {
		return controls
	}
	// Try as array of strings (just control IDs)
	var ids []string
	if err := json.Unmarshal(ind.RawControls, &ids); err == nil {
		controls = make([]ControlJSON, len(ids))
		for i, id := range ids {
			controls[i] = ControlJSON{ControlID: id}
		}
		return controls
	}
	return nil
}

// ControlJSON represents an SP 800-53 control reference
type ControlJSON struct {
	ControlID string `json:"control_id"`
	Title     string `json:"title"`
}

// ThemeJSON represents a KSI theme
type ThemeJSON struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Theme      string                     `json:"theme"`
	Indicators map[string]IndicatorJSON   `json:"indicators"`
}

// DocumentMetadata holds basic info about a FedRAMP document
type DocumentMetadata struct {
	Code        string
	Name        string
	Description string
}
