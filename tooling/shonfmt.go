package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputPath string
	var minify bool
	flag.StringVar(&inputPath, "in", "", "Input SHON file")
	flag.BoolVar(&minify, "minify", false, "Minify the SHON output")
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Usage: shonfmt -in file.shon [-minify]")
		os.Exit(1)
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read SHON file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	var out strings.Builder
	indent := 0
	inMultilineComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "/*") {
			inMultilineComment = true
		}

		if inMultilineComment {
			if !minify {
				out.WriteString(indentLine(indent, trimmed) + "\n")
			}
			if strings.Contains(trimmed, "*/") {
				inMultilineComment = false
			}
			continue
		}

		if trimmed == "" {
			if !minify {
				out.WriteString("\n")
			}
			continue
		}

		openBraces := strings.Count(trimmed, "{") + strings.Count(trimmed, "[")
		closeBraces := strings.Count(trimmed, "}") + strings.Count(trimmed, "]")

		if closeBraces > openBraces {
			indent--
		}

		if minify {
			out.WriteString(strings.TrimSpace(trimmed))
		} else {
			out.WriteString(indentLine(indent, trimmed) + "\n")
		}

		if openBraces > closeBraces {
			indent++
		}
	}

	fmt.Println(out.String())
}

func indentLine(level int, line string) string {
	return strings.Repeat("    ", level) + line
}
