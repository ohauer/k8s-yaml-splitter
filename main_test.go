package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSanitizeFilename tests the filename sanitization function
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "deployment", "deployment"},
		{"with-dash", "my-app", "my-app"},
		{"path-traversal-1", "../etc/passwd", "passwd"},
		{"path-traversal-2", "../../etc/passwd", "passwd"},
		{"path-traversal-3", "../../../etc/passwd", "passwd"},
		{"path-traversal-4", "....//....//etc//passwd", "passwd"},
		{"absolute-path", "/etc/passwd", "passwd"},
		{"windows-path", "C:\\Windows\\System32", "C--Windows-System32"}, // filepath.Base on Unix doesn't split Windows paths
		{"unc-path", "\\\\server\\share\\file", "--server-share-file"},   // filepath.Base on Unix doesn't split UNC paths
		{"empty", "", "unnamed"},
		{"dot", ".", "unnamed"},
		{"double-dot", "..", "unnamed"},
		{"colons", "app:v1.0", "app-v1.0"},
		{"spaces", "my app", "my-app"},
		{"slashes", "name/with/slashes", "slashes"},
		{"mixed", "../path:with/special chars", "special-chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}

			// Verify no path separators in result
			if filepath.Base(result) != result {
				t.Errorf("result contains path separators: %q", result)
			}

			// Verify no path traversal patterns
			if result != "unnamed" && (result == "." || result == "..") {
				t.Errorf("result is unsafe: %q", result)
			}
		})
	}
}

// TestGetOutputPath tests output path generation and validation
func TestGetOutputPath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		config    Config
		base      baseObject
		ext       string
		shouldErr bool
	}{
		{
			name: "normal-flat",
			config: Config{
				outputDir:     tmpDir,
				namespaceDirs: false,
			},
			base: baseObject{
				Kind: "Deployment",
				Meta: struct {
					Namespace string `yaml:"namespace"`
					Name      string `yaml:"name"`
				}{
					Namespace: "default",
					Name:      "app",
				},
			},
			ext:       ".yaml",
			shouldErr: false,
		},
		{
			name: "normal-namespace-dirs",
			config: Config{
				outputDir:     tmpDir,
				namespaceDirs: true,
			},
			base: baseObject{
				Kind: "Service",
				Meta: struct {
					Namespace string `yaml:"namespace"`
					Name      string `yaml:"name"`
				}{
					Namespace: "production",
					Name:      "web",
				},
			},
			ext:       ".yaml",
			shouldErr: false,
		},
		{
			name: "cluster-scoped",
			config: Config{
				outputDir:     tmpDir,
				namespaceDirs: true,
			},
			base: baseObject{
				Kind: "Namespace",
				Meta: struct {
					Namespace string `yaml:"namespace"`
					Name      string `yaml:"name"`
				}{
					Name: "test",
				},
			},
			ext:       ".yaml",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getOutputPath(tt.config, tt.base, tt.ext)

			if tt.shouldErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.shouldErr {
				// Verify path is absolute
				if !filepath.IsAbs(result) {
					t.Errorf("path is not absolute: %s", result)
				}

				// Verify path is within output directory
				absOutput, _ := filepath.Abs(tmpDir)
				if !filepath.HasPrefix(result, absOutput+string(filepath.Separator)) {
					t.Errorf("path escapes output directory: %s", result)
				}
			}
		})
	}
}

// TestValidateFileType tests file type validation
func TestValidateFileType(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	symlinkFile := filepath.Join(tmpDir, "symlink.txt")
	if err := os.Symlink(regularFile, symlinkFile); err != nil {
		t.Skip("Cannot create symlinks on this system")
	}

	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")

	tests := []struct {
		name      string
		path      string
		shouldErr bool
	}{
		{"regular-file", regularFile, false},
		{"non-existent", nonExistentFile, false},
		{"symlink", symlinkFile, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileType(tt.path)
			if tt.shouldErr && err == nil {
				t.Errorf("expected error for %s but got none", tt.path)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("unexpected error for %s: %v", tt.path, err)
			}
		})
	}
}

// TestShouldIncludeResource tests resource filtering logic
func TestShouldIncludeResource(t *testing.T) {
	tests := []struct {
		name         string
		kind         string
		includeKinds []string
		excludeKinds []string
		expected     bool
	}{
		{"no-filters", "Deployment", nil, nil, true},
		{"include-match", "Deployment", []string{"Deployment", "Service"}, nil, true},
		{"include-no-match", "ConfigMap", []string{"Deployment", "Service"}, nil, false},
		{"exclude-match", "Secret", nil, []string{"Secret", "ConfigMap"}, false},
		{"exclude-no-match", "Deployment", nil, []string{"Secret", "ConfigMap"}, true},
		{"include-and-exclude", "Deployment", []string{"Deployment"}, []string{"Deployment"}, false},
		{"case-insensitive", "deployment", []string{"Deployment"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIncludeResource(tt.kind, tt.includeKinds, tt.excludeKinds)
			if result != tt.expected {
				t.Errorf("shouldIncludeResource(%q, %v, %v) = %v, want %v",
					tt.kind, tt.includeKinds, tt.excludeKinds, result, tt.expected)
			}
		})
	}
}

