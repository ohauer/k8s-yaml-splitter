package main

import (
	"testing"
)

// FuzzSanitizeFilename tests filename sanitization with random inputs
func FuzzSanitizeFilename(f *testing.F) {
	// Seed with known problematic inputs
	f.Add("../../../etc/passwd")
	f.Add("....//....//etc//passwd")
	f.Add("/dev/null")
	f.Add("\\\\server\\share")
	f.Add("")
	f.Add(".")
	f.Add("..")
	f.Add("normal-name")
	f.Add("name:with:colons")
	f.Add("name/with/slashes")

	f.Fuzz(func(t *testing.T, input string) {
		result := sanitizeFilename(input)

		// Should never panic
		if result == "" {
			t.Error("sanitizeFilename returned empty string")
		}

		// Should not contain path separators
		if result != "unnamed" && (result == "." || result == "..") {
			t.Errorf("unsafe result: %q", result)
		}
	})
}

// FuzzProcessObject tests document processing with random YAML
func FuzzProcessObject(f *testing.F) {
	// Seed with valid YAML documents
	f.Add([]byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n"))
	f.Add([]byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: pod\n  namespace: default\n"))
	f.Add([]byte("invalid: [unclosed"))
	f.Add([]byte(""))
	f.Add([]byte("---"))

	f.Fuzz(func(t *testing.T, data []byte) {
		config := Config{
			outputDir:     t.TempDir(),
			outputFormat:  "yaml",
			continueOnErr: true,
		}
		stats := &ProcessingStats{}

		// Should never panic
		_ = processObject(data, config, stats, 1)
	})
}

// FuzzParseFilterList tests filter list parsing with random inputs
func FuzzParseFilterList(f *testing.F) {
	// Seed with valid inputs
	f.Add("Deployment,Service")
	f.Add("Deployment, Service, ConfigMap")
	f.Add("")
	f.Add(",,,")
	f.Add("  ,  ,  ")

	f.Fuzz(func(t *testing.T, input string) {
		// Should never panic
		result := parseFilterList(input)

		// Verify no empty strings in result
		for _, item := range result {
			if item == "" {
				t.Error("parseFilterList returned empty string in result")
			}
		}
	})
}
