package pkg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Record struct {
	Fields map[string]string
}

func ConvertFile(inputPath, outputPath string, sortKeys bool) error {
	inExt := strings.ToLower(filepath.Ext(inputPath))
	outExt := strings.ToLower(filepath.Ext(outputPath))

	switch inExt {
	case ".json":
		switch outExt {
		case ".shon":
			return JsonToShon(inputPath, outputPath, sortKeys)
		}
	case ".shon":
		switch outExt {
		case ".json":
			return ShonToJson(inputPath, outputPath)
		}
	case ".csv":
		switch outExt {
		case ".shon":
			return CSVToShon(inputPath, outputPath)
		}
	}

	return fmt.Errorf("unsupported conversion: %s → %s", inExt, outExt)
}

// ShonToJson converts a SHON file to a valid JSON file.
func ShonToJson(inputPath, outputPath string) error {
	// Read the SHON file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read SHON file: %w", err)
	}
	content := string(data)
	fmt.Println("🧾 Raw SHON input:\n-----------------\n" + content)

	// Step 1: Remove the $schema line
	content = regexp.MustCompile(`(?m)^\$schema:.*\n?`).ReplaceAllString(content, "")

	// Step 2: Extract content inside the first @namespace { ... } block
	nsMatch := regexp.MustCompile(`@[\w]+\s*{([\s\S]*)}`).FindStringSubmatch(content)
	if len(nsMatch) < 2 {
		return fmt.Errorf("failed to extract SHON body from namespace block")
	}
	body := strings.TrimSpace(nsMatch[1])

	// Step 3: Wrap the extracted content in braces to form a JSON-like structure
	wrapped := fmt.Sprintf("{%s}", body)
	fmt.Println("\n🧪 Cleaned & wrapped SHON:\n-----------------\n" + wrapped)

	// Step 4: Flatten the string (remove newlines and tabs)
	flat := strings.ReplaceAll(wrapped, "\n", "")
	flat = strings.ReplaceAll(flat, "\t", "")
	flat = strings.ReplaceAll(flat, "  ", "")
	wrapped = flat

	// Step 5: Replace $timestamp("...") with the inner string value.
	// reTimestamp := regexp.MustCompile(`\$timestamp$begin:math:text$"([^"]+)"$end:math:text$`)
	// wrapped = reTimestamp.ReplaceAllString(wrapped, `"$1"`)

	// Replace $decimal("...") with the inner string value.
	// reDecimal := regexp.MustCompile(`\$decimal$begin:math:text$"([^"]+)"$end:math:text$`)
	// wrapped = reDecimal.ReplaceAllString(wrapped, `"$1"`)

	wrapped = strings.ReplaceAll(wrapped, "$timestamp(", "")
	wrapped = strings.ReplaceAll(wrapped, "$decimal(", "")
	wrapped = strings.ReplaceAll(wrapped, ")", "")

	reTypes := regexp.MustCompile(`\$(?:timestamp|decimal)$begin:math:text$"((?:[^"\\\\]|\\\\.)*)"$end:math:text$`)
	wrapped = reTypes.ReplaceAllString(wrapped, `"$1"`)
	/*
		wrapped = strings.ReplaceAll(wrapped, "$timestamp(", "")
		wrapped = strings.ReplaceAll(wrapped, "$decimal(", "")
		wrapped = strings.ReplaceAll(wrapped, ")", "")
	*/
	fmt.Println("\n🔍 After $type cleanup:\n-----------------\n" + wrapped)

	// Fail fast if any $timestamp or $decimal expressions remain.
	if strings.Contains(wrapped, "$timestamp") || strings.Contains(wrapped, "$decimal") {
		return fmt.Errorf("🔥 $type expressions still present after cleanup")
	}

	// Step 6: Quote unquoted keys using a multiline regex.
	// This regex matches any '{' or ',' followed by optional whitespace and a key, then a colon.
	reKey := regexp.MustCompile(`([{,])\s*([a-zA-Z_][a-zA-Z0-9_]*):`)
	wrapped = reKey.ReplaceAllString(wrapped, `$1"$2":`)
	// Also, ensure the very first key is quoted.
	wrapped = regexp.MustCompile(`^{\s*([a-zA-Z_][a-zA-Z0-9_]*):`).ReplaceAllString(wrapped, `{"$1":`)
	fmt.Println("\n🔧 Final JSON-like string:\n-----------------\n" + wrapped)

	// Step 7: Parse the final string as JSON.
	var parsed interface{}
	if err := json.Unmarshal([]byte(wrapped), &parsed); err != nil {
		return fmt.Errorf("invalid SHON structure after conversion: %w", err)
	}

	// Step 8: Marshal the parsed object as pretty JSON and write it to the output file.
	out, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("✔ JSON file written to %s\n", filepath.Base(outputPath))
	return nil
}

