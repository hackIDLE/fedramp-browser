package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/ethanolivertroy/fedramp-tui/internal/cache"
	"github.com/ethanolivertroy/fedramp-tui/internal/model"
)

const BaseURL = "https://raw.githubusercontent.com/FedRAMP/docs/main"

// ConsolidatedFilename is the single consolidated FedRAMP document
const ConsolidatedFilename = "FRMR.documentation.json"

// DocumentFiles maps document codes to their metadata
var DocumentFiles = map[string]DocumentMetadata{
	"FRD": {Code: "FRD", Name: "FedRAMP Definitions", Description: "Terms and definitions"},
	"KSI": {Code: "KSI", Name: "Key Security Indicators", Description: "Security indicators with control mappings"},
	"VDR": {Code: "VDR", Name: "Vulnerability Detection & Response", Description: "Vulnerability management requirements"},
	"UCM": {Code: "UCM", Name: "Using Cryptographic Modules", Description: "Cryptographic module requirements"},
	"SCG": {Code: "SCG", Name: "Secure Configuration Guide", Description: "Secure configuration requirements"},
	"ADS": {Code: "ADS", Name: "Authorization Data Sharing", Description: "Data sharing requirements"},
	"CCM": {Code: "CCM", Name: "Collaborative Continuous Monitoring", Description: "Continuous monitoring requirements"},
	"FSI": {Code: "FSI", Name: "FedRAMP Security Inbox", Description: "Security inbox procedures"},
	"ICP": {Code: "ICP", Name: "Incident Communications Procedures", Description: "Incident communication requirements"},
	"MAS": {Code: "MAS", Name: "Minimum Assessment Scope", Description: "Assessment scope requirements"},
	"PVA": {Code: "PVA", Name: "Persistent Validation & Assessment", Description: "Validation and assessment requirements"},
	"SCN": {Code: "SCN", Name: "Significant Change Notifications", Description: "Change notification requirements"},
}

// DocumentOrder defines the display order of documents
var DocumentOrder = []string{"FRD", "KSI", "VDR", "UCM", "SCG", "ADS", "CCM", "FSI", "ICP", "MAS", "PVA", "SCN"}

// Client is an HTTP client for fetching FedRAMP documents
type Client struct {
	httpClient *http.Client
	baseURL    string
	cache      *cache.Cache
	refresh    bool
}

// ClientOption configures the client
type ClientOption func(*Client)

// WithRefresh forces fresh fetch, ignoring cache
func WithRefresh(refresh bool) ClientOption {
	return func(c *Client) {
		c.refresh = refresh
	}
}

// NewClient creates a new API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		baseURL:    BaseURL,
	}

	// Initialize cache (ignore errors, will just fetch fresh)
	if cache, err := cache.New(); err == nil {
		c.cache = cache
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// FetchConsolidatedDocument fetches the single consolidated FRMR.documentation.json
func (c *Client) FetchConsolidatedDocument() ([]byte, error) {
	return c.fetchDocument(ConsolidatedFilename)
}

func (c *Client) fetchDocument(filename string) ([]byte, error) {
	url := c.baseURL + "/" + filename

	// Check cache first (unless refresh is forced)
	if c.cache != nil && !c.refresh {
		if data, ok := c.cache.Get(url); ok {
			return data, nil
		}
	}

	// Fetch from network
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Save to cache
	if c.cache != nil {
		_ = c.cache.Set(url, data)
	}

	return data, nil
}

// ParseConsolidatedDocument parses the top-level structure of the consolidated file
func (c *Client) ParseConsolidatedDocument(data []byte) (*ConsolidatedDocument, error) {
	var doc ConsolidatedDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing consolidated document: %w", err)
	}
	return &doc, nil
}

