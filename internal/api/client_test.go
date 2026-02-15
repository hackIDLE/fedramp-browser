package api

import (
	"encoding/json"
	"testing"
)

func TestParseKSIStructure(t *testing.T) {
	client := NewClient()

	// Fetch the consolidated document
	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	t.Logf("Fetched %d bytes of consolidated data", len(data))

	// Try to parse just the raw structure first
	var rawDoc map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		t.Fatalf("Failed to unmarshal raw doc: %v", err)
	}

	t.Logf("Top-level keys: %v", getKeys(rawDoc))

	// Check if KSI key exists at top level
	ksiRaw, ok := rawDoc["KSI"]
	if !ok {
		t.Log("KSI key not found at top level, checking under FRR")
		// KSI may be under FRR as a process
		frrRaw, frrOk := rawDoc["FRR"]
		if !frrOk {
			t.Fatal("Neither KSI nor FRR found")
		}
		var frr map[string]json.RawMessage
		if err := json.Unmarshal(frrRaw, &frr); err != nil {
			t.Fatalf("Failed to unmarshal FRR: %v", err)
		}
		t.Logf("FRR keys: %v", getKeys(frr))
		return
	}

	t.Logf("KSI raw data length: %d bytes", len(ksiRaw))

	// Parse the KSI section
	var ksiMap map[string]json.RawMessage
	if err := json.Unmarshal(ksiRaw, &ksiMap); err != nil {
		t.Fatalf("Failed to unmarshal KSI section: %v", err)
	}

	t.Logf("KSI theme keys: %v", getKeys(ksiMap))

	// Try parsing one theme
	if afrRaw, ok := ksiMap["AFR"]; ok {
		var theme ThemeJSON
		if err := json.Unmarshal(afrRaw, &theme); err != nil {
			t.Fatalf("Failed to unmarshal AFR theme: %v", err)
		}
		t.Logf("AFR theme: id=%s, name=%s, indicators=%d", theme.ID, theme.Name, len(theme.Indicators))
	}

	// Verify we can parse themes correctly
	var ksiThemes map[string]ThemeJSON
	if err := json.Unmarshal(ksiRaw, &ksiThemes); err != nil {
		t.Fatalf("Failed to unmarshal KSI themes: %v", err)
	}

	t.Logf("Parsed %d themes", len(ksiThemes))
	totalIndicators := 0
	for code, theme := range ksiThemes {
		t.Logf("  Theme %s: %s (%d indicators)", code, theme.Name, len(theme.Indicators))
		totalIndicators += len(theme.Indicators)
	}
	t.Logf("Total indicators: %d", totalIndicators)

	if totalIndicators == 0 {
		t.Error("Expected indicators but got 0")
	}
}

func TestParseIndicators(t *testing.T) {
	client := NewClient()

	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	indicators, err := client.ParseIndicators(data)
	if err != nil {
		t.Fatalf("ParseIndicators failed: %v", err)
	}

	t.Logf("Parsed %d indicators", len(indicators))

	if len(indicators) == 0 {
		t.Error("Expected indicators but got 0")
	}

	// Print first few indicators
	for i, ind := range indicators {
		if i >= 3 {
			break
		}
		t.Logf("Indicator %d: %s - %s (fka: %s)", i, ind.ID, ind.Name, ind.FKA)
	}
}

func getKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestParseRequirements(t *testing.T) {
	client := NewClient()

	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	t.Logf("Fetched %d bytes of consolidated data", len(data))

	// Test parsing ADS requirements (known to exist)
	reqs, err := client.ParseRequirements(data, "ADS")
	if err != nil {
		t.Fatalf("ParseRequirements for ADS failed: %v", err)
	}

	t.Logf("ParseRequirements returned %d ADS requirements", len(reqs))

	if len(reqs) == 0 {
		t.Error("Expected requirements but got 0")
	}

	for i, req := range reqs {
		if i >= 3 {
			break
		}
		t.Logf("Requirement %d: %s - %s (keyword: %s, category: %s, applicability: %s, fka: %s)",
			i, req.ID, req.Name, req.PrimaryKeyWord, req.Category, req.Applicability, req.FKA)
	}
}

func TestParseDefinitions(t *testing.T) {
	client := NewClient()

	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	defs, err := client.ParseDefinitions(data)
	if err != nil {
		t.Fatalf("ParseDefinitions failed: %v", err)
	}

	t.Logf("Parsed %d definitions", len(defs))

	if len(defs) == 0 {
		t.Error("Expected definitions but got 0")
	}

	for i, def := range defs {
		if i >= 3 {
			break
		}
		t.Logf("Definition %d: %s - %s (fka: %s)", i, def.ID, def.Term, def.FKA)
	}
}

func TestParseAllProcessRequirements(t *testing.T) {
	client := NewClient()

	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	processCodes := []string{"VDR", "UCM", "SCG", "ADS", "CCM", "FSI", "ICP", "MAS", "PVA", "SCN", "KSI"}
	totalReqs := 0

	for _, code := range processCodes {
		reqs, err := client.ParseRequirements(data, code)
		if err != nil {
			t.Errorf("ParseRequirements for %s failed: %v", code, err)
			continue
		}
		t.Logf("  %s: %d requirements", code, len(reqs))
		totalReqs += len(reqs)
	}

	t.Logf("Total requirements across all processes: %d", totalReqs)
}

func TestParseDocumentInfo(t *testing.T) {
	client := NewClient()

	data, err := client.FetchConsolidatedDocument()
	if err != nil {
		t.Fatalf("Failed to fetch consolidated document: %v", err)
	}

	// Test FRD info
	frdInfo, err := client.ParseDocumentInfo(data, "FRD")
	if err != nil {
		t.Errorf("ParseDocumentInfo for FRD failed: %v", err)
	} else {
		t.Logf("FRD info: name=%s, short_name=%s", frdInfo.Name, frdInfo.ShortName)
	}

	// Test ADS info (a known FRR process)
	adsInfo, err := client.ParseDocumentInfo(data, "ADS")
	if err != nil {
		t.Errorf("ParseDocumentInfo for ADS failed: %v", err)
	} else {
		t.Logf("ADS info: name=%s, short_name=%s", adsInfo.Name, adsInfo.ShortName)
	}
}
