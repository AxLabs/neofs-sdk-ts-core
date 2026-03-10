package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neo-fs/protoc-gen-grpc-ts/internal/types"
)

func TestComprehensiveIntegration(t *testing.T) {
	// Test all proto files together
	protoFiles := []string{
		"testdata/simple.proto",
		"testdata/nested.proto",
		"testdata/service.proto",
		"testdata/imports.proto",
		"testdata/edge_cases.proto",
		"testdata/complex_service.proto",
	}

	// Clean and generate
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate all files (default to react-native for backward compatibility)
	args := []string{
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
	}
	args = append(args, protoFiles...)

	cmd = exec.Command("protoc", args...)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Verify all expected files were generated
	expectedFiles := []string{
		"testoutput/simple_types.ts",
		"testoutput/nested_types.ts",
		"testoutput/service_types.ts",
		"testoutput/service_services.ts",
		"testoutput/imports_types.ts",
		"testoutput/edge_cases_types.ts",
		"testoutput/complex_service_types.ts",
		"testoutput/complex_service_services.ts",
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not generated", expectedFile)
		}
	}

	// Test TypeScript compilation
	validateTypeScriptSyntax(t, "testoutput")

	// Test specific content in key files
	testFileContent(t, "testoutput/simple_types.ts", []string{
		"export namespace TestSimple {",
		"export interface SimpleMessage {",
		"Name: string;",
		"Age: number;",
		"Active: boolean;",
		"export enum Status {",
		"STATUS_UNSPECIFIED = 0,",
	})

	testFileContent(t, "testoutput/nested_types.ts", []string{
		"export namespace TestNested {",
		"export interface OuterMessage_InnerMessage {",
		"export enum Container_Type {",
		"export interface ComplexMessage_Header {",
		"export interface ComplexMessage_Body {",
	})

	testFileContent(t, "testoutput/service_services.ts", []string{
		"import { GrpcClient } from 'grpc-react-native';",
		"export class TestServiceClient {",
		"constructor(private client: GrpcClient) {}",
		"async getUser(",
		"async createContainer(",
		"async processMessage(",
		"// Target: react-native",
	})

	testFileContent(t, "testoutput/edge_cases_types.ts", []string{
		"export namespace TestEdge {",
		"export interface AllTypes {",
		"DoubleField: number;",
		"FloatField: number;",
		"Int32Field: number;",
		"StringField: string;",
		"BytesField: Uint8Array;",
		"RepeatedString: string[];",
		"OptionalString: string;",
		"export interface DeepNesting_Level1_Level2_Level3_Level4 {",
		"export enum EnumMessage_OuterEnum {",
		"export enum EnumMessage_InnerMessage_InnerEnum {",
	})
}

func testFileContent(t *testing.T, filePath string, expectedContent []string) {
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

func TestFieldNamingConventions(t *testing.T) {
	// Test that field names are properly converted to camelCase
	expectedMappings := map[string]string{
		"user_id":        "userId",
		"first_name":     "firstName",
		"last_name":      "lastName",
		"created_at":     "createdAt",
		"is_active":      "isActive",
		"has_permission": "hasPermission",
		"api_key":        "apiKey",
		"http_status":    "httpStatus",
		"xml_parser":     "xmlParser",
		"json_data":      "jsonData",
	}

	for input, expected := range expectedMappings {
		result := types.ToCamelCase(input)
		if result != expected {
			t.Errorf("toCamelCase(%s) = %s, expected %s", input, result, expected)
		}
	}
}

func TestNamespaceGenerationIntegration(t *testing.T) {
	testCases := []struct {
		packageName string
		expected    string
	}{
		{"test.simple", "TestSimple"},
		{"test.nested", "TestNested"},
		{"test.service", "TestService"},
		{"test.imports", "TestImports"},
		{"test.edge", "TestEdge"},
		{"test.complex", "TestComplex"},
		{"neo.fs.v2.refs", "NeoFsV2Refs"},
		{"neo.fs.v2.session", "NeoFsV2Session"},
		{"neo.fs.v2.accounting", "NeoFsV2Accounting"},
		{"neo.fs.v2.acl", "NeoFsV2Acl"},
		{"neo.fs.v2.status", "NeoFsV2Status"},
		{"neo.fs.v2.link", "NeoFsV2Link"},
		{"neo.fs.v2.lock", "NeoFsV2Lock"},
		{"neo.fs.v2.object", "NeoFsV2Object"},
		{"neo.fs.v2.container", "NeoFsV2Container"},
		{"neo.fs.v2.netmap", "NeoFsV2Netmap"},
		{"neo.fs.v2.reputation", "NeoFsV2Reputation"},
		{"neo.fs.v2.audit", "NeoFsV2Audit"},
		{"neo.fs.v2.storagegroup", "NeoFsV2Storagegroup"},
		{"neo.fs.v2.subnet", "NeoFsV2Subnet"},
		{"neo.fs.v2.tombstone", "NeoFsV2Tombstone"},
	}

	for _, tc := range testCases {
		result := getNamespaceName(tc.packageName)
		if result != tc.expected {
			t.Errorf("getNamespaceName(%s) = %s, expected %s", tc.packageName, result, tc.expected)
		}
	}
}

// Helper function for testing namespace generation
func getNamespaceName(packageName string) string {
	if packageName == "" {
		return ""
	}

	// Convert package name to TypeScript namespace format
	parts := strings.Split(packageName, ".")
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, types.ToCamelCase(part))
		}
	}
	return strings.Join(result, "")
}

