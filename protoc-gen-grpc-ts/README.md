# protoc-gen-grpc-react-native

A Protocol Buffers compiler plugin that generates TypeScript code for React Native gRPC clients with binary protobuf serialization support.

## Features

- **TypeScript Generation**: Generates TypeScript interfaces, classes, and enums
- **Binary Serialization**: Full protobuf wire format support compatible with grpc-react-native
- **Streaming Support**: All gRPC streaming patterns (unary, server, client, bidirectional)
- **Message Classes**: Generates classes with `serializeBinary()`, `deserializeBinary()`, and `toObject()` methods
- **Field Type Support**: All protobuf field types including repeated fields and maps
- **Cross-package Imports**: Proper handling of dependencies between packages
- **React Native Compatible**: Optimized for React Native with binary data handling

## Quick Start

### Build the Plugin

The easiest way to build the plugin is using the Makefile:

```bash
# Build the plugin (default target)
make

# Or explicitly
make build
```

This will create `protoc-gen-grpc-react-native` in the current directory.

### Install the Plugin (Optional)

To install the plugin globally so you can use it without specifying the path:

```bash
# Install to /usr/local/bin (requires sudo)
sudo make install

# Or install to a custom directory (e.g., ~/bin)
make install INSTALL_DIR=~/bin
```

### Generate TypeScript Code

Once built, use the plugin with `protoc`:

```bash
# Basic usage
protoc --plugin=./protoc-gen-grpc-react-native \
       --grpc-react-native_out=./output \
       -I./proto \
       ./proto/service.proto
```

## Installation

### Prerequisites

- Go 1.19 or later
- Protocol Buffers compiler (`protoc`)
- grpc-react-native client library

### Build the Plugin

```bash
cd grpc-react-native-plugin
make build
# or
go build -o protoc-gen-grpc-react-native main.go
```

### Install Plugin Globally (Optional)

To use the plugin without specifying the full path each time, you can install it to a directory in your `PATH`:

```bash
# Build the plugin
go build -o protoc-gen-grpc-react-native

# Install to a directory in PATH (e.g., /usr/local/bin or ~/bin)
sudo cp protoc-gen-grpc-react-native /usr/local/bin/
# or
cp protoc-gen-grpc-react-native ~/bin/

# Make sure the directory is in your PATH
export PATH=$PATH:~/bin
```

After installation, you can use the plugin without the `--plugin` flag:

```bash
protoc --grpc-react-native_out=./output -I./proto ./proto/service.proto
```

## Usage

### Prerequisites

1. **Build the plugin** (if not already built):
```bash
go build -o protoc-gen-grpc-react-native
```

2. **Install protoc** (if not already installed):
```bash
# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# macOS
brew install protobuf

# Or download from: https://github.com/protocolbuffers/protobuf/releases
```

### Basic Usage

#### Single Proto File
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/your_service.proto
```

#### Multiple Proto Files
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/service1.proto ./proto/service2.proto
```

