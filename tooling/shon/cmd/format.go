/*
Copyright © 2025 sottey

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sottey/shon/tooling/shon/pkg"
	"github.com/spf13/cobra"
)

var (
	minify bool
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Beautify or minimfy SHON",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.DebugPrint("Starting format", Verbose)

		if InputFile == "" {
			fmt.Println("No input file specified. Cancelling.")
			return
		}

		data, err := os.ReadFile(InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read SHON file: %v\n", err)
			os.Exit(1)
		}

		lines := strings.Split(string(data), "\n")
		var out strings.Builder
		level := 0
		inMultilineComment := false

		for _, line := range lines {
			pkg.DebugPrint("Trimming Space...", Verbose)
			trimmed := strings.TrimSpace(line)

			pkg.DebugPrint("Removing comments...", Verbose)
			if strings.HasPrefix(trimmed, "/*") {
				inMultilineComment = true
			}

			if inMultilineComment {
				if !minify {
					out.WriteString(pkg.IndentLine(level, Indentation, trimmed) + "\n")
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

			pkg.DebugPrint("Cleaning up brackets", Verbose)
			openBraces := strings.Count(trimmed, "{") + strings.Count(trimmed, "[")
			closeBraces := strings.Count(trimmed, "}") + strings.Count(trimmed, "]")

			if closeBraces > openBraces {
				level--
			}

			if minify {
				out.WriteString(strings.TrimSpace(trimmed))
			} else {
				out.WriteString(pkg.IndentLine(level, Indentation, trimmed) + "\n")
			}

			if openBraces > closeBraces {
				level++
			}
		}

		if OutputFile == "" {
			fmt.Println(out.String())
		} else {
			err := os.WriteFile(OutputFile, []byte(out.String()), 0644)
			if err != nil {
				log.Fatalf("failed to write to file %s: %v", OutputFile, err)
				return
			}
		}

		pkg.DebugPrint("Format complete", Verbose)
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().BoolVarP(&minify, "minify", "m", false, "Minify shon (remove whitespace)")
}
