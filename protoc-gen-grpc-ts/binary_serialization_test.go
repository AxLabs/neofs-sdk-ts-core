package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestBinarySerialization(t *testing.T) {
	// Clean and setup
	os.RemoveAll("testoutput")
	os.MkdirAll("testoutput", 0755)

	// Build the plugin
	cmd := exec.Command("go", "build", "-o", "protoc-gen-grpc-react-native", "main.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build plugin: %v", err)
	}

	// Generate simple service files for testing
	cmd = exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-react-native",
		"--grpc-react-native_out=testoutput",
		"-I", "testdata",
		"testdata/simple.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test BinaryWriter class generation
	typesFile := "testoutput/simple_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test BinaryWriter class
	if !strings.Contains(typesContent, "class BinaryWriter") {
		t.Error("BinaryWriter class not found")
	}

	// Test BinaryWriter methods
	expectedWriterMethods := []string{
		"writeDouble(fieldNumber: number, value: number): void",
		"writeFloat(fieldNumber: number, value: number): void",
		"writeInt32(fieldNumber: number, value: number): void",
		"writeInt64(fieldNumber: number, value: bigint): void",
		"writeUint32(fieldNumber: number, value: number): void",
		"writeUint64(fieldNumber: number, value: bigint): void",
		"writeBool(fieldNumber: number, value: boolean): void",
		"writeString(fieldNumber: number, value: string): void",
		"writeBytes(fieldNumber: number, value: Uint8Array): void",
		"writeMessage(fieldNumber: number, value: any): void",
		"writeEnum(fieldNumber: number, value: number): void",
		"getResultBuffer(): Uint8Array",
	}

	for _, method := range expectedWriterMethods {
		if !strings.Contains(typesContent, method) {
			t.Errorf("BinaryWriter method not found: %s", method)
		}
	}

	// Test BinaryReader class
	if !strings.Contains(typesContent, "class BinaryReader") {
		t.Error("BinaryReader class not found")
	}

	// Test BinaryReader methods
	expectedReaderMethods := []string{
		"readDouble(): number",
		"readFloat(): number",
		"readInt32(): number",
		"readInt64(): bigint",
		"readUint32(): number",
		"readUint64(): bigint",
		"readBool(): boolean",
		"readString(): string",
		"readBytes(): Uint8Array",
		"readMessage<T>(deserialize: (data: Uint8Array) => T): T",
		"readEnum(): number",
		"readVarint(): number",
		"readVarintBigInt(): bigint",
		"skipField(wireType: number): void",
	}

	for _, method := range expectedReaderMethods {
		if !strings.Contains(typesContent, method) {
			t.Errorf("BinaryReader method not found: %s", method)
		}
	}

	// Test TextEncoder/TextDecoder declarations
	if !strings.Contains(typesContent, "declare const TextEncoder") {
		t.Error("TextEncoder declaration not found")
	}
	if !strings.Contains(typesContent, "declare const TextDecoder") {
		t.Error("TextDecoder declaration not found")
	}

	// Test serializeBinary method in message classes
	if !strings.Contains(typesContent, "serializeBinary(): Uint8Array {") {
		t.Error("serializeBinary method not found in message classes")
	}
	if !strings.Contains(typesContent, "const writer = new BinaryWriter();") {
		t.Error("BinaryWriter instantiation not found in serializeBinary")
	}
	if !strings.Contains(typesContent, "return writer.getResultBuffer();") {
		t.Error("getResultBuffer call not found in serializeBinary")
	}

	// Test deserializeBinary method in message classes
	if !strings.Contains(typesContent, "static deserializeBinary(data: Uint8Array)") {
		t.Error("deserializeBinary method not found in message classes")
	}
	if !strings.Contains(typesContent, "const reader = new BinaryReader(data);") {
		t.Error("BinaryReader instantiation not found in deserializeBinary")
	}
	if !strings.Contains(typesContent, "while (reader.position < reader.buffer.length)") {
		t.Error("Wire format parsing loop not found in deserializeBinary")
	}
}

func TestBinarySerializationFieldHandling(t *testing.T) {
	// Test that different field types generate correct serialization code
	typesFile := "testoutput/simple_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test field serialization for different types
	expectedSerialization := map[string]string{
		"string":     "writer.writeString(",
		"number":     "writer.writeInt32(",
		"boolean":    "writer.writeBool(",
		"Uint8Array": "writer.writeBytes(",
	}

	for fieldType, expectedCall := range expectedSerialization {
		if !strings.Contains(typesContent, expectedCall) {
			t.Errorf("Serialization call for %s not found: %s", fieldType, expectedCall)
		}
	}

	// Test field deserialization
	expectedDeserialization := map[string]string{
		"string":     "reader.readString()",
		"number":     "reader.readInt32()",
		"boolean":    "reader.readBool()",
		"Uint8Array": "reader.readBytes()",
	}

	for fieldType, expectedCall := range expectedDeserialization {
		if !strings.Contains(typesContent, expectedCall) {
			t.Errorf("Deserialization call for %s not found: %s", fieldType, expectedCall)
		}
	}

	// Test wire type handling
	if !strings.Contains(typesContent, "const wireType = tag & 7;") {
		t.Error("Wire type extraction not found")
	}
	if !strings.Contains(typesContent, "switch (fieldNumber) {") {
		t.Error("Field number switch not found")
	}
	if !strings.Contains(typesContent, "reader.skipField(wireType);") {
		t.Error("Field skipping not found")
	}
}

func TestBinarySerializationWithAdvancedTypes(t *testing.T) {
	// Test with advanced types proto that has more complex field types
	cmd := exec.Command("protoc",
		"--plugin=./protoc-gen-grpc-react-native",
		"--grpc-react-native_out=testoutput",
		"-I", "testdata",
		"testdata/basic_types.proto")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %v\nOutput: %s", err, string(output))
	}

	// Test advanced types file
	typesFile := "testoutput/basic_types_types.ts"
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", typesFile, err)
	}
	typesContent := string(content)

	// Test repeated field handling
	if !strings.Contains(typesContent, "for (const item of this.") {
		t.Error("Repeated field serialization loop not found")
	}

	// Test map field handling
	if !strings.Contains(typesContent, "for (const [key, value] of this.") {
		t.Error("Map field serialization loop not found")
	}

	// Test bigint handling
	if !strings.Contains(typesContent, "writeInt64(") {
		t.Error("Int64 serialization not found")
	}
	if !strings.Contains(typesContent, "readInt64()") {
		t.Error("Int64 deserialization not found")
	}
}
