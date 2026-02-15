package model

import "fmt"

// DocumentItem wraps Document for the list component
type DocumentItem struct {
	Document
}

func (d DocumentItem) Title() string {
	return fmt.Sprintf("[%s] %s", d.Code, d.Name)
}

func (d DocumentItem) Description() string {
	if d.RequirementCount > 0 {
		return fmt.Sprintf("%s (%d requirements)", d.Document.Description, d.RequirementCount)
	}
	return d.Document.Description
}

func (d DocumentItem) FilterValue() string {
	return d.Code + " " + d.Name + " " + d.Document.Description
}

// RequirementItem wraps Requirement for the list component
type RequirementItem struct {
	Requirement
}

func (r RequirementItem) Title() string {
	return fmt.Sprintf("[%s] %s", r.ID, r.Name)
}

func (r RequirementItem) Description() string {
	keyword := r.PrimaryKeyWord
	if keyword == "" {
		keyword = "INFO"
	}
	desc := fmt.Sprintf("%s | %s | %s", r.DocumentCode, keyword, r.Impact.String())
	if r.Category != "" {
		desc += " | " + r.Category
	}
	return desc
}

func (r RequirementItem) FilterValue() string {
	val := r.ID + " " + r.Name + " " + r.Statement + " " + r.DocumentCode
	if r.FKA != "" {
		val += " " + r.FKA
	}
	return val
}

// DefinitionItem wraps Definition for the list component
type DefinitionItem struct {
	Definition
}

func (d DefinitionItem) Title() string {
	return d.Term
}

func (d DefinitionItem) Description() string {
	desc := d.Text
	if len(desc) > 100 {
		desc = desc[:97] + "..."
	}
	return desc
}

func (d DefinitionItem) FilterValue() string {
	result := d.Term + " " + d.Text
	for _, alt := range d.Alts {
		result += " " + alt
	}
	if d.FKA != "" {
		result += " " + d.FKA
	}
	return result
}

// IndicatorItem wraps Indicator for the list component
type IndicatorItem struct {
	Indicator
}

func (i IndicatorItem) Title() string {
	prefix := ""
	if i.Retired {
		prefix = "[RETIRED] "
	}
	return fmt.Sprintf("%s%s", prefix, i.Name)
}

func (i IndicatorItem) Description() string {
	controlCount := ""
	if len(i.Controls) > 0 {
		controlCount = fmt.Sprintf(" | %d controls", len(i.Controls))
	}
	return fmt.Sprintf("%s | %s%s", i.ThemeName, i.Impact.String(), controlCount)
}

func (i IndicatorItem) FilterValue() string {
	val := i.ID + " " + i.Name + " " + i.Statement + " " + i.ThemeName
	if i.FKA != "" {
		val += " " + i.FKA
	}
	return val
}
