package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var inputPath, outputPath string
	var sortKeys bool
	flag.StringVar(&inputPath, "in", "", "Input JSON file")
	flag.StringVar(&outputPath, "out", "", "Base output name")
	flag.BoolVar(&sortKeys, "sort", false, "Sort fields alphabetically")
	flag.Parse()

	if inputPath == "" || outputPath == "" {
		fmt.Println("Usage: json2shon -in file.json -out name [-sort]")
		os.Exit(1)
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var input interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
		os.Exit(1)
	}

	shonBody := convertToShon(input, 1, sortKeys)
	shon := fmt.Sprintf("$schema: \"%s.shos\"\n\n@data {%s\n}", filepath.Base(outputPath), shonBody)

	if err := os.WriteFile(outputPath+".shon", []byte(shon), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write SHON file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✔ SHON file written to %s.shon\n", outputPath)
}

func convertToShon(v interface{}, indent int, sortKeys bool) string {
	ind := strings.Repeat("    ", indent)
	switch val := v.(type) {
	case map[string]interface{}:
		var keys []string
		for k := range val {
			keys = append(keys, k)
		}
		if sortKeys {
			sort.Strings(keys)
		}
		var sb strings.Builder
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("\n%s%s: %s,", ind, k, convertToShon(val[k], indent+1, sortKeys)))
		}
		return sb.String()
	case []interface{}:
		var sb strings.Builder
		sb.WriteString("[")
		for i, item := range val {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(convertToShon(item, 0, sortKeys))
		}
		sb.WriteString("]")
		return sb.String()
	case string:
		// Detect decimal and timestamp patterns for SHON 0.6
		if _, err := strconv.ParseFloat(val, 64); err == nil && strings.Contains(val, ".") {
			return fmt.Sprintf("$decimal(\"%s\")", val)
		}
		if strings.Contains(val, "T") && strings.Contains(val, ":") {
			return fmt.Sprintf("$timestamp(\"%s\")", val)
		}
		return fmt.Sprintf("\"%s\"", val)
	case float64:
		if float64(int(val)) == val {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("$decimal(\"%f\")", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("\"%v\"", val)
	}
}
