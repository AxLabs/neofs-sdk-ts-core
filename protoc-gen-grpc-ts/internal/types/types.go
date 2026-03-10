package types

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetTypeScriptTypeForField returns the TypeScript type for a protobuf field
func GetTypeScriptTypeForField(field *protogen.Field, currentFile *protogen.File) string {
	baseType := GetBaseTypeScriptTypeForField(field, currentFile)

	// Handle repeated fields
	if field.Desc.IsList() {
		return baseType + "[]"
	}

	// Handle map fields
	if field.Desc.IsMap() {
		keyType := GetMapKeyType(field.Desc.MapKey().Kind())
		valueType := GetBaseTypeScriptTypeForField(field, currentFile)
		return "Map<" + keyType + ", " + valueType + ">"
	}

	return baseType
}

// GetBaseTypeScriptTypeForField returns the base TypeScript type for a protobuf field
func GetBaseTypeScriptTypeForField(field *protogen.Field, currentFile *protogen.File) string {
	switch field.Desc.Kind() {
	case protoreflect.DoubleKind,
		protoreflect.FloatKind:
		return "number"
	case protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Sint64Kind:
		return "bigint"
	case protoreflect.Int32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Sint32Kind:
		return "number"
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "Uint8Array"
	case protoreflect.MessageKind:
		result := GetMessageTypeName(field.Message, currentFile)
		// fmt.Printf("DEBUG: GetBaseTypeScriptTypeForField message=%s, result=%s\n", field.Message.GoIdent.GoName, result)
		return result
	case protoreflect.EnumKind:
		return GetEnumTypeName(field.Enum, currentFile)
	default:
		return "any"
	}
}

// GetMapKeyType returns the TypeScript type for map keys
func GetMapKeyType(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "bigint"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "number"
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.StringKind:
		return "string"
	default:
		return "string"
	}
}

// GetMessageTypeName returns the TypeScript type name for a message
func GetMessageTypeName(message *protogen.Message, currentFile *protogen.File) string {
	// Check if it's from the same file by comparing paths
	messageFilePath := message.Desc.ParentFile().Path()
	currentFilePath := currentFile.Desc.Path()

	if messageFilePath == currentFilePath {
		// Same file - check if it's a nested message
		if message.Desc.Parent() != nil {
			// Nested message - simplified for now
			return message.GoIdent.GoName
		}
		// Top-level message in same file
		return message.GoIdent.GoName
	}

	// External message - use package prefix
	packageName := GetPackageNameFromDescriptor(message.Desc.ParentFile())
	return packageName + "." + message.GoIdent.GoName
}

// GetMessageImplTypeName returns the TypeScript implementation type name for a message
func GetMessageImplTypeName(message *protogen.Message, currentFile *protogen.File) string {
	// Check if it's from the same file by comparing paths
	messageFilePath := message.Desc.ParentFile().Path()
	currentFilePath := currentFile.Desc.Path()

	if messageFilePath == currentFilePath {
		// Same file - check if it's a nested message
		if message.Desc.Parent() != nil {
			// Nested message - simplified for now
			return message.GoIdent.GoName + "Impl"
		}
		// Top-level message in same file
		return message.GoIdent.GoName + "Impl"
	}

	// External message - use package prefix
	packageName := GetPackageNameFromDescriptor(message.Desc.ParentFile())
	return packageName + "." + message.GoIdent.GoName + "Impl"
}

// GetEnumTypeName returns the TypeScript type name for an enum
func GetEnumTypeName(enum *protogen.Enum, currentFile *protogen.File) string {
	// Nested enum - use just the name
	if enum.Desc.Parent() != nil {
		return enum.GoIdent.GoName
	}

	// Check if it's from the same file by comparing full paths
	enumFilePath := enum.Desc.ParentFile().Path()
	currentFilePath := currentFile.Desc.Path()

	if enumFilePath != currentFilePath {
		// External enum - use package prefix
		packageName := GetPackageNameFromDescriptor(enum.Desc.ParentFile())
		return packageName + "." + enum.GoIdent.GoName
	}

	// Same file enum
	return enum.GoIdent.GoName
}

// GetPackageNameFromDescriptor extracts the package name from a file descriptor
func GetPackageNameFromDescriptor(fileDesc protoreflect.FileDescriptor) string {
	packageName := string(fileDesc.Package())
	if packageName == "" {
		return ""
	}

	// Convert package name to TypeScript namespace format
	parts := strings.Split(packageName, ".")
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, ToCamelCase(part))
		}
	}
	return strings.Join(result, "")
}

// GetPackageName extracts the package name from a proto file
func GetPackageName(file *protogen.File) string {
	packageName := string(file.Desc.Package())
	if packageName == "" {
		return ""
	}

	// Convert package name to TypeScript namespace format
	parts := strings.Split(packageName, ".")
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, ToCamelCase(part))
		}
	}
	return strings.Join(result, "")
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Handle special cases
	switch s {
	case "v2":
		return "V2"
	case "v1":
		return "V1"
	default:
		// Convert snake_case to camelCase
		parts := strings.Split(s, "_")
		if len(parts) == 1 {
			// No underscores, just capitalize first letter
			return strings.ToUpper(s[:1]) + s[1:]
		}

		// Multiple parts - first part lowercase, rest capitalized
		result := strings.ToLower(parts[0])
		for _, part := range parts[1:] {
			if part != "" {
				result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
		}
		return result
	}
}

// GetDefaultValueForField returns the default value for a field
func GetDefaultValueForField(field *protogen.Field, file *protogen.File) string {
	// Handle repeated fields
	if field.Desc.IsList() {
		return "[]"
	}

	// Handle map fields
	if field.Desc.IsMap() {
		return "new Map()"
	}

	switch field.Desc.Kind() {
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		return "0.0"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind, protoreflect.Sint64Kind:
		return "0n"
	case protoreflect.Int32Kind, protoreflect.Fixed32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		return "0"
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return `""`
	case protoreflect.BytesKind:
		return "new Uint8Array(0)"
	case protoreflect.MessageKind:
		// For message types, use undefined as the default
		// Creating new instances causes infinite recursion for self-referencing types
		return "undefined"
	case protoreflect.EnumKind:
		return "0"
	default:
		return "undefined"
	}
}
