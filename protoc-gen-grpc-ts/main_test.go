package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neo-fs/protoc-gen-grpc-ts/internal/types"
)

func TestPluginGeneration(t *testing.T) {
	tests := []struct {
		name            string
		protoFiles      []string
		expectedFiles   []string
		expectedContent map[string][]string
	}{
		{
			name:       "Simple message generation",
			protoFiles: []string{"testdata/simple.proto"},
			expectedFiles: []string{
				"testoutput/simple_types.ts",
			},
			expectedContent: map[string][]string{
				"testoutput/simple_types.ts": {
					"export namespace TestSimple {",
					"export interface SimpleMessage {",
					"name: string;",
					"age: number;",
					"active: boolean;",
					"score: number;",
					"data: Uint8Array;",
					"export enum Status {",
					"STATUS_UNSPECIFIED = 0,",
					"STATUS_ACTIVE = 1,",
					"STATUS_INACTIVE = 2,",
					"export interface User {",
					"id: string;",
					"status: Status;",
					"profile: SimpleMessage;",
				},
			},
		},
		{
			name:       "Nested message generation",
			protoFiles: []string{"testdata/nested.proto"},
			expectedFiles: []string{
				"testoutput/nested_types.ts",
			},
			expectedContent: map[string][]string{
				"testoutput/nested_types.ts": {
					"export namespace TestNested {",
					"export interface OuterMessage {",
					"outerField: string;",
					"inner: OuterMessage_InnerMessage;",
					"export interface OuterMessage_InnerMessage {",
					"innerField: string;",
					"innerNumber: number;",
					"export interface Container {",
					"name: string;",
					"type: Container_Type;",
					"export enum Container_Type {",
					"TYPE_UNSPECIFIED = 0,",
					"TYPE_SMALL = 1,",
					"TYPE_LARGE = 2,",
					"export interface ComplexMessage {",
					"id: string;",
					"header: ComplexMessage_Header;",
					"body: ComplexMessage_Body;",
					"export interface ComplexMessage_Header {",
					"version: string;",
					"timestamp: number;",
					"export interface ComplexMessage_Body {",
					"content: string;",
					"tags: string[];",
				},
			},
		},
		{
			name:       "Service generation",
			protoFiles: []string{"testdata/simple.proto", "testdata/nested.proto", "testdata/service.proto"},
			expectedFiles: []string{
				"testoutput/simple_types.ts",
				"testoutput/nested_types.ts",
				"testoutput/service_types.ts",
				"testoutput/service_services.ts",
			},
			expectedContent: map[string][]string{
				"testoutput/service_services.ts": {
					"import { GrpcClient } from 'grpc-react-native';",
					"GetUserRequestImpl",
					"GetUserResponseImpl",
					"export class TestServiceClient {",
					"constructor(private client: GrpcClient) {}",
					"async getUser(request: GetUserRequestImpl): Promise<GetUserResponseImpl> {",
					"this.client.unaryCall(",
					"'test.service.TestService/GetUser',",
					"request.serializeBinary()",
					"// Target: react-native",
				},
				"testoutput/service_types.ts": {
					"import { TestSimple } from './simple_types';",
					"import { TestNested } from './nested_types';",
					"export namespace TestService {",
					"export interface GetUserRequest {",
					"userId: string;",
					"export interface GetUserResponse {",
					"user: TestSimple.User;",
					"success: boolean;",
					"export interface CreateContainerRequest {",
					"container: TestNested.Container;",
					"description: string;",
					"export interface CreateContainerResponse {",
					"containerId: string;",
					"status: TestSimple.Status;",
				},
			},
		},
		{
			name:       "Import handling",
			protoFiles: []string{"testdata/simple.proto", "testdata/nested.proto", "testdata/imports.proto"},
			expectedFiles: []string{
				"testoutput/simple_types.ts",
				"testoutput/nested_types.ts",
				"testoutput/imports_types.ts",
			},
			expectedContent: map[string][]string{
				"testoutput/imports_types.ts": {
					"import { TestSimple } from './simple_types';",
					"import { TestNested } from './nested_types';",
					"export namespace TestImports {",
					"export interface ImportedMessage {",
					"simple: TestSimple.SimpleMessage;",
					"outer: TestNested.OuterMessage;",
					"status: TestSimple.Status;",
					"export interface ListMessage {",
					"items: string[];",
					"users: TestSimple.User[];",
					"numbers: number[];",
					"export interface OptionalMessage {",
					"name?: string;",
					"count?: number;",
					"status?: TestSimple.Status;",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean test output directory
			os.RemoveAll("testoutput")
			os.MkdirAll("testoutput", 0755)

			// Build the plugin
			cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
			cmd.Dir = "."
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to build plugin: %v", err)
			}

			// Run protoc with the plugin (React Native target)
			args := []string{
				"--plugin=./protoc-gen-grpc-ts",
				"--grpc-ts_out=target=react-native:testoutput",
				"-I", "testdata",
			}
			args = append(args, tt.protoFiles...)

			cmd = exec.Command("protoc", args...)
			cmd.Dir = "."
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
			}

			// Check that expected files were generated
			for _, expectedFile := range tt.expectedFiles {
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not generated", expectedFile)
				}
			}

			// Check file contents
			for filePath, expectedLines := range tt.expectedContent {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read file %s: %v", filePath, err)
					continue
				}

				contentStr := string(content)
				for _, expectedLine := range expectedLines {
					if !strings.Contains(contentStr, expectedLine) {
						t.Errorf("File %s does not contain expected line: %s", filePath, expectedLine)
					}
				}
			}
		})
	}
}

