package pkg_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sottey/shon/tooling/shon/pkg"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return tmp
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	return string(data)
}

func TestJsonToShon(t *testing.T) {
	input := `{
		"id": "001",
		"name": "Sean",
		"active": true,
		"balance": "1042.75",
		"created": "2025-03-22T14:45:00Z",
		"tags": ["dev", "golang"],
		"location": {
			"city": "Palm Springs",
			"state": "CA"
		}
	}`
	in := writeTempFile(t, "input.json", input)
	out := filepath.Join(t.TempDir(), "output.shon")

	err := pkg.JsonToShon(in, out, false)
	if err != nil {
		t.Fatalf("JsonToShon failed: %v", err)
	}

	result := readFile(t, out)
	if !strings.Contains(result, "@data") {
		t.Error("SHON output missing @data block")
	}
	if !strings.Contains(result, `$timestamp("2025-03-22T14:45:00Z")`) {
		t.Error("timestamp type not preserved")
	}
}

func TestShonToJson(t *testing.T) {
	input := `$schema: "output.shos"

@data {
	id: "001",
	name: "Sean",
	active: true,
	created: $timestamp("2025-03-22T14:45:00Z"),
	balance: $decimal("1042.75"),
	tags: ["dev", "golang"],
	location: {
		city: "Palm Springs",
		state: "CA"
	}
}`
	in := writeTempFile(t, "input.shon", input)
	out := filepath.Join(t.TempDir(), "output.json")

	err := pkg.ShonToJson(in, out)
	if err != nil {
		t.Fatalf("ShonToJson failed: %v", err)
	}

	result := readFile(t, out)
	if !strings.Contains(result, `"created": "2025-03-22T14:45:00Z"`) {
		t.Error("timestamp not correctly converted to JSON")
	}
}

func TestJsonShonRoundTrip(t *testing.T) {
	// Test json data
	jsonData := `{
		"id": "001",
		"name": "Sean",
		"active": true,
		"balance": "1042.75",
		"created": "2025-03-22T14:45:00Z",
		"tags": ["dev", "golang"],
		"location": {
			"city": "Palm Springs",
			"state": "CA"
		}
	}`

	// Create the temp file with the test data
	inputJsonFile := writeTempFile(t, "input.json", jsonData)

	// Build the output path
	outputShonFile := filepath.Join(t.TempDir(), "output.shon")

	// Convert the test data to SHON
	err := pkg.JsonToShon(inputJsonFile, outputShonFile, true)
	if err != nil {
		t.Fatalf("JsonToShon failed: %v", err)
	}

	// Read in the output (SHON data)
	shonData := readFile(t, outputShonFile)
	if !strings.Contains(shonData, "@data") {
		t.Error("SHON output missing @data block")
	}
	if !strings.Contains(shonData, `$timestamp("2025-03-22T14:45:00Z")`) {
		t.Error("timestamp type not preserved")
	}

	// Create a new output file for the json
	outputJsonFile := filepath.Join(t.TempDir(), "output.json")

	// Create the output file via conversion
	err = pkg.ShonToJson(outputShonFile, outputJsonFile)
	if err != nil {
		t.Fatalf("ShonToJson failed: %v", err)
	}

	// Normalize both JSON data items before comparison
	var originalObj, roundTripObj interface{}
	if err := json.Unmarshal([]byte(jsonData), &originalObj); err != nil {
		t.Fatalf("Failed to unmarshal input JSON: %v", err)
	}

	// Read in the final json output
	jsonDataResult := readFile(t, outputJsonFile)
	if err := json.Unmarshal([]byte(jsonDataResult), &roundTripObj); err != nil {
		t.Fatalf("Failed to unmarshal output JSON: %v", err)
	}

	// do the compare
	if !reflect.DeepEqual(originalObj, roundTripObj) {
		prettyInput, _ := json.MarshalIndent(originalObj, "", "  ")
		prettyOutput, _ := json.MarshalIndent(roundTripObj, "", "  ")
		t.Errorf("Round-trip data mismatch:\nInput:\n%s\n\nOutput:\n%s", prettyInput, prettyOutput)
	}
}

func TestCsvToShon(t *testing.T) {
	input := `name,address,title
Sean,1234 Main St,Engineer
Ellie,1234 Main St,CTO
Darcy,5678 2nd Ave,Engineer`
	in := writeTempFile(t, "input.csv", input)
	out := filepath.Join(t.TempDir(), "output.shon")

	err := pkg.CSVToShon(in, out)
	if err != nil {
		t.Fatalf("CSVToShon failed: %v", err)
	}

	result := readFile(t, out)
	if !strings.Contains(result, "@address") || !strings.Contains(result, "@title") {
		t.Error("expected ref blocks missing")
	}
	if !strings.Contains(result, "&address.") {
		t.Error("references not generated")
	}
}
