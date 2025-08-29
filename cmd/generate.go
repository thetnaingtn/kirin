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
	initProtoFolder     string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "Generate code from protobuf definitions using buf",
	Long: `Generate code from protobuf definitions using buf.
Prerequisites:
- buf must be installed (https://buf.build/docs/installation)
- proto directory with valid protobuf files
- buf.yaml configuration file in the proto directory
- buf.gen.yaml configuration file in the proto directory (for generation)`,
	Example: generateExample,
	RunE:    generateRunE,
}

// init subcommand for generate
var generateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize buf configuration files in existing proto directory",
	Long: `Initialize buf configuration files (buf.yaml and buf.gen.yaml) in an existing proto directory.
This command assumes you already have a proto directory with .proto files and adds the necessary 
buf configuration files to enable code generation.`,
	Example: `# Initialize buf config in default proto directory
kirin generate init

# Initialize buf config in custom proto directory  
kirin generate init --proto-folder=api/proto`,
	RunE: generateInitRunE,
}

func init() {
	generateCmd.Flags().StringVar(&generateProtoFolder, "proto-folder", "proto", "Proto directory name (default: proto)")
	generateInitCmd.Flags().StringVar(&initProtoFolder, "proto-folder", "proto", "Proto directory name (default: proto)")
	generateCmd.AddCommand(generateInitCmd)
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

	// Check if buf.yaml exists in proto directory
	bufYamlPath := filepath.Join(generateProtoFolder, "buf.yaml")
	if _, err := os.Stat(bufYamlPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file buf.yaml not found in %s directory. Please ensure buf.yaml exists", generateProtoFolder)
	}

	// Check if buf.gen.yaml exists in proto directory
	bufGenYamlPath := filepath.Join(generateProtoFolder, "buf.gen.yaml")
	if _, err := os.Stat(bufGenYamlPath); os.IsNotExist(err) {
		return fmt.Errorf("generation config buf.gen.yaml not found in %s directory. Please ensure buf.gen.yaml exists", generateProtoFolder)
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

	if err := buildCmd.Run(); err != nil {
		cmd.Println("Validation failed!")
		return fmt.Errorf("buf build failed with validation errors. Please fix the protobuf issues above before generating code")
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

func generateInitRunE(cmd *cobra.Command, args []string) error {
	// Load existing configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Apply configuration values if flags weren't explicitly set
	if !cmd.Flags().Changed("proto-folder") && cfg.Generate.ProtoFolder != "" {
		initProtoFolder = cfg.Generate.ProtoFolder
	}

	// Check if proto directory exists
	if _, err := os.Stat(initProtoFolder); os.IsNotExist(err) {
		return fmt.Errorf("proto directory '%s' not found. Please ensure the proto directory exists with your protobuf files", initProtoFolder)
	}

	bufYamlPath := filepath.Join(initProtoFolder, "buf.yaml")
	bufGenYamlPath := filepath.Join(initProtoFolder, "buf.gen.yaml")

	// Check if both configuration files already exist
	bufYamlExists := true
	bufGenYamlExists := true

	if _, err := os.Stat(bufYamlPath); os.IsNotExist(err) {
		bufYamlExists = false
	}

	if _, err := os.Stat(bufGenYamlPath); os.IsNotExist(err) {
		bufGenYamlExists = false
	}

	// If both files already exist, inform user and exit
	if bufYamlExists && bufGenYamlExists {
		cmd.Printf("ℹ️  Buf configuration files already exist in %s\n", initProtoFolder)
		cmd.Printf("   - buf.yaml ✅\n")
		cmd.Printf("   - buf.gen.yaml ✅\n")
		cmd.Println("Your proto directory is already configured for buf!")
		cmd.Println("You can run 'kirin generate' to generate code from your protobuf files.")
		return nil
	}

	// Create buf.yaml if it doesn't exist
	if !bufYamlExists {
		if err := os.WriteFile(bufYamlPath, []byte(bufYamlTemplate), 0644); err != nil {
			return fmt.Errorf("failed to create buf.yaml: %w", err)
		}
		cmd.Printf("✅ Created buf.yaml in %s\n", initProtoFolder)
	}

	// Create buf.gen.yaml if it doesn't exist
	if !bufGenYamlExists {
		if err := os.WriteFile(bufGenYamlPath, []byte(getBufGenYamlTemplate()), 0644); err != nil {
			return fmt.Errorf("failed to create buf.gen.yaml: %w", err)
		}
		cmd.Printf("✅ Created buf.gen.yaml in %s\n", initProtoFolder)
	}

	// Save configuration if any non-default values were used
	if config.HasNonDefaultGenerateValues(initProtoFolder) {
		cfg.Generate.ProtoFolder = initProtoFolder

		if err := config.SaveConfig(cfg); err != nil {
			cmd.Printf("Warning: failed to save configuration: %v\n", err)
		} else {
			cmd.Println("Configuration saved to .kirin.toml")
		}
	}

	cmd.Printf("Buf configuration initialized successfully in %s!\n", initProtoFolder)
	cmd.Println("You can now run 'kirin generate' to generate code from your protobuf files.")

	return nil
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
