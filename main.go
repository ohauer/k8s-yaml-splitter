package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"
)

// yamlSeparator defines the YAML document separator used by Kubernetes
const yamlSeparator = "\n---"

// Resource limits to prevent DoS attacks
const (
	MaxDocuments       = 10000
	MaxDocumentSize    = 10 * 1024 * 1024   // 10MB per document
	MaxTotalOutputSize = 1024 * 1024 * 1024 // 1GB total output
	MaxInputFileSize   = 100 * 1024 * 1024  // 100MB input file
)

// Exit codes
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitInvalidArgs  = 2
)

// Version information set at build time via ldflags
var (
	Version   = "dev"
	GitCommit = "dev"
	BuildDate = "unknown"
)

// Config holds all command-line configuration options
type Config struct {
	inputFile     string   // Path to input YAML file or "-" for stdin
	outputDir     string   // Directory to write split files
	outputFormat  string   // Output format: "yaml" or "json"
	sortKeys      bool     // Whether to sort keys in output
	createDir     bool     // Whether to create output directory if missing
	dryRun        bool     // Whether to run in dry-run mode (no file writes)
	namespaceDirs bool     // Whether to organize files by namespace directories
	includeKinds  []string // Resource kinds to include (empty = all)
	excludeKinds  []string // Resource kinds to exclude
	continueOnErr bool     // Whether to continue processing on errors
}

// ProcessingStats tracks processing statistics
type ProcessingStats struct {
	Total           int
	Processed       int
	Skipped         int
	Errors          int
	TotalOutputSize int64 // Track total output size for DoS prevention
}

// baseObject represents the minimal structure needed to identify and name Kubernetes resources
// We only unmarshal the fields we need for file naming to preserve original formatting
type baseObject struct {
	bytes  []byte // Original raw bytes for content preservation
	Kind   string `yaml:"kind"`
	ApiVer string `yaml:"apiVersion"`
	Meta   struct {
		Namespace string `yaml:"namespace"`
		Name      string `yaml:"name"`
	} `yaml:"metadata"`
}

// splitYAMLDocument is a bufio.SplitFunc that splits YAML streams into individual documents
// This function is adapted from kubernetes/apimachinery to handle multi-document YAML files
// It looks for "---" separators and properly handles edge cases like EOF and incomplete documents
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// Found a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// No more data after separator
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil // Need more data
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			// Found complete separator with newline
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil // Incomplete separator
	}

	// No separator found
	if atEOF {
		return len(data), data, nil // Return remaining data
	}
	return 0, nil, nil // Request more data
}

// sortMapKeys recursively sorts all map keys in a data structure
// This ensures consistent output ordering for both YAML and JSON formats
func sortMapKeys(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Create new sorted map
		sorted := make(map[string]interface{})
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// Recursively sort nested structures
		for _, k := range keys {
			sorted[k] = sortMapKeys(v[k])
		}
		return sorted
	case []interface{}:
		// Sort elements in arrays
		for i, item := range v {
			v[i] = sortMapKeys(item)
		}
		return v
	default:
		// Return primitive values unchanged
		return v
	}
}

// formatOutput converts raw YAML data to the requested output format
// It handles three scenarios:
// 1. JSON output: Always sorts keys and converts to JSON
// 2. YAML with sorting: Unmarshals, sorts, and re-marshals
// 3. YAML without sorting: Preserves original formatting to maintain comments and structure
func formatOutput(data []byte, format string, sortKeys bool) ([]byte, string, error) {
	switch format {
	case "json":
		// JSON output is always sorted for consistency
		var obj interface{}
		if err := yaml.Unmarshal(data, &obj); err != nil {
			return nil, "", err
		}
		obj = sortMapKeys(obj)
		output, err := json.MarshalIndent(obj, "", "  ")
		return output, ".json", err
	default:
		if sortKeys {
			// Sort YAML keys by unmarshaling and re-marshaling
			var obj interface{}
			if err := yaml.Unmarshal(data, &obj); err != nil {
				return nil, "", err
			}
			obj = sortMapKeys(obj)
			var buf bytes.Buffer
			buf.WriteString("---\n")
			encoder := yaml.NewEncoder(&buf)
			encoder.SetIndent(2)
			if err := encoder.Encode(obj); err != nil {
				return nil, "", err
			}
			encoder.Close()
			return buf.Bytes(), ".yaml", nil
		} else {
			// Preserve original formatting when not sorting
			// This maintains comments, spacing, and original structure
			var buf bytes.Buffer
			buf.WriteString("---\n")
			buf.Write(bytes.TrimRight(data, "\n"))
			buf.WriteString("\n")
			return buf.Bytes(), ".yaml", nil
		}
	}
}

