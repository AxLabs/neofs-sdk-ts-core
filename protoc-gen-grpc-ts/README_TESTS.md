# Protoc Plugin Test Suite

This directory contains a comprehensive test suite for the `protoc-gen-grpc-react-native` plugin.

## Test Structure

### Test Files
- `main_test.go` - Main test suite with unit and integration tests
- `integration_test.go` - Comprehensive integration tests
- `test_utils.go` - Test utilities and helper functions
- `run_tests.sh` - Test runner script
- `Makefile` - Make targets for different test scenarios

### Test Data
- `testdata/` - Directory containing test proto files
  - `simple.proto` - Basic message and enum types
  - `nested.proto` - Nested messages and enums
  - `service.proto` - Service definitions with RPC methods
  - `imports.proto` - Cross-package imports and references
  - `edge_cases.proto` - Edge cases and complex scenarios
  - `complex_service.proto` - Complex service with various method patterns

## Test Categories

### 1. Unit Tests
- **Field Naming**: Tests conversion of protobuf field names to camelCase
- **Namespace Generation**: Tests conversion of package names to TypeScript namespaces
- **Type Mapping**: Tests protobuf to TypeScript type conversions
- **File Generation**: Tests correct file structure generation

### 2. Integration Tests
- **Message Generation**: Tests TypeScript interface generation from protobuf messages
- **Enum Generation**: Tests TypeScript enum generation from protobuf enums
- **Service Generation**: Tests service client class generation
- **Import Handling**: Tests cross-file import generation
- **Nested Types**: Tests nested message and enum generation
- **Edge Cases**: Tests complex scenarios and edge cases

### 3. Compilation Tests
- **TypeScript Syntax**: Tests that generated TypeScript files have valid syntax
- **Type Safety**: Tests that generated types are properly typed
- **Import Resolution**: Tests that imports are correctly resolved

### 4. Performance Tests
- **Benchmark Tests**: Tests plugin performance with various file sizes
- **Memory Usage**: Tests memory usage during generation
- **Speed Tests**: Tests generation speed with multiple files

## Running Tests

### Quick Start
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
make test-specific TEST=TestFieldNaming

# Run quick test
make quick-test
```

### Detailed Commands

#### Unit Tests
```bash
# Run unit tests only
make test-unit

# Run with verbose output
make test-verbose

# Run with race detection
make test-race
```

#### Integration Tests
```bash
# Run integration tests only
make test-integration

# Run comprehensive test suite
make run-tests
```

#### Performance Tests
```bash
# Run benchmark tests
make benchmark

# Run full test suite (includes benchmarks)
make full-test
```

### Manual Testing
```bash
# Build the plugin
make build

# Run specific test
go test -v -run "TestFieldNaming"

# Run with coverage
go test -v -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

## Test Coverage

The test suite aims for comprehensive coverage of:
- ✅ Message generation (all field types)
- ✅ Enum generation (including nested enums)
- ✅ Service generation (all RPC method types)
- ✅ Import handling (cross-package references)
- ✅ Nested type generation
- ✅ Field naming conventions
- ✅ Namespace generation
- ✅ TypeScript compilation
- ✅ Error handling
- ✅ Performance benchmarks

## Test Data

### Simple Test (`simple.proto`)
- Basic message with all primitive types
- Simple enum
- Message with enum field

### Nested Test (`nested.proto`)
- Messages with nested messages
- Messages with nested enums
- Complex nested structures

### Service Test (`service.proto`)
- Service with various RPC methods
- Request/response message types
- Cross-package message references

### Imports Test (`imports.proto`)
- Cross-package imports
- Repeated fields
- Optional fields

### Edge Cases Test (`edge_cases.proto`)
- All protobuf field types
- Deeply nested structures
- Complex enum hierarchies
- Reserved fields
- Oneof fields

### Complex Service Test (`complex_service.proto`)
- Complex service with various method patterns
- Map fields
- Complex message types

## Expected Output

### Generated Files
For each proto file, the plugin generates:
- `*_types.ts` - TypeScript interfaces and enums
- `*_services.ts` - Service client classes (if services are defined)

### File Structure
```
testoutput/
├── simple_types.ts
├── nested_types.ts
├── service_types.ts
├── service_services.ts
├── imports_types.ts
├── edge_cases_types.ts
├── complex_service_types.ts
└── complex_service_services.ts
```

### TypeScript Features
- ✅ Proper namespace organization
- ✅ CamelCase field naming
- ✅ Type-safe interfaces
- ✅ Enum generation
- ✅ Import statements
- ✅ Service client classes
- ✅ Proper type annotations

## Troubleshooting

### Common Issues

1. **Protoc not found**
   ```bash
   # Install protoc
   # On Ubuntu/Debian:
   sudo apt-get install protobuf-compiler
   
   # On macOS:
   brew install protobuf
   ```

2. **TypeScript not found**
   ```bash
   # Install TypeScript
   npm install -g typescript
   ```

3. **Test failures**
   ```bash
   # Clean and rebuild
   make clean
   make build
   make test
   ```

### Debug Mode
```bash
# Run tests with verbose output
go test -v -count=1

# Run specific test with debug
go test -v -run "TestFieldNaming" -count=1
```

## Contributing

When adding new tests:
1. Add test cases to appropriate test files
2. Update test data if needed
3. Run full test suite to ensure no regressions
4. Update this README if adding new test categories

## Test Metrics

- **Total Test Cases**: 50+
- **Coverage Target**: 90%+
- **Performance**: < 1s for full test suite
- **Memory Usage**: < 100MB peak
