package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the full-stack application (frontend + backend)",
	Long: `Build the full-stack application by compiling both frontend and backend.

This command will:
1. Look for a 'web' directory containing a Vite-powered frontend
2. Detect the package manager (npm, yarn, or pnpm) from lock files
3. Run the frontend build process
4. Find the main.go file (cmd/main.go or nested in cmd/<app>/main.go)
5. Compile the Go backend server

Prerequisites:
- web directory with Vite configuration (vite.config.js/ts)
- package.json with build script in web directory
- main.go file in cmd/ directory or nested subdirectory
- Package manager installed (npm, yarn, or pnpm)`,
	Example: buildExample,
	RunE:    buildRunE,
}

func buildRunE(cmd *cobra.Command, args []string) error {
	// Check if web directory exists
	webDir := "web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		return fmt.Errorf("web directory not found. Make sure you're in a project root with a 'web' directory")
	}

	// Check for Vite configuration
	configPath := filepath.Join(webDir, "vite.config.ts")

	hasViteConfig := false
	if _, err := os.Stat(configPath); err == nil {
		hasViteConfig = true
	}

	if !hasViteConfig {
		return fmt.Errorf("vite configuration not found in web directory. Please ensure vite.config.ts exists")
	}

	// Check for package.json
	packageJsonPath := filepath.Join(webDir, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in web directory")
	}

	// Detect package manager
	packageManager, err := detectPackageManager(webDir)
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	cmd.Printf("Building frontend with %s...\n", packageManager)

	// Build frontend
	if err := buildFrontend(webDir, packageManager); err != nil {
		return fmt.Errorf("frontend build failed: %w", err)
	}

	cmd.Println("Frontend build completed successfully!")
	cmd.Println("Building backend...")

	// Find and build backend
	mainGoPath, err := findMainGo()
	if err != nil {
		return fmt.Errorf("failed to find main.go: %w", err)
	}

	if err := buildBackend(cmd, mainGoPath); err != nil {
		return fmt.Errorf("backend build failed: %w", err)
	}

	cmd.Println("✨ Full-stack build completed successfully!")
	return nil
}

// detectPackageManager detects the package manager based on lock files
func detectPackageManager(webDir string) (string, error) {
	// Check for lock files in order of preference
	lockFiles := map[string]string{
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
		"package-lock.json": "npm",
	}

	for lockFile, manager := range lockFiles {
		lockPath := filepath.Join(webDir, lockFile)
		if _, err := os.Stat(lockPath); err == nil {
			// Verify the package manager is installed
			if _, err := exec.LookPath(manager); err != nil {
				return "", fmt.Errorf("%s lock file found but %s is not installed", lockFile, manager)
			}
			return manager, nil
		}
	}

	// Fallback to npm if available
	if _, err := exec.LookPath("npm"); err == nil {
		return "npm", nil
	}

	return "", fmt.Errorf("no package manager found. Please install npm, yarn, or pnpm")
}

// buildFrontend runs the frontend build process
func buildFrontend(webDir, packageManager string) error {
	var buildCmd *exec.Cmd

	switch packageManager {
	case "npm":
		buildCmd = exec.Command("npm", "run", "build")
	case "yarn":
		buildCmd = exec.Command("yarn", "build")
	case "pnpm":
		buildCmd = exec.Command("pnpm", "build")
	default:
		return fmt.Errorf("unsupported package manager: %s", packageManager)
	}

	buildCmd.Dir = webDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	return buildCmd.Run()
}

// findMainGo finds the main.go file in cmd directory or nested subdirectories
func findMainGo() (string, error) {
	cmdDir := "cmd"
	if _, err := os.Stat(cmdDir); os.IsNotExist(err) {
		return "", fmt.Errorf("cmd directory not found")
	}

	var mainGoPath string

	// Check for cmd/main.go first
	directMainGo := filepath.Join(cmdDir, "main.go")
	if _, err := os.Stat(directMainGo); err == nil {
		return directMainGo, nil
	}

	// Look for main.go in nested cmd subdirectories
	err := filepath.Walk(cmdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "main.go" {
			// Check if this is actually a main package
			if isMainPackage(path) {
				mainGoPath = path
				return filepath.SkipDir // Found it, stop walking
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if mainGoPath == "" {
		return "", fmt.Errorf("main.go not found in cmd directory or its subdirectories")
	}

	return mainGoPath, nil
}

// isMainPackage checks if a Go file belongs to the main package
func isMainPackage(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			return strings.Contains(line, "package main")
		}
	}

	return false
}

// buildBackend compiles the Go backend
func buildBackend(cmd *cobra.Command, mainGoPath string) error {
	// Extract directory and create output binary name
	dir := filepath.Dir(mainGoPath)
	var outputName string

	// Check if custom output name is provided
	if customOutput, _ := cmd.Flags().GetString("output"); customOutput != "" {
		outputName = customOutput
	} else {
		// If it's cmd/main.go, use current directory name
		if filepath.Base(dir) == "cmd" {
			cwd, _ := os.Getwd()
			outputName = filepath.Base(cwd)
		} else {
			// If it's cmd/app/main.go, use the app name
			outputName = filepath.Base(dir)
		}
	}

	// Ensure build directory exists
	buildDir := "build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	outputPath := filepath.Join(buildDir, outputName)

	buildCmd := exec.Command("go", "build", "-o", outputPath, mainGoPath)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	return buildCmd.Run()
}

var buildExample = `
# Build the full-stack application
kirin build

# The command will:
# 1. Look for web/ directory with Vite configuration
# 2. Detect package manager (pnpm, yarn, or npm)
# 3. Run frontend build (npm/yarn/pnpm run build)
# 4. Find main.go in cmd/ directory or subdirectories
# 5. Compile Go backend to build/ directory
`

func init() {
	buildCmd.Flags().StringP("output", "o", "", "Output binary name (default: derived from directory)")
}