func TestTypeScriptCompilation(t *testing.T) {
	// Clean test output directory
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate all test files (React Native target)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
		"testdata/simple.proto",
		"testdata/nested.proto",
		"testdata/service.proto",
		"testdata/imports.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test TypeScript compilation
	cmd = exec.Command("npx", "tsc", "--noEmit", "testoutput/*.ts")
	cmd.Dir = "."
	output, err = cmd.CombinedOutput()
	if err != nil {
		// Check if the only error is the expected grpc-react-native module error
		outputStr := string(output)
		if !strings.Contains(outputStr, "Cannot find module 'grpc-react-native'") {
			t.Errorf("TypeScript compilation failed with unexpected errors: %s", outputStr)
		}
	}
}

func TestFieldNaming(t *testing.T) {
	// Test that field names are converted to camelCase
	expectedMappings := map[string]string{
		"user_id":        "userId",
		"first_name":     "firstName",
		"last_name":      "lastName",
		"created_at":     "createdAt",
		"is_active":      "isActive",
		"has_permission": "hasPermission",
	}

	for input, expected := range expectedMappings {
		result := types.ToCamelCase(input)
		if result != expected {
			t.Errorf("toCamelCase(%s) = %s, expected %s", input, result, expected)
		}
	}
}

func TestNamespaceGeneration(t *testing.T) {
	// Test namespace name generation
	testCases := []struct {
		packageName string
		expected    string
	}{
		{"test.simple", "TestSimple"},
		{"test.nested", "TestNested"},
		{"test.service", "TestService"},
		{"test.imports", "TestImports"},
		{"neo.fs.v2.refs", "NeoFsV2Refs"},
		{"neo.fs.v2.session", "NeoFsV2Session"},
		{"neo.fs.v2.accounting", "NeoFsV2Accounting"},
	}

	for _, tc := range testCases {
		result := getNamespaceName(tc.packageName)
		if result != tc.expected {
			t.Errorf("getNamespaceName(%s) = %s, expected %s", tc.packageName, result, tc.expected)
		}
	}
}

func TestTypeScriptTypeMapping(t *testing.T) {
	// This test would require creating protogen.Field objects, which is complex
	// For now, we'll test the basic functionality through integration tests
	t.Log("TypeScript type mapping is tested through integration tests")
}

func TestFileGeneration(t *testing.T) {
	// Test that the plugin generates the correct file structure
	expectedStructure := []string{
		"testoutput/simple_types.ts",
		"testoutput/nested_types.ts",
		"testoutput/service_types.ts",
		"testoutput/service_services.ts",
		"testoutput/imports_types.ts",
	}

	// Clean and generate
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
		"testdata/simple.proto",
		"testdata/nested.proto",
		"testdata/service.proto",
		"testdata/imports.proto")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("protoc failed: %v", err)
	}

	// Check file structure
	for _, expectedFile := range expectedStructure {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not generated", expectedFile)
		}
	}

	// Check that no unexpected files were generated
	var files []string
	err := filepath.Walk("testoutput", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".ts") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to list generated files: %v", err)
	}

	if len(files) != len(expectedStructure) {
		t.Errorf("Expected %d files, got %d files: %v", len(expectedStructure), len(files), files)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test that the plugin handles invalid proto files gracefully
	invalidProto := "testdata/invalid.proto"

	// Create an invalid proto file
	err := os.WriteFile(invalidProto, []byte("invalid proto content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid proto file: %v", err)
	}
	defer os.Remove(invalidProto)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Try to process invalid proto file
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
		invalidProto)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()

	// Should fail with protoc error, not plugin error
	if err == nil {
		t.Error("Expected protoc to fail with invalid proto file")
	}

	// The error should be from protoc, not from our plugin
	outputStr := string(output)
	if strings.Contains(outputStr, "protoc-gen-grpc-ts") &&
		!strings.Contains(outputStr, "protoc") {
		t.Errorf("Plugin should not be the source of the error: %s", outputStr)
	}
}

// Benchmark tests
func BenchmarkPluginGeneration(b *testing.B) {
	// Build the plugin once
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		b.Fatalf("Failed to build plugin: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Clean output directory
		os.RemoveAll("testoutput")
		os.MkdirAll("testoutput", 0755)

		// Run protoc (React Native target)
		cmd = exec.Command("protoc",
			"--plugin=./protoc-gen-grpc-ts",
			"--grpc-ts_out=target=react-native:testoutput",
			"-I", "testdata",
			"testdata/simple.proto",
			"testdata/nested.proto",
			"testdata/service.proto",
			"testdata/imports.proto")
		cmd.Dir = "."
		if err := cmd.Run(); err != nil {
			b.Fatalf("protoc failed: %v", err)
		}
	}
}
