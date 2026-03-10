package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestTypeMapping(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-react-native", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate files with various field types
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-react-native",
		"--grpc-react-native_out=testoutput",
		"-I", "testdata",
		"testdata/basic_types.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test type mapping
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test specific type mappings
	typeMappingTests := map[string]string{
		"int32":  "number",
		"int64":  "bigint",
		"uint32": "number",
		"uint64": "bigint",
		"float":  "number",
		"double": "number",
		"bool":   "boolean",
		"string": "string",
		"bytes":  "Uint8Array",
	}

	for protoType, expectedTS := range typeMappingTests {
		// Look for field declarations with the expected TypeScript type
		if !strings.Contains(typesContent, ": "+expectedTS) {
			t.Errorf("Type mapping for %s -> %s not found", protoType, expectedTS)
		}
	}

	// Test repeated fields
	if !strings.Contains(typesContent, "string[]") {
		t.Error("Repeated string field not found")
	}
	if !strings.Contains(typesContent, "number[]") {
		t.Error("Repeated number field not found")
	}

	// Test map fields
	if !strings.Contains(typesContent, "Map<") {
		t.Error("Map field type not found")
	}
}

func TestDefaultValues(t *testing.T) {
	// Test that default values are generated correctly
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test default value assignments in constructor
	expectedDefaults := map[string]string{
		"string":     `""`,
		"number":     "0",
		"bigint":     "0n",
		"boolean":    "false",
		"Uint8Array": "new Uint8Array(0)",
		"[]":         "[]",
		"Map":        "new Map()",
	}

	for fieldType, expectedDefault := range expectedDefaults {
		if !strings.Contains(typesContent, expectedDefault) {
			t.Errorf("Default value for %s not found: %s", fieldType, expectedDefault)
		}
	}

	// Test constructor parameter handling
	if !strings.Contains(typesContent, "constructor(data?: Partial<") {
		t.Error("Constructor with optional parameter not found")
	}
	if !strings.Contains(typesContent, "data?.") {
		t.Error("Optional parameter access not found")
	}
	if !strings.Contains(typesContent, "?? ") {
		t.Error("Nullish coalescing operator not found")
	}
}

func TestRepeatedFields(t *testing.T) {
	// Test repeated field handling
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test repeated field type declarations
	if !strings.Contains(typesContent, "Tags: string[]") {
		t.Error("Repeated string field type not found")
	}
	if !strings.Contains(typesContent, "Numbers: number[]") {
		t.Error("Repeated number field type not found")
	}

	// Test repeated field serialization
	if !strings.Contains(typesContent, "for (const item of this.") {
		t.Error("Repeated field serialization loop not found")
	}

	// Test repeated field deserialization
	if !strings.Contains(typesContent, ".push(") {
		t.Error("Repeated field deserialization push not found")
	}

	// Test default value for repeated fields
	if !strings.Contains(typesContent, "?? []") {
		t.Error("Default empty array for repeated fields not found")
	}
}

func TestMapFields(t *testing.T) {
	// Test map field handling
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test map field type declarations
	if !strings.Contains(typesContent, "Map<") {
		t.Error("Map field type declaration not found")
	}

	// Test map field serialization
	if !strings.Contains(typesContent, "for (const [key, value] of this.") {
		t.Error("Map field serialization loop not found")
	}

	// Test map field deserialization
	if !strings.Contains(typesContent, ".set(key, value)") {
		t.Error("Map field deserialization set not found")
	}

	// Test default value for map fields
	if !strings.Contains(typesContent, "new Map()") {
		t.Error("Default empty map for map fields not found")
	}
}

func TestNestedMessages(t *testing.T) {
	// Test nested message handling
	cmd := exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-react-native",
		"--grpc-react-native_out=testoutput",
		"-I", "testdata",
		"testdata/nested.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	typesFile := "testoutput/nested_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test nested message interfaces
	if !strings.Contains(typesContent, "export interface OuterMessage") {
		t.Error("Outer message interface not found")
	}
	if !strings.Contains(typesContent, "export interface OuterMessage_InnerMessage") {
		t.Error("Nested message interface not found")
	}

	// Test nested message classes
	if !strings.Contains(typesContent, "export class OuterMessageImpl") {
		t.Error("Outer message class not found")
	}
	if !strings.Contains(typesContent, "export class OuterMessage_InnerMessageImpl") {
		t.Error("Nested message class not found")
	}

	// Test nested field references
	if !strings.Contains(typesContent, "Inner: OuterMessage_InnerMessage") {
		t.Error("Nested field reference not found")
	}
}

func TestEnumGeneration(t *testing.T) {
	// Test enum generation
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test enum declaration
	if !strings.Contains(typesContent, "export enum Status") {
		t.Error("Enum declaration not found")
	}

	// Test enum values
	expectedEnumValues := []string{
		"STATUS_UNSPECIFIED = 0,",
		"STATUS_ACTIVE = 1,",
		"STATUS_INACTIVE = 2,",
	}

	for _, enumValue := range expectedEnumValues {
		if !strings.Contains(typesContent, enumValue) {
			t.Errorf("Enum value not found: %s", enumValue)
		}
	}

	// Test enum field usage
	if !strings.Contains(typesContent, "Status: Status") {
		t.Error("Enum field usage not found")
	}
}