// TestParseFilterList tests comma-separated list parsing
func TestParseFilterList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty", "", nil},
		{"single", "Deployment", []string{"Deployment"}},
		{"multiple", "Deployment,Service,ConfigMap", []string{"Deployment", "Service", "ConfigMap"}},
		{"with-spaces", "Deployment, Service, ConfigMap", []string{"Deployment", "Service", "ConfigMap"}},
		{"trailing-comma", "Deployment,Service,", []string{"Deployment", "Service"}},
		{"extra-spaces", "  Deployment  ,  Service  ", []string{"Deployment", "Service"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFilterList(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseFilterList(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("parseFilterList(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestResourceLimits tests resource limit enforcement
func TestResourceLimits(t *testing.T) {
	config := Config{
		outputDir:     t.TempDir(),
		outputFormat:  "yaml",
		continueOnErr: false,
	}

	t.Run("document-size-limit", func(t *testing.T) {
		// Create document larger than MaxDocumentSize
		largeData := make([]byte, MaxDocumentSize+1)
		stats := &ProcessingStats{}

		err := processObject(largeData, config, stats, 1)
		if err == nil {
			t.Error("expected error for oversized document")
		}
		if stats.Errors != 1 {
			t.Errorf("expected 1 error, got %d", stats.Errors)
		}
	})

	t.Run("document-count-limit", func(t *testing.T) {
		stats := &ProcessingStats{Total: MaxDocuments}
		smallData := []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n")

		err := processObject(smallData, config, stats, MaxDocuments+1)
		if err == nil {
			t.Error("expected error for exceeding document count")
		}
		if stats.Errors != 1 {
			t.Errorf("expected 1 error, got %d", stats.Errors)
		}
	})
}

// TestFormatOutput tests YAML and JSON formatting
func TestFormatOutput(t *testing.T) {
	yamlData := []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n")

	t.Run("yaml-no-sort", func(t *testing.T) {
		output, ext, err := formatOutput(yamlData, "yaml", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ext != ".yaml" {
			t.Errorf("expected .yaml extension, got %s", ext)
		}
		if len(output) == 0 {
			t.Error("output is empty")
		}
	})

	t.Run("yaml-with-sort", func(t *testing.T) {
		output, ext, err := formatOutput(yamlData, "yaml", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ext != ".yaml" {
			t.Errorf("expected .yaml extension, got %s", ext)
		}
		if len(output) == 0 {
			t.Error("output is empty")
		}
	})

	t.Run("json", func(t *testing.T) {
		output, ext, err := formatOutput(yamlData, "json", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ext != ".json" {
			t.Errorf("expected .json extension, got %s", ext)
		}
		if len(output) == 0 {
			t.Error("output is empty")
		}
	})
}

// TestProcessObject tests end-to-end document processing
func TestProcessObject(t *testing.T) {
	tmpDir := t.TempDir()
	config := Config{
		outputDir:     tmpDir,
		outputFormat:  "yaml",
		continueOnErr: true,
		createDir:     true,
	}

	t.Run("valid-document", func(t *testing.T) {
		data := []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test-config\n  namespace: default\ndata:\n  key: value\n")
		stats := &ProcessingStats{}

		err := processObject(data, config, stats, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.Processed != 1 {
			t.Errorf("expected 1 processed, got %d", stats.Processed)
		}
	})

	t.Run("missing-kind", func(t *testing.T) {
		data := []byte("apiVersion: v1\nmetadata:\n  name: test\n")
		stats := &ProcessingStats{}

		err := processObject(data, config, stats, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.Skipped != 1 {
			t.Errorf("expected 1 skipped, got %d", stats.Skipped)
		}
	})

	t.Run("malformed-yaml", func(t *testing.T) {
		data := []byte("invalid: [unclosed")
		stats := &ProcessingStats{}

		err := processObject(data, config, stats, 1)
		if err == nil {
			t.Error("expected error for malformed YAML")
		}
		if stats.Errors != 1 {
			t.Errorf("expected 1 error, got %d", stats.Errors)
		}
	})
}

// TestRandomString tests random string generation
func TestRandomString(t *testing.T) {
	// Generate multiple random strings
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s := randomString(8)
		if len(s) != 8 {
			t.Errorf("expected length 8, got %d", len(s))
		}
		if seen[s] {
			t.Errorf("duplicate random string: %s", s)
		}
		seen[s] = true
	}
}

// TestWriteFileAtomic tests atomic file writing
func TestWriteFileAtomic(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("successful-write", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		data := []byte("test content")

		err := writeFileAtomic(path, data, 0644)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file exists and has correct content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(content) != string(data) {
			t.Errorf("content mismatch: got %q, want %q", content, data)
		}

		// Verify no temp files left behind
		matches, _ := filepath.Glob(filepath.Join(tmpDir, "*.tmp.*"))
		if len(matches) > 0 {
			t.Errorf("temp files not cleaned up: %v", matches)
		}
	})

	t.Run("overwrite-existing", func(t *testing.T) {
		path := filepath.Join(tmpDir, "overwrite.txt")

		// Write initial content
		os.WriteFile(path, []byte("old"), 0644)

		// Overwrite atomically
		newData := []byte("new content")
		err := writeFileAtomic(path, newData, 0644)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify new content
		content, _ := os.ReadFile(path)
		if string(content) != string(newData) {
			t.Errorf("content not updated: got %q, want %q", content, newData)
		}
	})
}

// TestValidateDirectory tests directory validation
func TestValidateDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid-directory", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "valid")
		os.Mkdir(dir, 0755)

		err := validateDirectory(dir)
		if err != nil {
			t.Errorf("unexpected error for valid directory: %v", err)
		}
	})

	t.Run("non-existent", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "nonexistent")

		err := validateDirectory(dir)
		if err != nil {
			t.Errorf("unexpected error for non-existent directory: %v", err)
		}
	})

	t.Run("file-not-directory", func(t *testing.T) {
		file := filepath.Join(tmpDir, "file.txt")
		os.WriteFile(file, []byte("test"), 0644)

		err := validateDirectory(file)
		if err == nil {
			t.Error("expected error for file instead of directory")
		}
	})

	t.Run("symlink-directory", func(t *testing.T) {
		realDir := filepath.Join(tmpDir, "real")
		symlinkDir := filepath.Join(tmpDir, "symlink")
		os.Mkdir(realDir, 0755)

		if err := os.Symlink(realDir, symlinkDir); err != nil {
			t.Skip("Cannot create symlinks on this system")
		}

		err := validateDirectory(symlinkDir)
		if err == nil {
			t.Error("expected error for symlink directory")
		}
	})
}

// TestCleanupOrphanedTempFiles tests temp file cleanup
func TestCleanupOrphanedTempFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some temp files
	os.WriteFile(filepath.Join(tmpDir, "file.tmp.abc123"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.tmp.xyz789"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "normal.yaml"), []byte("test"), 0644)

	cleanupOrphanedTempFiles(tmpDir)

	// Verify temp files removed
	matches, _ := filepath.Glob(filepath.Join(tmpDir, "*.tmp.*"))
	if len(matches) > 0 {
		t.Errorf("temp files not cleaned up: %v", matches)
	}

	// Verify normal file still exists
	if _, err := os.Stat(filepath.Join(tmpDir, "normal.yaml")); err != nil {
		t.Error("normal file was removed")
	}
}

// TestSplitYAMLDocument tests YAML document splitting
func TestSplitYAMLDocument(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "single-document",
			input:    "apiVersion: v1\nkind: Pod",
			expected: 1,
		},
		{
			name:     "two-documents",
			input:    "apiVersion: v1\n---\napiVersion: v2",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			scanner.Split(splitYAMLDocument)

			count := 0
			for scanner.Scan() {
				count++
			}

			if count != tt.expected {
				t.Errorf("got %d documents, want %d", count, tt.expected)
			}
		})
	}
}

// TestSortMapKeys tests recursive map key sorting
func TestSortMapKeys(t *testing.T) {
	t.Run("simple-map", func(t *testing.T) {
		input := map[string]interface{}{
			"z": "last",
			"a": "first",
			"m": "middle",
		}
		result := sortMapKeys(input)
		sorted := result.(map[string]interface{})

		// Just verify it returns a map with same keys
		if len(sorted) != 3 {
			t.Errorf("map size changed: got %d, want 3", len(sorted))
		}
		if sorted["a"] != "first" || sorted["z"] != "last" {
			t.Error("values changed during sorting")
		}
	})

	t.Run("nested-map", func(t *testing.T) {
		input := map[string]interface{}{
			"outer": map[string]interface{}{
				"z": "last",
				"a": "first",
			},
		}
		result := sortMapKeys(input)
		sorted := result.(map[string]interface{})
		nested := sorted["outer"].(map[string]interface{})

		if len(nested) != 2 {
			t.Errorf("nested map size changed: got %d, want 2", len(nested))
		}
		if nested["a"] != "first" || nested["z"] != "last" {
			t.Error("nested values changed during sorting")
		}
	})

	t.Run("array", func(t *testing.T) {
		input := []interface{}{
			map[string]interface{}{"z": 1, "a": 2},
			"string",
			123,
		}
		result := sortMapKeys(input)
		arr := result.([]interface{})

		if len(arr) != 3 {
			t.Errorf("array length changed: got %d, want 3", len(arr))
		}
	})
}

// TestReadInput tests input reading from file and stdin
func TestReadInput(t *testing.T) {
	t.Run("read-file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test.yaml")
		content := []byte("test content")
		os.WriteFile(tmpFile, content, 0644)

		result, err := readInput(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != string(content) {
			t.Errorf("content mismatch: got %q, want %q", result, content)
		}
	})

	t.Run("file-not-found", func(t *testing.T) {
		_, err := readInput("/nonexistent/file.yaml")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}
