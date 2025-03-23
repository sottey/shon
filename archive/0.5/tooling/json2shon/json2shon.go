
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "sort"
    "strings"
)

func main() {
    var inputPath, outputBase, namespace, schemaPath, objType string
    var sortKeys, compactArrays, trailingCommas bool
    flag.StringVar(&inputPath, "in", "", "Path to input JSON file")
    flag.StringVar(&outputBase, "out", "", "Base name for output files (no extension)")
    flag.StringVar(&namespace, "namespace", "data", "Top-level SHON namespace")
    flag.StringVar(&schemaPath, "schema", "./schema.shos", "Schema path to embed in $schema")
    flag.StringVar(&objType, "type", "", "Optional $type value to annotate")
    flag.BoolVar(&sortKeys, "sort", false, "Sort keys alphabetically (default: false)")
    flag.BoolVar(&compactArrays, "compact-arrays", true, "Compact simple object arrays (default: true)")
    flag.BoolVar(&trailingCommas, "trailing-commas", false, "Add trailing commas to objects and arrays")
    flag.Parse()

    if inputPath == "" || outputBase == "" {
        fmt.Println("Usage: json2shon -in input.json -out output_name [flags]")
        os.Exit(1)
    }

    inputData, err := os.ReadFile(inputPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
        os.Exit(1)
    }

    var jsonObj interface{}
    if err := json.Unmarshal(inputData, &jsonObj); err != nil {
        fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
        os.Exit(1)
    }

    shonBody := convertToShon(jsonObj, 1, objType, compactArrays, !sortKeys, trailingCommas)
    shonBody = strings.TrimPrefix(shonBody, "{\n")
    shonBody = strings.TrimSuffix(shonBody, "}")
    fullOutput := fmt.Sprintf("$schema: \"%s\"\n\n@%s {\n%s\n}\n", schemaPath, namespace, shonBody)

    shonPath := outputBase + ".shon"
    if err := os.WriteFile(shonPath, []byte(fullOutput), 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing SHON file: %v\n", err)
        os.Exit(1)
    }

    shosPath := outputBase + ".shos"
    os.WriteFile(shosPath, []byte("// schema placeholder\n"), 0644)

    fmt.Printf("✔ SHON written to %s\n", filepath.Base(shonPath))
    fmt.Printf("✔ Schema placeholder written to %s\n", filepath.Base(shosPath))
}

func convertToShon(data interface{}, indent int, objType string, compactArrays, preserveOrder, trailingCommas bool) string {
    pad := strings.Repeat("    ", indent)
    switch val := data.(type) {
    case map[string]interface{}:
        var sb strings.Builder
        sb.WriteString("{\n")

        if objType != "" {
            sb.WriteString(fmt.Sprintf("%s$type: \"%s\"", pad, objType))
            if len(val) > 0 || trailingCommas {
                sb.WriteString(",\n")
            } else {
                sb.WriteString("\n")
            }
        }

        keys := make([]string, 0, len(val))
        for k := range val {
            keys = append(keys, k)
        }
        if !preserveOrder {
            sort.Strings(keys)
        }

        for i, k := range keys {
            v := val[k]
            sb.WriteString(fmt.Sprintf("%s%s: %s", pad, k, convertToShon(v, indent+1, "", compactArrays, preserveOrder, trailingCommas)))
            if i < len(keys)-1 || trailingCommas {
                sb.WriteString(",\n")
            } else {
                sb.WriteString("\n")
            }
        }

        closingPad := ""
        if indent > 0 {
            closingPad = strings.Repeat("    ", indent-1)
        }
        sb.WriteString(closingPad + "}")
        return sb.String()

    case []interface{}:
        if compactArrays && allSimpleObjects(val) {
            var sb strings.Builder
            sb.WriteString("[\n")
            for i, v := range val {
                item := v.(map[string]interface{})
                keys := make([]string, 0, len(item))
                for k := range item {
                    keys = append(keys, k)
                }
                sort.Strings(keys)

                fields := make([]string, 0, len(keys))
                for _, k := range keys {
                    fields = append(fields, fmt.Sprintf("%s: %s", k, convertToShon(item[k], 0, "", compactArrays, preserveOrder, trailingCommas)))
                }
                itemPad := strings.Repeat("    ", indent+1)
                sb.WriteString(fmt.Sprintf("%s{ %s }", itemPad, strings.Join(fields, ", ")))
                if i < len(val)-1 || trailingCommas {
                    sb.WriteString(",\n")
                } else {
                    sb.WriteString("\n")
                }
            }
            sb.WriteString(pad + "]")
            return sb.String()
        }

        var sb strings.Builder
        sb.WriteString("[\n")
        for i, v := range val {
            item := convertToShon(v, indent+1, "", compactArrays, preserveOrder, trailingCommas)
            sb.WriteString(fmt.Sprintf("%s%s", strings.Repeat("    ", indent+1), item))
            if i < len(val)-1 || trailingCommas {
                sb.WriteString(",\n")
            } else {
                sb.WriteString("\n")
            }
        }
        sb.WriteString(fmt.Sprintf("%s]", pad))
        return sb.String()

    case string:
        escaped := strings.ReplaceAll(val, "\"", "\\\"")
        return fmt.Sprintf("\"%s\"", escaped)

    case float64, bool, nil:
        return fmt.Sprintf("%v", val)

    default:
        return fmt.Sprintf("\"%v\"", val)
    }
}

func allSimpleObjects(array []interface{}) bool {
    for _, item := range array {
        obj, ok := item.(map[string]interface{})
        if !ok {
            return false
        }
        for _, v := range obj {
            kind := reflect.TypeOf(v).Kind()
            if kind == reflect.Map || kind == reflect.Slice {
                return false
            }
        }
    }
    return true
}
