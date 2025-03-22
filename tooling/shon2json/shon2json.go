
package main

import (
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
)

func main() {
    var inputPath, outputPath string
    var sortKeys bool
    flag.StringVar(&inputPath, "in", "", "Input SHON file")
    flag.StringVar(&outputPath, "out", "", "Output JSON file")
    flag.BoolVar(&sortKeys, "sort", false, "Sort top-level object keys alphabetically (default: false)")
    flag.Parse()

    if inputPath == "" || outputPath == "" {
        fmt.Println("Usage: shon2json -in input.shon -out output.json [--sort]")
        os.Exit(1)
    }

    input, err := os.ReadFile(inputPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read SHON file: %v\n", err)
        os.Exit(1)
    }

    jsonText, err := shonToJSON(string(input), sortKeys)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to convert SHON to JSON: %v\n", err)
        os.Exit(1)
    }

    err = os.WriteFile(outputPath, []byte(jsonText), 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to write output file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✔ JSON written to %s\n", outputPath)
}

func shonToJSON(shon string, sortKeys bool) (string, error) {
    lines := strings.Split(shon, "\n")
    var clean []string
    insideNamespace := false

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "$schema") || trimmed == "" {
            continue
        }
        if strings.HasPrefix(trimmed, "@") && strings.Contains(trimmed, "{") {
            insideNamespace = true
            continue
        }
        if insideNamespace && trimmed == "}" {
            insideNamespace = false
            continue
        }

        if insideNamespace {
            line = regexp.MustCompile(`//.*$`).ReplaceAllString(line, "")
            clean = append(clean, line)
        }
    }

    raw := "{" + strings.Join(clean, "\n") + "}}"

    // Quote top-level and nested keys
    raw = strings.ReplaceAll(raw, "'", "\"")
    raw = regexp.MustCompile(`(?m)^(\s*)([a-zA-Z0-9_.$]+)\s*:`).ReplaceAllString(raw, `$1"$2":`)
    raw = regexp.MustCompile(`^{\s*([a-zA-Z0-9_.$]+)\s*:`).ReplaceAllString(raw, `{"$1":`)
    raw = regexp.MustCompile(`([{,]\s*)([a-zA-Z0-9_.$]+)\s*:`).ReplaceAllString(raw, `$1"$2":`)

    // Fix missing commas between top-level keys
    raw = regexp.MustCompile(`}\s*\n\s*"[^"]+":`).ReplaceAllStringFunc(raw, func(match string) string {
        return "},\n" + match[strings.Index(match, "\""):]
    })
    raw = regexp.MustCompile(`}\s*("[^"]+":)`).ReplaceAllString(raw, "}, $1")

    raw = strings.ReplaceAll(raw, ",}", "}")
    raw = strings.ReplaceAll(raw, ",]", "]")

    var parsed map[string]interface{}
    if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
        return "", errors.New("invalid SHON structure after conversion: " + err.Error())
    }

    if sortKeys {
        sorted := make(map[string]interface{}, len(parsed))
        keys := make([]string, 0, len(parsed))
        for k := range parsed {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        for _, k := range keys {
            sorted[k] = parsed[k]
        }
        parsed = sorted
    }

    jsonBytes, err := json.MarshalIndent(parsed, "", "  ")
    if err != nil {
        return "", err
    }

    return string(jsonBytes), nil
}
