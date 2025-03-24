package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "in", "", "Input SHON file")
	flag.StringVar(&outputPath, "out", "", "Output JSON file")
	flag.Parse()

	if inputPath == "" || outputPath == "" {
		fmt.Println("Usage: shon2json -in file.shon -out output.json")
		os.Exit(1)
	}

	content, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read SHON file: %v\n", err)
		os.Exit(1)
	}

	raw := string(content)

	// Remove $schema line
	raw = regexp.MustCompile(`(?m)^\$schema:.*\n?`).ReplaceAllString(raw, "")

	// Remove multi-line and single-line comments
	raw = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`(?m)^\s*//.*`).ReplaceAllString(raw, "")

	// Convert typed values
	raw = strings.ReplaceAll(raw, "$decimal(", "")
	raw = strings.ReplaceAll(raw, "$timestamp(", "")
	raw = strings.ReplaceAll(raw, "$tuple(", "[")
	raw = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_]*)\(`).ReplaceAllString(raw, "[") // named tuple

	// Close open parentheses
	raw = strings.ReplaceAll(raw, ")", "]")

	// Replace refs with strings
	raw = regexp.MustCompile(`&([a-zA-Z0-9_\.]+)`).ReplaceAllString(raw, `"$1"`)

	// Quote keys
	raw = regexp.MustCompile(`(?m)^\s*([a-zA-Z_][a-zA-Z0-9_]*):`).ReplaceAllString(raw, `"$1":`)
	raw = regexp.MustCompile(`(?m)([^"\s][a-zA-Z0-9_]*):`).ReplaceAllString(raw, `"$1":`)

	// Remove @namespace {
	raw = regexp.MustCompile(`(?m)^@([a-zA-Z0-9_]+)\s*\{`).ReplaceAllString(raw, `{`)

	// Final cleanup
	raw = strings.TrimSpace(raw)
	if !strings.HasSuffix(raw, "}") {
		raw += "}"
	}

	var parsed interface{}
	err = json.Unmarshal([]byte(raw), &parsed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to convert SHON to JSON: %v\n", err)
		os.Exit(1)
	}

	jsonOut, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath+".json", jsonOut, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✔ JSON written to %s.json\n", outputPath)
}
