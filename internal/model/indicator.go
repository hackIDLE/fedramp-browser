package model

// Control represents an SP 800-53 control reference
type Control struct {
	ControlID string
	Title     string
}

// Indicator represents a Key Security Indicator
type Indicator struct {
	ID           string
	FKA          string // Formerly known as (previous ID)
	ThemeCode    string
	ThemeName    string
	ThemeDesc    string
	Name         string
	Statement    string
	Impact       Impact
	Controls     []Control
	Reference    string
	ReferenceURL string
	Note         string
	Retired      bool
}

// HasControls returns true if the indicator has control mappings
func (i Indicator) HasControls() bool {
	return len(i.Controls) > 0
}