#### All Proto Files in Directory
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/*.proto
```

### Advanced Usage Examples

#### With Custom Output Directory
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./src/generated -I./proto ./proto/accounting/service.proto
```

#### With Multiple Include Paths
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto -I./common ./proto/service.proto
```

#### With NeoFS Protobuf Definitions
```bash
# Single service
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./neofs-api neofs-api/accounting/service.proto

# All NeoFS services
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./neofs-api neofs-api/accounting/service.proto neofs-api/object/service.proto neofs-api/container/service.proto
```

#### With Streaming Services
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/streaming_service.proto
```

#### With Complex Types (Maps, Repeated Fields)
```bash
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/advanced_types.proto
```

### Generated Files

The plugin generates the following files:

- **`*_types.ts`** - TypeScript interfaces, classes, enums, and binary serialization
- **`*_services.ts`** - gRPC service clients (only generated if proto contains services)

#### Example Output Structure
```
output/
├── simple_types.ts           # Simple message types
├── service_types.ts          # Service message types  
├── service_services.ts       # Service client
├── streaming_service_types.ts # Streaming message types
└── streaming_service_services.ts # Streaming service client
```

### Command Line Options

#### Plugin Options
```bash
# Basic usage
--grpc-react-native_out=./output

# With custom options (future enhancement)
--grpc-react-native_out=./output:option1=value1,option2=value2
```

#### Include Paths (`-I` flag)
```bash
# Single include path
-I./proto

# Multiple include paths
-I./proto -I./common -I./third_party

# Relative to current directory
-I. -I./proto -I./neofs-api
```

### Integration with Build Systems

#### Makefile Example
```makefile
PROTO_FILES := $(wildcard proto/*.proto)
GENERATED_FILES := $(PROTO_FILES:.proto=_types.ts) $(PROTO_FILES:.proto=_services.ts)

.PHONY: generate
generate: $(GENERATED_FILES)

%_types.ts %_services.ts: %.proto
	protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=. -I./proto $<
```

#### npm Script Example
```json
{
  "scripts": {
    "generate": "protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./src/generated -I./proto ./proto/*.proto",
    "generate:neofs": "protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./src/generated -I./neofs-api neofs-api/accounting/service.proto neofs-api/object/service.proto"
  }
}
```

#### Gradle Example (Android)
```gradle
task generateProto {
    doLast {
        exec {
            commandLine 'protoc', 
                '--plugin=./protoc-gen-grpc-react-native',
                '--grpc-react-native_out=./src/generated',
                '-I./proto',
                './proto/service.proto'
        }
    }
}
```

### Troubleshooting

#### Common Issues

1. **Plugin not found**:
```bash
# Make sure plugin is built and executable
chmod +x protoc-gen-grpc-react-native
```

2. **Import errors**:
```bash
# Use correct include paths
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto -I./common ./proto/service.proto
```

3. **Permission denied**:
```bash
# Make plugin executable
chmod +x protoc-gen-grpc-react-native
```

4. **Output directory doesn't exist**:
```bash
# Create output directory first
mkdir -p output
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto ./proto/service.proto
```

#### Debug Mode
```bash
# Verbose output
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./output -I./proto --verbose ./proto/service.proto
```

## Generated Code

### Message Classes

The plugin generates TypeScript classes with serialization methods:

```typescript
export class SimpleMessageImpl {
  constructor(data?: Partial<SimpleMessage>) {
    this.name = data?.name ?? "";
    this.age = data?.age ?? 0;
    this.active = data?.active ?? false;
  }

  name!: string;
  age!: number;
  active!: boolean;

  serializeBinary(): Uint8Array {
    // Binary protobuf serialization
  }

  static deserializeBinary(data: Uint8Array): SimpleMessageImpl {
    // Binary protobuf deserialization
  }

  toObject(): SimpleMessage {
    // Convert to plain object
  }
}
```

### Service Clients

The plugin generates service client classes with streaming support:

```typescript
export class ObjectServiceClient {
  constructor(private client: GrpcClient) {}

  // Unary call
  async delete(request: DeleteRequestImpl): Promise<DeleteResponseImpl> {
    const response = await this.client.unaryCall(
      'neo.fs.v2.object.ObjectService/Delete',
      request.serializeBinary()
    );
    return DeleteResponseImpl.deserializeBinary(response.data as Uint8Array);
  }

  // Server streaming
  async* get(request: GetRequestImpl): AsyncGenerator<GetResponseImpl> {
    const stream = this.client.serverStreamCall(
      'neo.fs.v2.object.ObjectService/Get',
      request.serializeBinary()
    );
    for await (const response of stream) {
      if (response.done) {
        break; // End of stream
      }
      if (response.data) {
        yield GetResponseImpl.deserializeBinary(response.data as Uint8Array);
      }
    }
  }

  // Client streaming - accepts AsyncIterable of requests
  async put(requests: AsyncIterable<PutRequestImpl>): Promise<PutResponseImpl> {
    async function* serialize() {
      for await (const req of requests) {
        yield req.serializeBinary();
      }
    }
    const response = await this.client.clientStreamCall(
      'neo.fs.v2.object.ObjectService/Put',
      serialize()
    );
    return PutResponseImpl.deserializeBinary(response.data as Uint8Array);
  }

  // Bidirectional streaming
  async* chat(requests: AsyncIterable<ChatMessageImpl>): AsyncGenerator<ChatMessageImpl> {
    async function* serialize() {
      for await (const req of requests) {
        yield req.serializeBinary();
      }
    }
    const stream = this.client.bidirectionalStreamCall(
      'neo.fs.v2.object.ObjectService/Chat',
      serialize()
    );
    for await (const response of stream) {
      if (response.done) {
        break; // End of stream
      }
      if (response.data) {
        yield ChatMessageImpl.deserializeBinary(response.data as Uint8Array);
      }
    }
  }
}
```

### Type Definitions

The plugin generates TypeScript interfaces and enums:

```typescript
export interface SimpleMessage {
  name: string;
  age: number;
  active: boolean;
  score: number;
  data: Uint8Array;
}

export enum Status {
  STATUS_UNSPECIFIED = 0,
  STATUS_ACTIVE = 1,
  STATUS_INACTIVE = 2,
}
```

## Supported Features

### Field Types

- **Scalar Types**: int32, int64, uint32, uint64, sint32, sint64, fixed32, fixed64, sfixed32, sfixed64, float, double, bool, string, bytes
- **Repeated Fields**: Arrays of any scalar or message type
- **Map Fields**: Maps with string/numeric keys
- **Nested Messages**: Proper type references for nested message types
- **Enums**: TypeScript enums with proper value mappings

### Streaming Methods

- **Unary**: `async method(request: RequestImpl): Promise<ResponseImpl>`
- **Server Streaming**: `async* method(request: RequestImpl): AsyncGenerator<ResponseImpl>`
- **Client Streaming**: `async method(requests: AsyncIterable<RequestImpl>): Promise<ResponseImpl>`
- **Bidirectional Streaming**: `async* method(requests: AsyncIterable<RequestImpl>): AsyncGenerator<ResponseImpl>`

### Binary Serialization

The plugin includes a complete BinaryWriter and BinaryReader implementation:

```typescript
class BinaryWriter {
  writeString(fieldNumber: number, value: string): void
  writeBytes(fieldNumber: number, value: Uint8Array): void
  writeInt32(fieldNumber: number, value: number): void
  writeMessage(fieldNumber: number, value: any): void
  // ... more methods
}

class BinaryReader {
  readString(fieldNumber: number): string
  readBytes(fieldNumber: number): Uint8Array
  readInt32(fieldNumber: number): number
  readMessage(fieldNumber: number, deserializer: Function): any
  // ... more methods
}
```

## Integration with grpc-react-native

The generated code is designed to work seamlessly with the grpc-react-native client:

```typescript
import { GrpcClient } from 'grpc-react-native';
import { AccountingServiceClient } from './generated/accounting_service_services';
import { BalanceRequestImpl } from './generated/accounting_types';

// Create gRPC client
const grpcClient = new GrpcClient({
  host: 'localhost',
  port: 50051,
  useTls: false
});

// Create service client
const accountingClient = new AccountingServiceClient(grpcClient);

// Create request message
const request = new BalanceRequestImpl({
  ownerId: "test-owner-id"
});

// Make gRPC call
const response = await accountingClient.balance(request);
console.log('Response:', response);
```

## Examples

### Simple Message

```protobuf
syntax = "proto3";

message SimpleMessage {
  string name = 1;
  int32 age = 2;
  bool active = 3;
  double score = 4;
  bytes data = 5;
}
```

Generated TypeScript:

```typescript
export interface SimpleMessage {
  name: string;
  age: number;
  active: boolean;
  score: number;
  data: Uint8Array;
}

export class SimpleMessageImpl {
  constructor(data?: Partial<SimpleMessage>) {
    this.name = data?.name ?? "";
    this.age = data?.age ?? 0;
    this.active = data?.active ?? false;
    this.score = data?.score ?? 0.0;
    this.data = data?.data ?? new Uint8Array();
  }

  serializeBinary(): Uint8Array { /* ... */ }
  static deserializeBinary(data: Uint8Array): SimpleMessageImpl { /* ... */ }
  toObject(): SimpleMessage { /* ... */ }
}
```

### Streaming Service

```protobuf
service StreamingService {
  rpc GetData(GetDataRequest) returns (stream GetDataResponse);
  rpc UploadData(stream UploadDataRequest) returns (UploadDataResponse);
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

Generated TypeScript:

```typescript
export class StreamingServiceClient {
  constructor(private client: GrpcClient) {}

  // Server streaming
  async* getData(request: GetDataRequestImpl): AsyncGenerator<GetDataResponseImpl> {
    const stream = this.client.serverStreamCall(
      'test.streaming.StreamingService/GetData',
      request.serializeBinary()
    );
    for await (const response of stream) {
      if (response.done) {
        break;
      }
      if (response.data) {
        yield GetDataResponseImpl.deserializeBinary(response.data as Uint8Array);
      }
    }
  }

  // Client streaming - accepts AsyncIterable
  async uploadData(requests: AsyncIterable<UploadDataRequestImpl>): Promise<UploadDataResponseImpl> {
    async function* serialize() {
      for await (const req of requests) {
        yield req.serializeBinary();
      }
    }
    const response = await this.client.clientStreamCall(
      'test.streaming.StreamingService/UploadData',
      serialize()
    );
    return UploadDataResponseImpl.deserializeBinary(response.data as Uint8Array);
  }

  // Bidirectional streaming
  async* chat(requests: AsyncIterable<ChatMessageImpl>): AsyncGenerator<ChatMessageImpl> {
    async function* serialize() {
      for await (const req of requests) {
        yield req.serializeBinary();
      }
    }
    const stream = this.client.bidirectionalStreamCall(
      'test.streaming.StreamingService/Chat',
      serialize()
    );
    for await (const response of stream) {
      if (response.done) {
        break;
      }
      if (response.data) {
        yield ChatMessageImpl.deserializeBinary(response.data as Uint8Array);
      }
    }
  }
}
```

### Usage Example

```typescript
import { GrpcClient } from 'grpc-react-native';
import { StreamingServiceClient } from './generated/streaming_service_services';
import { UploadDataRequestImpl } from './generated/streaming_service_types';

const client = new GrpcClient({
  host: 'localhost',
  port: 50051,
  useTls: false
});

await client.initialize();

const serviceClient = new StreamingServiceClient(client);

// Client streaming example
async function* generateUploadRequests() {
  for (let i = 0; i < 10; i++) {
    yield new UploadDataRequestImpl({
      chunk: new TextEncoder().encode(`chunk-${i}`)
    });
  }
}

const response = await serviceClient.uploadData(generateUploadRequests());
console.log('Upload complete:', response);
```

## Testing

The plugin includes comprehensive test cases covering all features:

### Running Tests

#### All Tests
```bash
# Run all tests (some legacy tests may fail due to file structure changes)
go test -v
```

#### New Comprehensive Tests Only
```bash
# Run only the new comprehensive tests that cover all features
go test -v -run "TestStreaming|TestBinary|TestType|TestDefault|TestRepeated|TestMap|TestNested|TestEnum|TestFieldNaming|TestNamespaceGeneration|TestErrorHandling|TestPerformance"
```

#### Individual Test Categories
```bash
# Test streaming methods
go test -v -run "TestStreaming"

# Test binary serialization
go test -v -run "TestBinary"

# Test type mapping
go test -v -run "TestType"

# Test default values
go test -v -run "TestDefault"

# Test repeated fields
go test -v -run "TestRepeated"

# Test map fields
go test -v -run "TestMap"

# Test nested messages
go test -v -run "TestNested"

# Test enum generation
go test -v -run "TestEnum"
```

### Test Coverage

The test suite covers:

- **Streaming Methods**: Server, client, and bidirectional streaming
- **Binary Serialization**: Complete BinaryWriter/BinaryReader functionality
- **Type Mapping**: All protobuf types (string, number, bigint, Uint8Array, etc.)
- **Default Values**: Proper initialization of all field types
- **Repeated Fields**: Array handling and serialization
- **Map Fields**: Map<K,V> generation and serialization
- **Nested Messages**: Complex message structures
- **Enum Generation**: TypeScript enum generation
- **Field Naming**: PascalCase conversion
- **Namespace Generation**: Proper namespace handling
- **Error Handling**: Edge cases and error conditions
- **Performance**: Large message handling

### Manual Testing

#### Test with Sample Proto Files
```bash
# Test with simple message
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./testoutput -I./testdata testdata/simple.proto

# Test with streaming service
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./testoutput -I./testdata testdata/streaming_service.proto

# Test with complex types
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./testoutput -I./testdata testdata/basic_types.proto

# Test with NeoFS protobuf definitions
protoc --plugin=./protoc-gen-grpc-react-native --grpc-react-native_out=./testoutput -I./testdata/neofs testdata/neofs/accounting/service.proto
```

#### Verify Generated Code
```bash
# Check generated files
ls -la testoutput/

# Verify TypeScript syntax (if tsc is available)
npx tsc --noEmit testoutput/*.ts
```

## Development

### Building

```bash
go build -o protoc-gen-grpc-react-native
```

### Testing

```bash
go test ./...
```

### Adding New Features

1. Modify `main.go` to add new generation logic
2. Add test cases in `testdata/`
3. Update expected output in `testoutput/`
4. Run tests to verify changes

## License

This project is part of the NeoFS TypeScript SDK and follows the same license terms.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Changelog

### v1.0.0
- Initial release with TypeScript generation
- Binary protobuf serialization support
- Full streaming method support
- Cross-package import handling
- NeoFS protobuf compatibility
