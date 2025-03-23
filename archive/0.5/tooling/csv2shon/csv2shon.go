
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Record struct {
	Fields map[string]string
}

func main() {
	var inputPath, outputPath, namespace string
	flag.StringVar(&inputPath, "in", "", "Input CSV file")
	flag.StringVar(&outputPath, "out", "", "Output SHON file name (without extension)")
	flag.StringVar(&namespace, "ns", "data", "Top-level namespace name")
	flag.Parse()

	if inputPath == "" || outputPath == "" {
		fmt.Println("Usage: csv2shon -in input.csv -out outputName [-ns namespace]")
		os.Exit(1)
	}

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open CSV: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read CSV: %v\n", err)
		os.Exit(1)
	}

	if len(rows) < 2 {
		fmt.Fprintln(os.Stderr, "CSV must have a header and at least one data row")
		os.Exit(1)
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
		if len(values) < len(rows)-1 { // column has at least one duplicate
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
	builder.WriteString(fmt.Sprintf("@%s {\n", namespace))

	builder.WriteString("    records: [\n")
	for _, rec := range records {
		builder.WriteString("        {\n")
		for k, v := range rec.Fields {
			builder.WriteString(fmt.Sprintf("            %s: %s,\n", k, v))
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ],\n")

	// Add @refs blocks for columns that used them
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
	shonPath := outputPath + ".shon"
	err = os.WriteFile(shonPath, []byte(builder.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write SHON file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✔ SHON written to %s\n", shonPath)
}