// sanitizeFilename removes characters that are invalid in filenames and prevents path traversal
// Kubernetes resource names can contain characters like "/" and ":" that are problematic in filenames
func sanitizeFilename(name string) string {
	// Extract base filename FIRST to prevent path traversal
	name = filepath.Base(name)
	// Remove path traversal attempts
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, " ", "-")
	// Ensure filename is not empty after sanitization
	if name == "" || name == "." || name == ".." {
		name = "unnamed"
	}
	return name
}

// validateFileType checks if a path is safe to write to (not a symlink or special file)
func validateFileType(path string) error {
	info, err := os.Lstat(path) // Don't follow symlinks
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, OK to create
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Reject symlinks
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to write to symlink: %s", path)
	}

	// Reject special files (devices, pipes, sockets)
	if !info.Mode().IsRegular() {
		return fmt.Errorf("refusing to write to special file: %s", path)
	}

	return nil
}

// randomString generates a random hex string of length n
func randomString(n int) string {
	bytes := make([]byte, n/2+1)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based if crypto/rand fails
		return fmt.Sprintf("%d", os.Getpid())
	}
	return hex.EncodeToString(bytes)[:n]
}

// writeFileAtomic writes data to a file atomically using a random temp file
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmpFile := path + ".tmp." + randomString(8)

	// Clean up temp file on error
	defer func() {
		if _, err := os.Stat(tmpFile); err == nil {
			os.Remove(tmpFile)
		}
	}()

	// Write to temp file
	if err := os.WriteFile(tmpFile, data, perm); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// cleanupOrphanedTempFiles removes orphaned temp files from previous runs
func cleanupOrphanedTempFiles(dir string) {
	pattern := filepath.Join(dir, "*.tmp.*")
	matches, _ := filepath.Glob(pattern)
	for _, file := range matches {
		os.Remove(file) // Ignore errors
	}
}

// validateDirectory checks if a path is a valid directory (not a symlink)
func validateDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, OK to create
		}
		return fmt.Errorf("failed to stat directory: %w", err)
	}

	// Reject symlinks
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to use symlink as directory: %s", path)
	}

	// Must be a directory
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	return nil
}

// shouldIncludeResource determines if a resource should be processed based on filters
func shouldIncludeResource(kind string, includeKinds, excludeKinds []string) bool {
	// Check exclude list first
	for _, exclude := range excludeKinds {
		if strings.EqualFold(kind, exclude) {
			return false
		}
	}

	// If include list is empty, include all (except excluded)
	if len(includeKinds) == 0 {
		return true
	}

	// Check include list
	for _, include := range includeKinds {
		if strings.EqualFold(kind, include) {
			return true
		}
	}

	return false
}

