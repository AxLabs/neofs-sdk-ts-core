package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestNodeJSServiceGeneration tests Node.js service generation
func TestNodeJSServiceGeneration(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate service files (Node.js target)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=nodejs:testoutput",
		"-I", "testdata",
		"testdata/simple.proto",
		"testdata/service.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test that files were generated
	expectedFiles := []string{
		"testoutput/simple_types.ts",
		"testoutput/service_types.ts",
		"testoutput/service_services.ts",
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not generated", file)
		}
	}

	// Test service file content
	servicesFile := "testoutput/service_services.ts"
	content, err := os.ReadFile(servicesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", servicesFile, err)
	}
	servicesContent := string(content)

	// Test Node.js specific imports
	if !strings.Contains(servicesContent, "import * as grpc from '@grpc/grpc-js';") {
		t.Error("Node.js import not found")
	}
	if !strings.Contains(servicesContent, "// Target: nodejs") {
		t.Error("Node.js target comment not found")
	}

	// Test service class structure
	if !strings.Contains(servicesContent, "export class TestServiceClient {") {
		t.Error("Service client class not found")
	}
	if !strings.Contains(servicesContent, "private client: grpc.Client;") {
		t.Error("Client field not found")
	}
	if !strings.Contains(servicesContent, "constructor(address: string, credentials: grpc.ChannelCredentials, options?: grpc.ChannelOptions)") {
		t.Error("Constructor not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "close(): void {") {
		t.Error("Close method not found")
	}

	// Test unary method
	if !strings.Contains(servicesContent, "getUser(request: GetUserRequestImpl, metadata?: grpc.Metadata, options?: grpc.CallOptions): Promise<GetUserResponseImpl>") {
		t.Error("Unary method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.makeUnaryRequest") {
		t.Error("makeUnaryRequest not found")
	}
	if !strings.Contains(servicesContent, "Buffer.from(arg.serializeBinary())") {
		t.Error("Buffer serialization not found")
	}
	if !strings.Contains(servicesContent, ".deserializeBinary(new Uint8Array(buf))") {
		t.Error("Buffer deserialization not found")
	}
}

// TestNodeJSStreamingMethods tests Node.js streaming method generation
func TestNodeJSStreamingMethods(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate streaming service files (Node.js target)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=nodejs:testoutput",
		"-I", "testdata",
		"testdata/streaming_service.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test that files were generated
	expectedFiles := []string{
		"testoutput/streaming_service_types.ts",
		"testoutput/streaming_service_services.ts",
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not generated", file)
		}
	}

	// Test streaming service content
	servicesFile := "testoutput/streaming_service_services.ts"
	content, err := os.ReadFile(servicesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", servicesFile, err)
	}
	servicesContent := string(content)

	// Test Node.js imports
	if !strings.Contains(servicesContent, "import * as grpc from '@grpc/grpc-js';") {
		t.Error("Node.js import not found")
	}
	if !strings.Contains(servicesContent, "// Target: nodejs") {
		t.Error("Node.js target comment not found")
	}

	// Test server streaming method
	if !strings.Contains(servicesContent, "getData(request: GetDataRequestImpl, metadata?: grpc.Metadata, options?: grpc.CallOptions): grpc.ClientReadableStream<GetDataResponseImpl>") {
		t.Error("Server streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.makeServerStreamRequest") {
		t.Error("makeServerStreamRequest not found")
	}

	// Test client streaming method
	if !strings.Contains(servicesContent, "uploadData(metadata?: grpc.Metadata, options?: grpc.CallOptions, callback?: (error: grpc.ServiceError | null, response?: UploadDataResponseImpl) => void): grpc.ClientWritableStream<UploadDataRequestImpl>") {
		t.Error("Client streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.makeClientStreamRequest") {
		t.Error("makeClientStreamRequest not found")
	}

	// Test bidirectional streaming method
	if !strings.Contains(servicesContent, "chat(metadata?: grpc.Metadata, options?: grpc.CallOptions): grpc.ClientDuplexStream<ChatMessageImpl, ChatMessageImpl>") {
		t.Error("Bidirectional streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.makeBidiStreamRequest") {
		t.Error("makeBidiStreamRequest not found")
	}

	// Test unary method
	if !strings.Contains(servicesContent, "getStatus(request: GetStatusRequestImpl, metadata?: grpc.Metadata, options?: grpc.CallOptions): Promise<GetStatusResponseImpl>") {
		t.Error("Unary method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.makeUnaryRequest") {
		t.Error("makeUnaryRequest not found")
	}

	// Test service class structure
	if !strings.Contains(servicesContent, "export class StreamingServiceClient {") {
		t.Error("Service client class not found")
	}
	if !strings.Contains(servicesContent, "private client: grpc.Client;") {
		t.Error("Client field not found")
	}
}

// TestNodeJSTargetFlag tests that the target flag works correctly
func TestNodeJSTargetFlag(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Test Node.js target
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=nodejs:testoutput",
		"-I", "testdata",
		"testdata/service.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Check Node.js specific content
	content, err := os.ReadFile("testoutput/service_services.ts")
	if err != nil {
		t.Fatalf("Failed to read service file: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "import * as grpc from '@grpc/grpc-js';") {
		t.Error("Node.js target should import @grpc/grpc-js")
	}
	if !strings.Contains(contentStr, "// Target: nodejs") {
		t.Error("Node.js target comment should be present")
	}
	if strings.Contains(contentStr, "import { GrpcClient } from 'grpc-react-native';") {
		t.Error("Node.js target should not import grpc-react-native")
	}

	// Clean and test React Native target
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
		"-I", "testdata",
		"testdata/service.proto")
	cmd.Dir = "."
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Check React Native specific content
	content, err = os.ReadFile("testoutput/service_services.ts")
	if err != nil {
		t.Fatalf("Failed to read service file: %v", err)
	}
	contentStr = string(content)

	if !strings.Contains(contentStr, "import { GrpcClient } from 'grpc-react-native';") {
		t.Error("React Native target should import grpc-react-native")
	}
	if !strings.Contains(contentStr, "// Target: react-native") {
		t.Error("React Native target comment should be present")
	}
	if strings.Contains(contentStr, "import * as grpc from '@grpc/grpc-js';") {
		t.Error("React Native target should not import @grpc/grpc-js")
	}
}

// TestInvalidTarget tests that invalid target flags default to react-native
func TestInvalidTarget(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Test invalid target (should default to react-native)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=invalid:testoutput",
		"-I", "testdata",
		"testdata/service.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Should default to react-native generation
	content, err := os.ReadFile("testoutput/service_services.ts")
	if err != nil {
		t.Fatalf("Failed to read service file: %v", err)
	}
	contentStr := string(content)

	// Should use react-native (default behavior)
	if !strings.Contains(contentStr, "import { GrpcClient } from 'grpc-react-native';") {
		t.Error("Invalid target should default to react-native")
	}
}
