package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStreamingMethods(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-ts", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate streaming service files (React Native target)
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-ts",
		"--grpc-ts_out=target=react-native:testoutput",
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

	// Test server streaming method (React Native pattern)
	if !strings.Contains(servicesContent, "async* getData(request: GetDataRequestImpl): AsyncGenerator<GetDataResponseImpl>") {
		t.Error("Server streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.serverStreamCall") {
		t.Error("Server streaming call not found")
	}

	// Test client streaming method (React Native pattern)
	if !strings.Contains(servicesContent, "async uploadData(requests: AsyncIterable<UploadDataRequestImpl>): Promise<UploadDataResponseImpl>") {
		t.Error("Client streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.clientStreamCall") {
		t.Error("Client streaming call not found")
	}

	// Test bidirectional streaming method (React Native pattern)
	if !strings.Contains(servicesContent, "async* chat(requests: AsyncIterable<ChatMessageImpl>): AsyncGenerator<ChatMessageImpl>") {
		t.Error("Bidirectional streaming method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.bidirectionalStreamCall") {
		t.Error("Bidirectional streaming call not found")
	}

	// Test unary method (React Native pattern)
	if !strings.Contains(servicesContent, "async getStatus(request: GetStatusRequestImpl): Promise<GetStatusResponseImpl>") {
		t.Error("Unary method not found or incorrect signature")
	}
	if !strings.Contains(servicesContent, "this.client.unaryCall") {
		t.Error("Unary call not found")
	}

	// Test service class structure
	if !strings.Contains(servicesContent, "export class StreamingServiceClient") {
		t.Error("Service client class not found")
	}
	if !strings.Contains(servicesContent, "constructor(private client: GrpcClient)") {
		t.Error("Constructor not found")
	}
}

func TestStreamingTypesGeneration(t *testing.T) {
	// Test that streaming types are generated correctly
	typesFile := "testoutput/streaming_service_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test namespace
	if !strings.Contains(typesContent, "export namespace TestStreaming") {
		t.Error("Namespace not found")
	}

	// Test message interfaces
	expectedInterfaces := []string{
		"export interface GetDataRequest",
		"export interface GetDataResponse",
		"export interface UploadDataRequest",
		"export interface UploadDataResponse",
		"export interface ChatMessage",
		"export interface GetStatusRequest",
		"export interface GetStatusResponse",
	}

	for _, iface := range expectedInterfaces {
		if !strings.Contains(typesContent, iface) {
			t.Errorf("Interface %s not found", iface)
		}
	}

	// Test implementation classes
	expectedClasses := []string{
		"export class GetDataRequestImpl",
		"export class GetDataResponseImpl",
		"export class UploadDataRequestImpl",
		"export class UploadDataResponseImpl",
		"export class ChatMessageImpl",
		"export class GetStatusRequestImpl",
		"export class GetStatusResponseImpl",
	}

	for _, class := range expectedClasses {
		if !strings.Contains(typesContent, class) {
			t.Errorf("Implementation class %s not found", class)
		}
	}

	// Test field types
	expectedFields := map[string][]string{
		"GetDataRequest":     {"Query: string", "Limit: number"},
		"GetDataResponse":    {"Data: string", "HasMore: boolean"},
		"UploadDataRequest":  {"Filename: string", "Chunk: Uint8Array", "IsLast: boolean"},
		"UploadDataResponse": {"FileId: string", "TotalSize: bigint", "Success: boolean"},
		"ChatMessage":        {"UserId: string", "Message: string", "Timestamp: bigint"},
		"GetStatusRequest":   {"Service: string"},
		"GetStatusResponse":  {"Status: string", "Uptime: number"},
	}

	for messageName, fields := range expectedFields {
		for _, field := range fields {
			if !strings.Contains(typesContent, field) {
				t.Errorf("Field %s not found in %s", field, messageName)
			}
		}
	}

	// Test serialization methods
	if !strings.Contains(typesContent, "serializeBinary(): Uint8Array") {
		t.Error("serializeBinary method not found")
	}
	if !strings.Contains(typesContent, "static deserializeBinary(data: Uint8Array)") {
		t.Error("deserializeBinary method not found")
	}
}
