package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thetnaingtn/kirin/internal/config"
)

var (
	generateProtoFolder string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "Generate code from protobuf definitions using buf",
	Example: generateExample,
	RunE:    generateRunE,
}

func init() {
	generateCmd.Flags().StringVar(&generateProtoFolder, "proto-folder", "proto", "Proto directory name (default: proto)")
}

func generateRunE(cmd *cobra.Command, args []string) error {
	// Load existing configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Apply configuration values if flags weren't explicitly set
	if !cmd.Flags().Changed("proto-folder") && cfg.Generate.ProtoFolder != "" {
		generateProtoFolder = cfg.Generate.ProtoFolder
	}

	// Check if buf is installed
	if _, err := exec.LookPath("buf"); err != nil {
		return fmt.Errorf("buf is not installed or not found in PATH. Please install buf from https://buf.build/docs/installation")
	}

	// Look for proto directory
	if _, err := os.Stat(generateProtoFolder); os.IsNotExist(err) {
		return fmt.Errorf("proto directory '%s' not found. Make sure you're in a project root with a proto directory", generateProtoFolder)
	}

	// Check if proto directory contains .proto files
	hasProtoFiles, err := hasProtobufFiles(generateProtoFolder)
	if err != nil {
		return fmt.Errorf("error checking proto directory: %w", err)
	}
	if !hasProtoFiles {
		return fmt.Errorf("no .proto files found in the '%s' directory", generateProtoFolder)
	}

	// Check if both configuration files already exist
	bufYamlExists := true
	bufGenYamlExists := true

	bufYamlPath := filepath.Join(generateProtoFolder, "buf.yaml")
	bufGenYamlPath := filepath.Join(generateProtoFolder, "buf.gen.yaml")

	if _, err := os.Stat(bufYamlPath); os.IsNotExist(err) {
		bufYamlExists = false
	}

	if _, err := os.Stat(bufGenYamlPath); os.IsNotExist(err) {
		bufGenYamlExists = false
	}

	// Check if buf.yaml exists in proto directory
	if !bufYamlExists {
		cmd.Println("Configuration file buf.yaml not found. Creating a buf.yaml...")

		if err := os.WriteFile(bufYamlPath, []byte(bufYamlTemplate), 0644); err != nil {
			return fmt.Errorf("failed to create buf.yaml: %w", err)
		}

		cmd.Printf("✅ Created buf.yaml in %s\n", generateProtoFolder)
	}

	// Create buf.gen.yaml if it doesn't exist
	if !bufGenYamlExists {
		cmd.Println("Configuration file buf.gen.yaml not found. Creating a buf.gen.yaml...")
		if err := os.WriteFile(bufGenYamlPath, []byte(getBufGenYamlTemplate()), 0644); err != nil {
			return fmt.Errorf("failed to create buf.gen.yaml: %w", err)
		}
		cmd.Printf("✅ Created buf.gen.yaml in %s\n", generateProtoFolder)
	}

	// Get current working directory to restore later
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Change to proto directory
	if err := os.Chdir(generateProtoFolder); err != nil {
		return fmt.Errorf("failed to change to proto directory: %w", err)
	}

	cmd.Println("Validating protobuf files with buf build...")

	// Run buf build to validate (from within proto directory)
	buildCmd := exec.Command("buf", "build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	err = buildCmd.Run()

	if err != nil {
		cmd.Println("Updating dependencies with buf dep update...")
		// try to update dependencies
		depsUpdateCmd := exec.Command("buf", "dep", "update")
		depsUpdateCmd.Stdout = os.Stdout
		depsUpdateCmd.Stderr = os.Stderr

		if err := depsUpdateCmd.Run(); err != nil {
			cmd.Println("Failed to update dependencies!")
			return fmt.Errorf("buf build failed with validation errors. Please fix the protobuf issues above before generating code")
		}
	}

	cmd.Println("Validation passed!")
	cmd.Println("Generating code from protobuf definitions...")

	// Run buf generate (from within proto directory)
	generateCmd := exec.Command("buf", "generate")
	generateCmd.Stdout = os.Stdout
	generateCmd.Stderr = os.Stderr

	if err := generateCmd.Run(); err != nil {
		return fmt.Errorf("buf generate failed: %w", err)
	}

	// Change back to original directory before saving config
	if err := os.Chdir(currentDir); err != nil {
		cmd.Printf("Warning: failed to change back to original directory: %v\n", err)
	}

	// Save configuration if any non-default values were used
	if config.HasNonDefaultGenerateValues(generateProtoFolder) {
		cfg.Generate.ProtoFolder = generateProtoFolder

		if err := config.SaveConfig(cfg); err != nil {
			cmd.Printf("Warning: failed to save configuration: %v\n", err)
		} else {
			cmd.Println("Configuration saved to .kirin.toml")
		}
	}

	cmd.Println("✨ Code generation completed successfully!")
	return nil
}

// hasProtobufFiles checks if the directory contains any .proto files
func hasProtobufFiles(dir string) (bool, error) {
	var hasProto bool

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			hasProto = true
			return filepath.SkipDir // We found at least one, no need to continue
		}
		return nil
	})

	return hasProto, err
}

// getModuleName reads the module name from go.mod file in the current directory
func getModuleName() string {
	goModPath := "go.mod"

	// Check if go.mod exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "your/module/name" // fallback if no go.mod found
	}

	file, err := os.Open(goModPath)
	if err != nil {
		return "your/module/name" // fallback on error
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if after, ok := strings.CutPrefix(line, "module "); ok {
			moduleName := after
			moduleName = strings.TrimSpace(moduleName)
			return moduleName
		}
	}

	return "your/module/name" // fallback if module line not found
}

// getBufGenYamlTemplate returns the buf.gen.yaml template with the current module name
func getBufGenYamlTemplate() string {
	moduleName := getModuleName()
	return fmt.Sprintf(`version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package
      module: buf.build/googleapis/googleapis
  override:
    - file_option: go_package_prefix
      value: %s
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: gen
    opt: paths=source_relative
  - remote: buf.build/grpc-ecosystem/gateway
    out: gen
    opt: paths=source_relative
  - remote: buf.build/community/stephenh-ts-proto
    out: ../%s/src/types/proto
    opt:
      - env=browser
      - useOptionals=messages
      - outputServices=generic-definitions
      - outputJsonMethods=false
      - useExactTypes=false
      - esModuleInterop=true
      - stringEnums=true
`, moduleName, buildFrontendFolder)
}

var (
	bufYamlTemplate = `version: v2
deps:
  - buf.build/googleapis/googleapis
lint:
  use:
    - BASIC
  except:
    - ENUM_VALUE_PREFIX
    - FIELD_NOT_REQUIRED
    - PACKAGE_DIRECTORY_MATCH
    - PACKAGE_NO_IMPORT_CYCLE
    - PACKAGE_VERSION_SUFFIX
  disallow_comment_ignores: true
breaking:
  use:
    - FILE
  except:
    - EXTENSION_NO_DELETE
    - FIELD_SAME_DEFAULT
`

	generateExample = `
# Generate code from protobuf definitions
kirin generate

# Or use aliases
kirin gen
kirin g

# With custom proto directory
kirin generate --proto-folder protos
kirin gen --proto-folder api/proto
`
)
