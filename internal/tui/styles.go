package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	PrimaryColor   = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	WarningColor   = lipgloss.Color("#FFCC00")
	ErrorColor     = lipgloss.Color("#FF5F56")
	SubtleColor    = lipgloss.Color("#626262")
	WhiteColor     = lipgloss.Color("#FFFFFF")
	BlackColor     = lipgloss.Color("#000000")

	// Document-specific colors
	FRDColor = lipgloss.Color("#FF6B6B") // Definitions - Red
	KSIColor = lipgloss.Color("#4ECDC4") // Indicators - Teal
	VDRColor = lipgloss.Color("#45B7D1") // Vuln Detection - Blue
	UCMColor = lipgloss.Color("#96CEB4") // Crypto Modules - Green
	SCGColor = lipgloss.Color("#FFEAA7") // Secure Config Guide - Yellow
	ADSColor = lipgloss.Color("#DDA0DD") // Auth Data Sharing - Plum
	CCMColor = lipgloss.Color("#98D8C8") // Continuous Monitor - Mint
	FSIColor = lipgloss.Color("#F7DC6F") // Security Inbox - Gold
	ICPColor = lipgloss.Color("#BB8FCE") // Incident Comms - Purple
	MASColor = lipgloss.Color("#85C1E9") // Min Assessment - Sky Blue
	PVAColor = lipgloss.Color("#F8B500") // Persistent Valid - Orange
	SCNColor = lipgloss.Color("#58D68D") // Significant Change - Lime
)

// DocumentColors maps document codes to colors
var DocumentColors = map[string]lipgloss.Color{
	"FRD": FRDColor,
	"KSI": KSIColor,
	"VDR": VDRColor,
	"UCM": UCMColor,
	"SCG": SCGColor,
	"ADS": ADSColor,
	"CCM": CCMColor,
	"FSI": FSIColor,
	"ICP": ICPColor,
	"MAS": MASColor,
	"PVA": PVAColor,
	"SCN": SCNColor,
}

// Styles
var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	DetailTitleStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				MarginBottom(1)

	DetailLabelStyle = lipgloss.NewStyle().
				Foreground(SubtleColor).
				Width(18)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(WhiteColor)

	DetailURLStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Underline(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(SubtleColor).
			MarginTop(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(WhiteColor)

	DimStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	StatementStyle = lipgloss.NewStyle().
			Foreground(WhiteColor).
			MarginTop(1).
			MarginBottom(1)

	NoteStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Italic(true)

	RetiredStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	ControlStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	// Impact text styles (colored text, no background)
	ImpactHighStyle = lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)
	ImpactModStyle  = lipgloss.NewStyle().Foreground(WarningColor).Bold(true)
	ImpactLowStyle  = lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)

	// Keyword subtle indicator styles
	KeywordMustStyle    = lipgloss.NewStyle().Foreground(ErrorColor)
	KeywordShouldStyle  = lipgloss.NewStyle().Foreground(WarningColor)
	KeywordMayStyle     = lipgloss.NewStyle().Foreground(SecondaryColor)
	KeywordMustNotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF2D2D"))
)

// Badge styles
func DocumentBadge(code string) string {
	color, ok := DocumentColors[code]
	if !ok {
		color = SubtleColor
	}
	return lipgloss.NewStyle().
		Foreground(BlackColor).
		Background(color).
		Padding(0, 1).
		Bold(true).
		Render(code)
}

func KeywordBadge(keyword string) string {
	var bg lipgloss.Color
	switch keyword {
	case "MUST":
		bg = ErrorColor
	case "MUST NOT":
		bg = lipgloss.Color("#FF2D2D")
	case "SHOULD":
		bg = WarningColor
	case "SHOULD NOT":
		bg = lipgloss.Color("#CC9900")
	case "MAY":
		bg = SecondaryColor
	default:
		bg = SubtleColor
	}
	return lipgloss.NewStyle().
		Foreground(BlackColor).
		Background(bg).
		Padding(0, 1).
		Bold(true).
		Render(keyword)
}

func ImpactBadge(low, moderate, high bool) string {
	var text string
	var bg lipgloss.Color

	if high {
		text = "HIGH"
		bg = ErrorColor
	} else if moderate {
		text = "MOD"
		bg = WarningColor
	} else if low {
		text = "LOW"
		bg = SecondaryColor
	} else {
		text = "N/A"
		bg = SubtleColor
	}

	return lipgloss.NewStyle().
		Foreground(BlackColor).
		Background(bg).
		Padding(0, 1).
		Bold(true).
		Render(text)
}

func RetiredBadge() string {
	return lipgloss.NewStyle().
		Foreground(WhiteColor).
		Background(ErrorColor).
		Padding(0, 1).
		Bold(true).
		Render("RETIRED")
}

func ViewBadge(name string, active bool) string {
	style := lipgloss.NewStyle().Padding(0, 1)
	if active {
		return style.Foreground(BlackColor).Background(PrimaryColor).Bold(true).Render(name)
	}
	return style.Foreground(SubtleColor).Render(name)
}

// ImpactText returns colored text for impact level (no background)
func ImpactText(low, moderate, high bool) string {
	if high {
		return ImpactHighStyle.Render("HIGH")
	} else if moderate {
		return ImpactModStyle.Render("MOD")
	} else if low {
		return ImpactLowStyle.Render("LOW")
	}
	return DimStyle.Render("N/A")
}

// KeywordIndicator returns a subtle colored dot for keyword display
func KeywordIndicator(keyword string) string {
	switch keyword {
	case "MUST":
		return KeywordMustStyle.Render("●") + " " + KeywordMustStyle.Render("MUST")
	case "MUST NOT":
		return KeywordMustNotStyle.Render("●") + " " + KeywordMustNotStyle.Render("MUST NOT")
	case "SHOULD":
		return KeywordShouldStyle.Render("●") + " " + KeywordShouldStyle.Render("SHOULD")
	case "SHOULD NOT":
		return KeywordShouldStyle.Render("●") + " " + KeywordShouldStyle.Render("SHOULD NOT")
	case "MAY":
		return KeywordMayStyle.Render("●") + " " + KeywordMayStyle.Render("MAY")
	default:
		return ""
	}
}