// getOutputPath generates the output file path based on configuration
func getOutputPath(config Config, base baseObject, ext string) (string, error) {
	var fileName string
	var dirPath string

	if config.namespaceDirs {
		// Organize by namespace directories
		if base.Meta.Namespace != "" {
			dirPath = filepath.Join(config.outputDir, sanitizeFilename(base.Meta.Namespace))
		} else {
			dirPath = filepath.Join(config.outputDir, "cluster-scoped")
		}
		// Simpler filename without namespace prefix
		fileName = fmt.Sprintf("%s-%s%s",
			sanitizeFilename(base.Kind),
			sanitizeFilename(base.Meta.Name),
			ext)
	} else {
		// Original flat structure
		dirPath = config.outputDir
		if base.Meta.Namespace != "" {
			fileName = fmt.Sprintf("%s-%s-%s%s",
				sanitizeFilename(base.Kind),
				sanitizeFilename(base.Meta.Namespace),
				sanitizeFilename(base.Meta.Name),
				ext)
		} else {
			fileName = fmt.Sprintf("%s-%s%s",
				sanitizeFilename(base.Kind),
				sanitizeFilename(base.Meta.Name),
				ext)
		}
	}

	filePath := filepath.Join(dirPath, fileName)

	// Validate path stays within output directory
	absOutput, err := filepath.Abs(config.outputDir)
	if err != nil {
		return "", fmt.Errorf("invalid output directory: %w", err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Ensure path is within output directory
	if !strings.HasPrefix(absPath, absOutput+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes output directory: %s", absPath)
	}

	return absPath, nil
}

// processObject handles a single YAML document from the input stream
// It extracts metadata for naming, formats the output, and writes the file
func processObject(data []byte, config Config, stats *ProcessingStats, docNum int) error {
	stats.Total++

	// Enforce document count limit
	if stats.Total > MaxDocuments {
		stats.Errors++
		return fmt.Errorf("exceeded maximum documents (%d)", MaxDocuments)
	}

	// Enforce document size limit
	if len(data) > MaxDocumentSize {
		stats.Errors++
		return fmt.Errorf("document %d exceeds maximum size (%d bytes)", docNum, MaxDocumentSize)
	}

	var base baseObject
	base.bytes = data

	// Parse only the metadata we need for file naming
	if err := yaml.Unmarshal(data, &base); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: malformed YAML: %w", docNum, err)
	}

	// Skip documents that don't have required Kubernetes fields
	if base.Kind == "" || base.ApiVer == "" {
		stats.Skipped++
		if !config.dryRun {
			fmt.Printf("Skipped document %d: missing kind or apiVersion\n", docNum)
		}
		return nil
	}

	// Apply resource filtering
	if !shouldIncludeResource(base.Kind, config.includeKinds, config.excludeKinds) {
		stats.Skipped++
		if !config.dryRun {
			fmt.Printf("Filtered out: %s | %s\n", base.Kind, base.Meta.Name)
		}
		return nil
	}

	// Format the output according to user preferences
	output, ext, err := formatOutput(data, config.outputFormat, config.sortKeys)
	if err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: formatting error: %w", docNum, err)
	}

	// Enforce total output size limit
	if stats.TotalOutputSize+int64(len(output)) > MaxTotalOutputSize {
		stats.Errors++
		return fmt.Errorf("exceeded maximum total output size (%d bytes)", MaxTotalOutputSize)
	}

	// Generate output path with validation
	filePath, err := getOutputPath(config, base, ext)
	if err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: invalid output path: %w", docNum, err)
	}

	// Display resource information
	fmt.Printf("Found! type: %s | apiVersion: %s | name: %s | namespace: %s\n",
		base.Kind, base.ApiVer, base.Meta.Name, base.Meta.Namespace)

	// In dry-run mode, only show what would be written
	if config.dryRun {
		fmt.Printf("==> DryRun: Writing %s\n", filePath)
		stats.Processed++
		return nil
	}

	// Create directory if needed
	dirPath := filepath.Dir(filePath)
	if err := validateDirectory(dirPath); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: %w", docNum, err)
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: failed to create directory %s: %w", docNum, dirPath, err)
	}
	// Validate directory after creation
	if err := validateDirectory(dirPath); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: directory validation failed after creation: %w", docNum, err)
	}

	// Validate file type before writing (reject symlinks and special files)
	if err := validateFileType(filePath); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: %w", docNum, err)
	}

	// Write the formatted output to file atomically
	if err := writeFileAtomic(filePath, output, 0644); err != nil {
		stats.Errors++
		return fmt.Errorf("document %d: failed to write file %s: %w", docNum, filePath, err)
	}

	fmt.Printf("* Writing %s\n", filePath)
	fmt.Printf("* Wrote %d bytes to %s\n", len(output), filePath)
	stats.Processed++
	stats.TotalOutputSize += int64(len(output))
	return nil
}

