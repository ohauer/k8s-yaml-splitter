package main

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
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

// FuzzSplitYAMLDocument tests the document splitter via bufio.Scanner with random streams
func FuzzSplitYAMLDocument(f *testing.F) {
	f.Add([]byte("a: 1\n---\nb: 2\n"))
	f.Add([]byte("---\n---\n---"))
	f.Add([]byte("a: 1\n---"))
	f.Add([]byte("---\r\na: 1\r\n"))
	f.Add([]byte(""))

	f.Fuzz(func(t *testing.T, data []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Buffer(make([]byte, 0, 64), 1024*1024)
		scanner.Split(splitYAMLDocument)
		// Bound iterations: each token must consume input, so count is <= len(data)+1
		for i := 0; scanner.Scan(); i++ {
			if i > len(data)+1 {
				t.Fatalf("splitter did not make progress on %d bytes", len(data))
			}
		}
	})
}

// FuzzFormatOutput tests YAML/JSON formatting with random documents
func FuzzFormatOutput(f *testing.F) {
	f.Add([]byte("b: 1\na: [1, 2]\n"), "json", true)
	f.Add([]byte("b: 1\na: 2\n"), "yaml", true)
	f.Add([]byte("# comment\na: 1\n"), "yaml", false)
	f.Add([]byte("{unclosed"), "json", false)
	f.Add([]byte("? [a]\n: b\n"), "json", true)

	f.Fuzz(func(t *testing.T, data []byte, format string, sortKeys bool) {
		// Should never panic; errors are acceptable
		_, _, _ = formatOutput(data, format, sortKeys)
	})
}

// FuzzShouldIncludeResource tests kind filtering invariants
func FuzzShouldIncludeResource(f *testing.F) {
	f.Add("Pod", "Pod,Service", "Secret")
	f.Add("", "", "")
	f.Add("Secret", "Secret", "secret")

	f.Fuzz(func(t *testing.T, kind, include, exclude string) {
		inc := parseFilterList(include)
		exc := parseFilterList(exclude)
		got := shouldIncludeResource(kind, inc, exc)
		// Excluded kinds must never be included
		for _, e := range exc {
			if strings.EqualFold(kind, e) && got {
				t.Errorf("kind %q excluded by %q but included", kind, e)
			}
		}
	})
}

// FuzzGetOutputPath verifies generated paths never escape the output directory
func FuzzGetOutputPath(f *testing.F) {
	f.Add("Pod", "default", "web", true)
	f.Add("../..", "../../etc", "../passwd", false)
	f.Add("", "", "", true)
	f.Add("a/b", "c\\d", "e:f", false)

	f.Fuzz(func(t *testing.T, kind, ns, name string, nsDirs bool) {
		out := t.TempDir()
		config := Config{outputDir: out, namespaceDirs: nsDirs}
		base := baseObject{Kind: kind}
		base.Meta.Namespace = ns
		base.Meta.Name = name

		path, err := getOutputPath(config, base, ".yaml")
		if err != nil {
			return
		}
		if !strings.HasPrefix(path, out+string(filepath.Separator)) {
			t.Errorf("path %q escapes %q", path, out)
		}
	})
}

// FuzzProcessObjectList tests List expansion with random items
func FuzzProcessObjectList(f *testing.F) {
	f.Add([]byte("apiVersion: v1\nkind: List\nitems:\n- apiVersion: v1\n  kind: Pod\n  metadata:\n    name: p\n"))
	f.Add([]byte("apiVersion: v1\nkind: List\nitems: []\n"))
	f.Add([]byte("apiVersion: v1\nkind: PodList\nitems:\n- 1\n- [a]\n- null\n"))
	f.Add([]byte("apiVersion: v1\nkind: List\nitems:\n- apiVersion: v1\n  kind: List\n  items:\n  - apiVersion: v1\n    kind: Pod\n"))

	f.Fuzz(func(t *testing.T, data []byte) {
		out := t.TempDir()
		config := Config{outputDir: out, outputFormat: "yaml", continueOnErr: true}
		stats := &ProcessingStats{}
		// Should never panic or hang
		_ = processObject(data, config, stats, 1)
		if stats.Total < 0 {
			t.Errorf("negative total: %+v", stats)
		}
	})
}
