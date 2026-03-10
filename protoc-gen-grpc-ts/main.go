package main

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/neo-fs/protoc-gen-grpc-ts/internal/generator"
	"github.com/neo-fs/protoc-gen-grpc-ts/internal/services"
)

func main() {
	var target string = "react-native" // default target

	protogen.Options{
		ParamFunc: func(name, value string) error {
			if name == "target" {
				target = value
			}
			return nil
		},
	}.Run(func(gen *protogen.Plugin) error {
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			// Generate types file for all files (shared between platforms)
			// Pass target so generator can create appropriate exports
			generator.GenerateTypesFile(gen, file, target)

			// Generate services file only if there are services
			if len(file.Services) > 0 {
				// Route to appropriate service generator based on target
				switch strings.ToLower(target) {
				case "nodejs", "node":
					services.GenerateNodeJSServicesFile(gen, file)
				case "react-native", "reactnative", "rn":
					services.GenerateReactNativeServicesFile(gen, file)
				default:
					// Default to react-native for backward compatibility
					services.GenerateReactNativeServicesFile(gen, file)
				}
			}
		}
		return nil
	})
}
