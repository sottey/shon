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
	"os"

	"github.com/spf13/cobra"
)

var (
	InputFile   string
	OutputFile  string
	SortKeys    bool
	Verbose     bool
	Indentation int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shon",
	Short: "Conversion and formatting for SHON files",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "File path of input file")
	rootCmd.PersistentFlags().StringVarP(&OutputFile, "output", "o", "", "File path of output file")
	rootCmd.PersistentFlags().BoolVarP(&SortKeys, "sort", "s", false, "If present, keys will be sorted alphabetically")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "If present, additional information will be displayed")
	rootCmd.PersistentFlags().IntVarP(&Indentation, "indentation-size", "n", 4, "Number of spaces to use for each level of indentation (Default: 4)")
}
