package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
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

func generateRunE(cmd *cobra.Command, args []string) error {
	// Check if buf is installed
	if _, err := exec.LookPath("buf"); err != nil {
		return fmt.Errorf("buf is not installed or not found in PATH. Please install buf from https://buf.build/docs/installation")
	}

	// Get proto directory from flag
	protoDir, _ := cmd.Flags().GetString("proto-folder")
	if _, err := os.Stat(protoDir); os.IsNotExist(err) {
		return fmt.Errorf("proto directory '%s' not found. Make sure you're in a project root with a proto directory", protoDir)
	}

	// Check if proto directory contains .proto files
	hasProtoFiles, err := hasProtobufFiles(protoDir)
	if err != nil {
		return fmt.Errorf("error checking proto directory: %w", err)
	}
	if !hasProtoFiles {
		return fmt.Errorf("no .proto files found in the '%s' directory", protoDir)
	}

	// Check if buf.yaml exists in proto directory
	bufYamlPath := filepath.Join(protoDir, "buf.yaml")
	if _, err := os.Stat(bufYamlPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file buf.yaml not found in %s directory. Please ensure buf.yaml exists", protoDir)
	}

	// Check if buf.gen.yaml exists in proto directory
	bufGenYamlPath := filepath.Join(protoDir, "buf.gen.yaml")
	if _, err := os.Stat(bufGenYamlPath); os.IsNotExist(err) {
		return fmt.Errorf("generation config buf.gen.yaml not found in %s directory. Please ensure buf.gen.yaml exists", protoDir)
	}

	// Get current working directory to restore later
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Change to proto directory
	if err := os.Chdir(protoDir); err != nil {
		return fmt.Errorf("failed to change to proto directory: %w", err)
	}

	// Ensure we change back to original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			cmd.Printf("Warning: failed to change back to original directory: %v\n", err)
		}
	}()

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

var generateExample = `
# Generate code from protobuf definitions
kirin generate

# Or use aliases
kirin gen
kirin g

# With custom proto directory
kirin generate --proto-folder protos
kirin gen --proto-folder api/proto
`

func init() {
	generateCmd.Flags().StringP("proto-folder", "p", "proto", "Proto directory name (default: proto)")
}