func CSVToShon(inputFile, outputFile string) error {
	if inputFile == "" {
		return fmt.Errorf("no input file specified")
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("CSV must have a header and at least one data row")
	}

	header := rows[0]
	var records []Record
	occurrence := make(map[string]map[string]int)
	refs := make(map[string]map[string]string)
	refColumns := make(map[string]bool)

	// Count occurrences per column
	for _, field := range header {
		occurrence[field] = make(map[string]int)
	}
	for _, row := range rows[1:] {
		for i, val := range row {
			key := header[i]
			occurrence[key][val]++
		}
	}

	// Decide if the column should use refs
	for key, values := range occurrence {
		if len(values) < len(rows)-1 {
			refColumns[key] = true
			refs[key] = make(map[string]string)
			for val := range values {
				refID := fmt.Sprintf("%s_%d", strings.ToLower(key), len(refs[key])+1)
				refs[key][val] = refID
			}
		}
	}

	// Build SHON records
	for _, row := range rows[1:] {
		rec := Record{Fields: make(map[string]string)}
		for i, val := range row {
			key := header[i]
			if refColumns[key] {
				refID := refs[key][val]
				rec.Fields[key] = "&" + strings.ToLower(key) + "." + refID
			} else {
				rec.Fields[key] = fmt.Sprintf("\"%s\"", val)
			}
		}
		records = append(records, rec)
	}

	// Generate SHON
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("$schema: \"%s.shos\"\n\n", outputFile))
	// builder.WriteString(fmt.Sprintf("@%s {\n", namespace))
	builder.WriteString("    records: [\n")
	for _, rec := range records {
		builder.WriteString("        {\n")
		for k, v := range rec.Fields {
			builder.WriteString(fmt.Sprintf("            %s: %s,\n", k, v))
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ],\n")

	// Add @refs blocks
	for key, values := range refs {
		refName := strings.ToLower(key)
		builder.WriteString(fmt.Sprintf("    @%s {\n", refName))
		for val, refID := range values {
			builder.WriteString(fmt.Sprintf("        %s: \"%s\",\n", refID, val))
		}
		builder.WriteString("    },\n")
	}

	builder.WriteString("}\n")

	// Write output
	shonPath := outputFile
	err = os.WriteFile(shonPath, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write SHON file: %v", err)
	}
	fmt.Printf("✔ SHON written to %s\n", shonPath)
	return nil
}

func JsonToShon(inputFile, outputFile string, sortKeys bool) error {
	if inputFile == "" {
		return fmt.Errorf("no input file specified")
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var input interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
		os.Exit(1)
	}

	shonBody := convertToShon(input, 1, sortKeys)
	shon := fmt.Sprintf("$schema: \"%s.shos\"\n\n@data %s", filepath.Base(outputFile), shonBody)

	if outputFile == "" {
		fmt.Println(shon)
	} else {
		if err := os.WriteFile(outputFile, []byte(shon), 0644); err != nil {
			return fmt.Errorf("failed to write SHON file: %v", err)
		}
		fmt.Printf("✔ SHON file written to %s\n", outputFile)
	}

	return nil
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
		sb.WriteString("{")
		for i, k := range keys {
			sb.WriteString(fmt.Sprintf("\n%s%s: %s", ind, k, convertToShon(val[k], indent+1, sortKeys)))
			if i < len(keys)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(fmt.Sprintf("\n%s}", strings.Repeat("    ", indent-1)))
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
