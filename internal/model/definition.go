package model

// Definition represents a FedRAMP term definition
type Definition struct {
	ID           string
	FKA          string // Formerly known as (previous ID)
	Term         string
	Alts         []string
	Text         string // The definition text
	Note         string
	Reference    string
	ReferenceURL string
}

// HasAlternatives returns true if there are alternative terms
func (d Definition) HasAlternatives() bool {
	return len(d.Alts) > 0
}

// HasReference returns true if there's a reference
func (d Definition) HasReference() bool {
	return d.Reference != "" || d.ReferenceURL != ""
}
