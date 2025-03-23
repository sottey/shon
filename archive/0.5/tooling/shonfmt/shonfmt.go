
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputPath, outputPath string
	var useTrailingCommas bool
	var minify bool

	flag.StringVar(&inputPath, "in", "", "Input SHON file to format")
	flag.StringVar(&outputPath, "out", "", "Output file (optional)")
	flag.BoolVar(&useTrailingCommas, "trailing-commas", false, "Add trailing commas to arrays/objects")
	flag.BoolVar(&minify, "minify", false, "Minify the SHON output (no indentation or newlines)")
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Usage: shonfmt -in input.shon [-out output.shon] [--trailing-commas] [--minify]")
		os.Exit(1)
	}

	input, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var formatted string
	if minify {
		formatted = minifySHON(string(input))
	} else {
		formatted = formatSHON(string(input), useTrailingCommas)
	}

	if outputPath != "" {
		err := os.WriteFile(outputPath, []byte(formatted), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✔ SHON written to %s\n", outputPath)
	} else {
		fmt.Println(formatted)
	}
}

func formatSHON(input string, trailingCommas bool) string {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	indentLevel := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Preserve $schema or comments
		if strings.HasPrefix(line, "$schema") || strings.HasPrefix(line, "//") {
			result = append(result, line)
			continue
		}
		if line == "" {
			result = append(result, "")
			continue
		}

		opening := strings.Count(line, "{") + strings.Count(line, "[")
		closing := strings.Count(line, "}") + strings.Count(line, "]")

		if closing > opening {
			indentLevel--
		}

		indent := strings.Repeat("    ", indentLevel)
		formattedLine := indent + line

		if trailingCommas && (strings.HasSuffix(line, "}") || strings.HasSuffix(line, "]")) &&
			!strings.HasSuffix(line, "},") && !strings.HasSuffix(line, "],") {
			formattedLine += ","
		}

		result = append(result, formattedLine)

		if opening > closing {
			indentLevel++
		}
	}
	return strings.Join(result, "\n")
}

func minifySHON(input string) string {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		result = append(result, trimmed)
	}
	return strings.Join(result, "")
}