// ParseDefinitions parses the FRD section into Definition models
func (c *Client) ParseDefinitions(data []byte) ([]model.Definition, error) {
	var doc struct {
		FRD FRDSection `json:"FRD"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	var definitions []model.Definition
	for id, d := range doc.FRD.Data.Both {
		note := d.Note
		if len(d.Notes) > 0 && note == "" {
			for j, n := range d.Notes {
				if j > 0 {
					note += " "
				}
				note += n
			}
		}

		defID := d.ID
		if defID == "" {
			defID = id
		}

		definitions = append(definitions, model.Definition{
			ID:           defID,
			FKA:          d.FKA,
			Term:         d.Term,
			Alts:         d.Alts,
			Text:         d.Definition,
			Note:         note,
			Reference:    d.Reference,
			ReferenceURL: d.ReferenceURL,
		})
	}

	// Sort by ID for stable ordering
	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].ID < definitions[j].ID
	})

	return definitions, nil
}

// ParseIndicators parses the KSI section into Indicator models
func (c *Client) ParseIndicators(data []byte) ([]model.Indicator, error) {
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, fmt.Errorf("parsing raw document: %w", err)
	}

	// Extract the KSI section directly
	ksiRaw, ok := rawDoc["KSI"]
	if !ok {
		return nil, fmt.Errorf("KSI section not found in document")
	}

	var ksiThemes map[string]ThemeJSON
	if err := json.Unmarshal(ksiRaw, &ksiThemes); err != nil {
		return nil, fmt.Errorf("parsing KSI themes: %w", err)
	}

	var indicators []model.Indicator

	for themeCode, theme := range ksiThemes {
		for indID, ind := range theme.Indicators {
			parsedControls := ind.ParseControls()
			controls := make([]model.Control, len(parsedControls))
			for j, ctrl := range parsedControls {
				controls[j] = model.Control{
					ControlID: ctrl.ControlID,
					Title:     ctrl.Title,
				}
			}

			id := ind.ID
			if id == "" {
				id = indID
			}

			indicators = append(indicators, model.Indicator{
				ID:           id,
				FKA:          ind.FKA,
				ThemeCode:    themeCode,
				ThemeName:    theme.Name,
				ThemeDesc:    theme.Theme,
				Name:         ind.Name,
				Statement:    ind.Statement,
				Impact: model.Impact{
					Low:      ind.Impact.Low,
					Moderate: ind.Impact.Moderate,
					High:     ind.Impact.High,
				},
				Controls:     controls,
				Reference:    ind.Reference,
				ReferenceURL: ind.ReferenceURL,
				Note:         ind.Note,
				Retired:      ind.Retired,
			})
		}
	}

	// Sort by ID for stable ordering
	sort.Slice(indicators, func(i, j int) bool {
		return indicators[i].ID < indicators[j].ID
	})

	return indicators, nil
}

// ParseRequirements parses requirements for a given process from the consolidated document
func (c *Client) ParseRequirements(data []byte, docCode string) ([]model.Requirement, error) {
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	frrData, ok := rawDoc["FRR"]
	if !ok {
		return nil, nil
	}

	var frr map[string]json.RawMessage
	if err := json.Unmarshal(frrData, &frr); err != nil {
		return nil, err
	}

	processData, ok := frr[docCode]
	if !ok {
		return nil, nil
	}

	var process FRRProcess
	if err := json.Unmarshal(processData, &process); err != nil {
		return nil, fmt.Errorf("parsing FRR process %s: %w", docCode, err)
	}

	var requirements []model.Requirement

	// Iterate: applicability → label → id → requirement
	for applicability, labels := range process.Data {
		for label, entries := range labels {
			for id, rawReq := range entries {
				var req RequirementJSON
				if err := json.Unmarshal(rawReq, &req); err != nil {
					continue
				}
				req.UnmarshalFollowingInfo()

				reqID := req.ID
				if reqID == "" {
					reqID = id
				}

				r := model.Requirement{
					ID:             reqID,
					FKA:            req.FKA,
					DocumentCode:   docCode,
					Category:       label,
					Applicability:  applicability,
					Statement:      req.Statement,
					Name:           req.Name,
					Impact: model.Impact{
						Low:      req.Impact.Low,
						Moderate: req.Impact.Moderate,
						High:     req.Impact.High,
					},
					Affects:        req.Affects,
					PrimaryKeyWord: req.PrimaryKeyWord,
					Note:           req.Note,
				}

				// Handle varies_by_level
				if req.VariesByLevel != nil {
					if req.PrimaryKeyWord == "" || req.PrimaryKeyWord == "varies_by_level" {
						// Use highest applicable level's keyword
						if req.VariesByLevel.High != nil && req.VariesByLevel.High.PrimaryKeyWord != "" {
							r.PrimaryKeyWord = req.VariesByLevel.High.PrimaryKeyWord
						} else if req.VariesByLevel.Moderate != nil && req.VariesByLevel.Moderate.PrimaryKeyWord != "" {
							r.PrimaryKeyWord = req.VariesByLevel.Moderate.PrimaryKeyWord
						} else if req.VariesByLevel.Low != nil && req.VariesByLevel.Low.PrimaryKeyWord != "" {
							r.PrimaryKeyWord = req.VariesByLevel.Low.PrimaryKeyWord
						}
					}
				}

				// Combine notes if needed
				if r.Note == "" && len(req.Notes) > 0 {
					for j, n := range req.Notes {
						if j > 0 {
							r.Note += " "
						}
						r.Note += n
					}
				}

				requirements = append(requirements, r)

				// Also extract nested following_information requirements
				if len(req.FollowingInformation) > 0 {
					requirements = append(requirements, c.extractFollowingRequirements(req.FollowingInformation, docCode, label, applicability)...)
				}
			}
		}
	}

	// Sort by ID for stable ordering
	sort.Slice(requirements, func(i, j int) bool {
		return requirements[i].ID < requirements[j].ID
	})

	return requirements, nil
}

func (c *Client) extractFollowingRequirements(reqs []RequirementJSON, docCode, category, applicability string) []model.Requirement {
	var requirements []model.Requirement

	for i := range reqs {
		r := &reqs[i]
		r.UnmarshalFollowingInfo()

		req := model.Requirement{
			ID:             r.ID,
			FKA:            r.FKA,
			DocumentCode:   docCode,
			Category:       category,
			Applicability:  applicability,
			Statement:      r.Statement,
			Name:           r.Name,
			Impact: model.Impact{
				Low:      r.Impact.Low,
				Moderate: r.Impact.Moderate,
				High:     r.Impact.High,
			},
			Affects:        r.Affects,
			PrimaryKeyWord: r.PrimaryKeyWord,
			Note:           r.Note,
		}
		requirements = append(requirements, req)

		if len(r.FollowingInformation) > 0 {
			requirements = append(requirements, c.extractFollowingRequirements(r.FollowingInformation, docCode, category, applicability)...)
		}
	}

	return requirements
}

// GetDocumentMetadata returns metadata for all documents
func GetDocumentMetadata() []model.Document {
	docs := make([]model.Document, 0, len(DocumentOrder))
	for _, code := range DocumentOrder {
		meta := DocumentFiles[code]
		docs = append(docs, model.Document{
			Code:        meta.Code,
			Name:        meta.Name,
			Description: meta.Description,
		})
	}
	return docs
}

// ParseDocumentInfo extracts the info section for a specific document
func (c *Client) ParseDocumentInfo(data []byte, docCode string) (*DocumentInfo, error) {
	switch docCode {
	case "FRD":
		var doc struct {
			FRD struct {
				Info DocumentInfo `json:"info"`
			} `json:"FRD"`
		}
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
		return &doc.FRD.Info, nil
	default:
		// FRR processes
		var rawDoc map[string]json.RawMessage
		if err := json.Unmarshal(data, &rawDoc); err != nil {
			return nil, err
		}
		frrData, ok := rawDoc["FRR"]
		if !ok {
			return nil, fmt.Errorf("FRR section not found")
		}
		var frr map[string]json.RawMessage
		if err := json.Unmarshal(frrData, &frr); err != nil {
			return nil, err
		}
		processData, ok := frr[docCode]
		if !ok {
			return nil, fmt.Errorf("process %s not found in FRR", docCode)
		}
		var process struct {
			Info DocumentInfo `json:"info"`
		}
		if err := json.Unmarshal(processData, &process); err != nil {
			return nil, err
		}
		return &process.Info, nil
	}
}

// EnrichDocument populates a Document with info from the JSON
func EnrichDocument(doc *model.Document, info *DocumentInfo) {
	if info == nil {
		return
	}

	// Purpose and expected outcomes
	doc.Purpose = info.FrontMatter.Purpose
	doc.ExpectedOutcomes = info.FrontMatter.ExpectedOutcomes

	// Authority references
	for _, auth := range info.FrontMatter.Authority {
		doc.Authority = append(doc.Authority, model.Authority{
			Reference:    auth.Reference,
			ReferenceURL: auth.ReferenceURL,
			Description:  auth.Description,
		})
	}

	// Releases
	for _, rel := range info.Releases {
		doc.Releases = append(doc.Releases, model.Release{
			ID:            rel.ID,
			PublishedDate: rel.PublishedDate,
			Description:   rel.Description,
		})
	}

	// Effective info (program status)
	if len(info.Effective) > 0 {
		doc.EffectiveInfo = make(map[string]model.EffectiveStatus)
		for version, eff := range info.Effective {
			doc.EffectiveInfo[version] = model.EffectiveStatus{
				Is:            eff.Is,
				CurrentStatus: eff.CurrentStatus,
				StartDate:     eff.StartDate,
				EndDate:       eff.EndDate,
				SignupURL:     eff.SignupURL,
				Comments:      eff.Comments,
			}
		}
	}
}