// parseFilterList parses comma-separated filter list
func parseFilterList(filterStr string) []string {
	if filterStr == "" {
		return nil
	}
	parts := strings.Split(filterStr, ",")
	var result []string
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// printStats displays processing statistics
func printStats(stats ProcessingStats) {
	fmt.Printf("\n=== Processing Summary ===\n")
	fmt.Printf("Total documents: %d\n", stats.Total)
	fmt.Printf("Processed: %d\n", stats.Processed)
	fmt.Printf("Skipped: %d\n", stats.Skipped)
	fmt.Printf("Errors: %d\n", stats.Errors)
}

// readInput reads content from file or stdin based on the input parameter
func readInput(inputFile string) ([]byte, error) {
	if inputFile == "-" {
		// Read from stdin
		return io.ReadAll(os.Stdin)
	}
	// Read from file
	return os.ReadFile(inputFile)
}

// printHelp displays usage information
func printHelp() {
	fmt.Printf("Usage: %s -f <input-file> <output-dir>\n", os.Args[0])
	fmt.Printf("   or: %s -f - <output-dir>  (read from stdin)\n\n", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  help, -h, --help     Show this help message")
	fmt.Println("  version, -v          Show version information")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// printVersion displays version information
func printVersion() {
	fmt.Printf("k8s-yaml-splitter %s\n", Version)
	fmt.Printf("  Git Commit: %s\n", GitCommit)
	fmt.Printf("  Build Date: %s\n", BuildDate)
}

// main function handles command-line parsing and orchestrates the splitting process
func main() {
	// Handle subcommands before flag parsing
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			printHelp()
			os.Exit(0)
		case "version", "-v":
			printVersion()
			os.Exit(0)
		}
	}

	var config Config
	var includeFilter, excludeFilter string
	var showVersion bool

	// Define command-line flags
	flag.StringVar(&config.inputFile, "f", "", "Input file (use '-' for stdin)")
	flag.StringVar(&config.outputFormat, "o", "yaml", "Output format (yaml or json)")
	flag.BoolVar(&config.sortKeys, "s", false, "Sort keys in output")
	flag.BoolVar(&config.createDir, "d", false, "Create output directory if it doesn't exist")
	flag.BoolVar(&config.dryRun, "dry-run", false, "Dry run mode")
	flag.BoolVar(&config.namespaceDirs, "namespace-dirs", false, "Organize output files by namespace directories")
	flag.StringVar(&includeFilter, "include", "", "Comma-separated list of resource kinds to include (e.g., 'Deployment,Service')")
	flag.StringVar(&excludeFilter, "exclude", "", "Comma-separated list of resource kinds to exclude (e.g., 'Secret,ConfigMap')")
	flag.BoolVar(&config.continueOnErr, "continue-on-error", true, "Continue processing on individual document errors")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	// Handle version flag (for backward compatibility)
	if showVersion {
		printVersion()
		os.Exit(0)
	}

	// Parse filter lists
	config.includeKinds = parseFilterList(includeFilter)
	config.excludeKinds = parseFilterList(excludeFilter)

	// Validate required parameters
	args := flag.Args()
	if config.inputFile == "" || len(args) != 1 {
		fmt.Printf("Usage: %s -f <input-file> <output-dir>\n", os.Args[0])
		fmt.Printf("   or: %s -f - <output-dir>  (read from stdin)\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(ExitInvalidArgs)
	}

	config.outputDir = args[0]

	// Clean up orphaned temp files from previous runs
	if !config.dryRun {
		cleanupOrphanedTempFiles(config.outputDir)
	}

	// Validate output format
	if config.outputFormat != "yaml" && config.outputFormat != "json" {
		log.Fatal("Output format must be 'yaml' or 'json'")
	}

	// Validate input file exists (unless reading from stdin)
	if config.inputFile != "-" {
		info, err := os.Stat(config.inputFile)
		if os.IsNotExist(err) {
			log.Fatalf("Input file %s does not exist", config.inputFile)
		}
		if err != nil {
			log.Fatalf("Failed to stat input file: %v", err)
		}

		// Validate input file is a regular file
		if !info.Mode().IsRegular() {
			log.Fatalf("Input file must be a regular file, not a directory or special file")
		}

		// Enforce input file size limit
		if info.Size() > MaxInputFileSize {
			log.Fatalf("Input file size (%d bytes) exceeds maximum (%d bytes)", info.Size(), MaxInputFileSize)
		}
	}

	// Handle output directory creation or validation
	if config.createDir {
		if err := os.MkdirAll(config.outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	} else if _, err := os.Stat(config.outputDir); os.IsNotExist(err) {
		log.Fatalf("Output directory %s does not exist (use -d to create)", config.outputDir)
	}

	// Read input content
	content, err := readInput(config.inputFile)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, exiting...")
		os.Exit(130)
	}()

	// Set up scanner to split YAML documents
	// Use large buffer to handle big Kubernetes manifests (like CRDs)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 4096), 1024*1024) // 1MB max buffer
	scanner.Split(splitYAMLDocument)

	// Initialize statistics
	var stats ProcessingStats
	docNum := 0

	// Process each YAML document
	for scanner.Scan() {
		docNum++
		// Copy scanner bytes to avoid mutation issues
		data := make([]byte, len(scanner.Bytes()))
		copy(data, scanner.Bytes())

		if err := processObject(data, config, &stats, docNum); err != nil {
			fmt.Printf("Error processing document %d: %v\n", docNum, err)
			if !config.continueOnErr {
				log.Fatalf("Stopping due to error (use -continue-on-error to continue)")
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Print processing statistics
	printStats(stats)

	// Exit with error code if there were processing errors
	if stats.Errors > 0 {
		os.Exit(ExitGeneralError)
	}
}