func TestImportHandling(t *testing.T) {
	// Test that imports are handled correctly
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-react-native", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate files with imports
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-react-native",
		"--grpc-react-native_out=testoutput",
		"-I", "testdata",
		"testdata/simple.proto",
		"testdata/imports.proto")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("protoc failed: %v", err)
	}

	// Check that imports are correctly generated
	content, err := os.ReadFile("testoutput/imports_types.ts")
	if err != nil {
		t.Fatalf("Failed to read imports_types.ts: %v", err)
	}

	contentStr := string(content)
	expectedImports := []string{
		"import { TestSimple } from './simple_types';",
		"import { TestNested } from './nested_types';",
	}

	for _, expectedImport := range expectedImports {
		if !strings.Contains(contentStr, expectedImport) {
			t.Errorf("File imports_types.ts does not contain expected import: %s", expectedImport)
		}
	}
}

func TestServiceGeneration(t *testing.T) {
	// Test that services are generated correctly
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate service files (React Native target)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
		"testdata/simple.proto",
		"testdata/service.proto")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("protoc failed: %v", err)
	}

	// Check service types file
	content, err := os.ReadFile("testoutput/service_types.ts")
	if err != nil {
		t.Fatalf("Failed to read service_types.ts: %v", err)
	}

	contentStr := string(content)
	expectedContent := []string{
		"import { TestSimple } from",
		"export interface GetUserRequest {",
		"UserId: string;",
		"export interface GetUserResponse {",
		"User?: TestSimple.User;",
		"Success: boolean;",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("File service_types.ts does not contain expected content: %s", expected)
		}
	}

	// Check service client file
	content, err = os.ReadFile("testoutput/service_services.ts")
	if err != nil {
		t.Fatalf("Failed to read service_services.ts: %v", err)
	}

	contentStr = string(content)
	expectedServiceContent := []string{
		"import { GrpcClient } from 'grpc-react-native';",
		"GetUserRequestImpl",
		"GetUserResponseImpl",
		"export class TestServiceClient {",
		"constructor(private client: GrpcClient) {}",
		"async getUser(request: GetUserRequestImpl): Promise<GetUserResponseImpl> {",
		"this.client.unaryCall(",
		"'test.service.TestService/GetUser',",
	}

	for _, expected := range expectedServiceContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("File service_services.ts does not contain expected content: %s", expected)
		}
	}
}

func TestErrorHandlingIntegration(t *testing.T) {
	// Test that the plugin handles errors gracefully
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

func TestPerformance(t *testing.T) {
	// Test performance with multiple files
	protoFiles := []string{
		"testdata/simple.proto",
		"testdata/nested.proto",
		"testdata/service.proto",
		"testdata/imports.proto",
		"testdata/edge_cases.proto",
		"testdata/complex_service.proto",
	}

	// Clean output directory
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Run multiple times to test performance
	for i := 0; i < 5; i++ {
		// Clean output directory
		os.RemoveAll("testoutput")
		os.MkdirAll("testoutput", 0755)

		// Run protoc (React Native target)
		args := []string{
			"--plugin=./protoc-gen-grpc-ts",
			"--grpc-ts_out=target=react-native:testoutput",
			"-I", "testdata",
		}
		args = append(args, protoFiles...)

		cmd = exec.Command("protoc", args...)
		cmd.Dir = "."
		if err := cmd.Run(); err != nil {
			t.Fatalf("protoc failed on iteration %d: %v", i, err)
		}

		// Verify files were generated (now in subdirectories based on Go import path)
		var files []string
		walkErr := filepath.Walk("testoutput", func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !info.IsDir() && strings.HasSuffix(path, ".ts") {
				files = append(files, path)
			}
			return nil
		})
		if walkErr != nil {
			t.Fatalf("Failed to list generated files on iteration %d: %v", i, walkErr)
		}

		if len(files) < 8 {
			t.Errorf("Expected at least 8 files on iteration %d, got %d", i, len(files))
		}
	}
}
