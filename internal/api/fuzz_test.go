package api

import "testing"

func FuzzParseDefinitions(f *testing.F) {
	// Seed with minimal valid structure (new consolidated format)
	f.Add([]byte(`{"FRD":{"data":{"both":{}}}}`))
	f.Add([]byte(`{"FRD":{"data":{"both":{"FRD-TEST":{"term":"Test Term","definition":"A test"}}}}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseDefinitions(data)
	})
}

func FuzzParseIndicators(f *testing.F) {
	// Seed with minimal valid structure
	f.Add([]byte(`{"KSI":{}}`))
	f.Add([]byte(`{"KSI":{"AFR":{"id":"AFR","name":"Test","indicators":[]}}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseIndicators(data)
	})
}

func FuzzParseRequirements(f *testing.F) {
	// Seed with minimal valid structure (new consolidated format)
	f.Add([]byte(`{"FRR":{"ADS":{"data":{}}}}`), "ADS")
	f.Add([]byte(`{"FRR":{"VDR":{"data":{"both":{"CSO":{"VDR-CSO-001":{"statement":"test","primary_key_word":"MUST"}}}}}}}`), "VDR")
	f.Add([]byte(`{}`), "ADS")
	f.Add([]byte(`invalid json`), "")

	f.Fuzz(func(t *testing.T, data []byte, docCode string) {
		client := NewClient()
		// Should not panic on any input
		_, _ = client.ParseRequirements(data, docCode)
	})
}
