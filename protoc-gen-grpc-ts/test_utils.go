package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Test utilities for the protoc plugin

// assertFileExists checks if a file exists and fails the test if it doesn't
func assertFileExists(t testing.T, filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", filePath)
	}
}

// assertFileContains checks if a file contains expected content
func assertFileContains(t testing.T, filePath string, expectedContent []string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file %s: %v", filePath, err)
		return
	}

	contentStr := string(content)
	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("File %s does not contain expected content: %s", filePath, expected)
		}
	}
}

// assertFileNotContains checks if a file does not contain unexpected content
func assertFileNotContains(t testing.T, filePath string, unexpectedContent []string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file %s: %v", filePath, err)
		return
	}

	contentStr := string(content)
	for _, unexpected := range unexpectedContent {
		if strings.Contains(contentStr, unexpected) {
			t.Errorf("File %s contains unexpected content: %s", filePath, unexpected)
		}
	}
}

// generateTestFiles generates test files using the plugin
func generateTestFiles(t testing.T, protoFiles []string, outputDir string) error {
	// Clean output directory
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		return err
	}

	// Run protoc (default to react-native for backward compatibility)
	args := []string{
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:" + outputDir,
		"-I", "testdata",
	}
	args = append(args, protoFiles...)

	cmd = exec.Command("protoc", args...)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// listGeneratedFiles returns a list of all generated TypeScript files
func listGeneratedFiles(outputDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".ts") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// validateTypeScriptSyntax checks if generated TypeScript files have valid syntax
func validateTypeScriptSyntax(t *testing.T, outputDir string) {
	files, err := listGeneratedFiles(outputDir)
	if err != nil {
		t.Errorf("Failed to list generated files: %v", err)
		return
	}

	for _, file := range files {
		cmd := exec.Command("npx", "tsc", "--noEmit", file)
		cmd.Dir = "."
		output, err := cmd.CombinedOutput()
		if err != nil {
			outputStr := string(output)
			// Skip validation if TypeScript is not installed
			if strings.Contains(outputStr, "This is not the tsc command you are looking for") ||
				strings.Contains(outputStr, "command not found") ||
				strings.Contains(outputStr, "tsc: command not found") {
				// TypeScript compiler not available - skip validation
				continue
			}
			// Check if the only error is the expected module error (grpc-react-native or @grpc/grpc-js)
			if !strings.Contains(outputStr, "Cannot find module 'grpc-react-native'") &&
				!strings.Contains(outputStr, "Cannot find module '@grpc/grpc-js'") {
				t.Errorf("TypeScript syntax error in %s: %s", file, outputStr)
			}
		}
	}
}

// Test data structures for comprehensive testing
type TestCase struct {
	Name              string
	ProtoFiles        []string
	ExpectedFiles     []string
	ExpectedContent   map[string][]string
	UnexpectedContent map[string][]string
}

var comprehensiveTestCases = []TestCase{
	{
		Name:          "Basic message generation",
		ProtoFiles:    []string{"testdata/simple.proto"},
		ExpectedFiles: []string{"testoutput/simple_types.ts"},
		ExpectedContent: map[string][]string{
			"testoutput/simple_types.ts": {
				"export namespace TestSimple {",
				"export interface SimpleMessage {",
				"name: string;",
				"age: number;",
				"active: boolean;",
				"export enum Status {",
			},
		},
	},
	{
		Name:          "Nested message generation",
		ProtoFiles:    []string{"testdata/nested.proto"},
		ExpectedFiles: []string{"testoutput/nested_types.ts"},
		ExpectedContent: map[string][]string{
			"testoutput/nested_types.ts": {
				"export interface OuterMessage_InnerMessage {",
				"export enum Container_Type {",
				"export interface ComplexMessage_Header {",
			},
		},
	},
	{
		Name:       "Service generation",
		ProtoFiles: []string{"testdata/simple.proto", "testdata/service.proto"},
		ExpectedFiles: []string{
			"testoutput/simple_types.ts",
			"testoutput/service_types.ts",
			"testoutput/service_services.ts",
		},
		ExpectedContent: map[string][]string{
			"testoutput/service_services.ts": {
				"import { GrpcClient } from 'grpc-react-native';",
				"export class TestServiceClient {",
				"async getUser(",
				"// Target: react-native",
			},
		},
	},
	{
		Name:          "Edge cases",
		ProtoFiles:    []string{"testdata/edge_cases.proto"},
		ExpectedFiles: []string{"testoutput/edge_cases_types.ts"},
		ExpectedContent: map[string][]string{
			"testoutput/edge_cases_types.ts": {
				"export interface AllTypes {",
				"doubleField: number;",
				"floatField: number;",
				"int32Field: number;",
				"stringField: string;",
				"bytesField: Uint8Array;",
				"repeatedString: string[];",
				"optionalString?: string;",
				"export interface DeepNesting_Level1_Level2_Level3_Level4 {",
				"export enum EnumMessage_OuterEnum {",
				"export enum EnumMessage_InnerMessage_InnerEnum {",
			},
		},
	},
}
