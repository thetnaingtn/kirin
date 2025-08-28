package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thetnaingtn/kirin/internal/config"
)

var (
	buildOutput         string
	buildFrontendFolder string
	buildMainFolder     string
	buildPkgManager     string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the full-stack application (frontend + backend)",
	Long: `Build the full-stack application by compiling both frontend and backend.
	
Prerequisites:
- Frontend directory with Vite configuration (vite.config.js/ts)
- package.json with build script in frontend directory
- main.go file in main directory or nested subdirectory
- Package manager installed (npm, yarn, or pnpm)`,
	Example: buildExample,
	RunE:    buildRunE,
}

func buildRunE(cmd *cobra.Command, args []string) error {
	// Load existing configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Apply configuration values if flags weren't explicitly set
	if !cmd.Flags().Changed("frontend-folder") && cfg.Build.FrontendFolder != "" {
		buildFrontendFolder = cfg.Build.FrontendFolder
	}
	if !cmd.Flags().Changed("main-folder") && cfg.Build.MainFolder != "" {
		buildMainFolder = cfg.Build.MainFolder
	}
	if !cmd.Flags().Changed("pkg-manager") && cfg.Build.PkgManager != "" {
		buildPkgManager = cfg.Build.PkgManager
	}
	if !cmd.Flags().Changed("output") && cfg.Build.Output != "" {
		buildOutput = cfg.Build.Output
	}

	// Check if frontend directory exists
	if _, err := os.Stat(buildFrontendFolder); os.IsNotExist(err) {
		return fmt.Errorf("frontend directory '%s' not found. Make sure you're in a project root with a frontend directory", buildFrontendFolder)
	}

	configPath := filepath.Join(buildFrontendFolder, "vite.config.ts")
	// Check for Vite configuration
	hasViteConfig := false
	if _, err := os.Stat(configPath); err == nil {
		hasViteConfig = true
	}

	if !hasViteConfig {
		return fmt.Errorf("vite configuration not found in %s directory. Please ensure vite.config.ts exists", buildFrontendFolder)
	}

	// Check for package.json
	packageJsonPath := filepath.Join(buildFrontendFolder, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in %s directory", buildFrontendFolder)
	}

	// Detect or use specified package manager
	var packageManager string

	if buildPkgManager != "" {
		// Validate user-specified package manager
		if buildPkgManager != "npm" && buildPkgManager != "yarn" && buildPkgManager != "pnpm" {
			return fmt.Errorf("invalid package manager '%s'. Supported: npm, yarn, pnpm", buildPkgManager)
		}

		// Verify the specified package manager is installed
		if _, err := exec.LookPath(buildPkgManager); err != nil {
			return fmt.Errorf("specified package manager '%s' is not installed", buildPkgManager)
		}

		packageManager = buildPkgManager
	} else {
		// Auto-detect package manager from lock files
		packageManager, err = detectPackageManager(buildFrontendFolder)
		if err != nil {
			return fmt.Errorf("failed to detect package manager: %w", err)
		}
	}

	cmd.Printf("Installing dependencies with %s...\n", packageManager)

	// Install dependencies
	if err := installDependencies(buildFrontendFolder, packageManager); err != nil {
		return fmt.Errorf("dependency installation failed: %w", err)
	}

	cmd.Printf("Building frontend with %s...\n", packageManager)

	// Build frontend
	if err := buildFrontend(buildFrontendFolder, packageManager); err != nil {
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

	// Save configuration if any non-default values were used
	if config.HasNonDefaultBuildValues(buildOutput, buildFrontendFolder, buildMainFolder, buildPkgManager) {
		cfg.Build.Output = buildOutput
		cfg.Build.FrontendFolder = buildFrontendFolder
		cfg.Build.MainFolder = buildMainFolder
		cfg.Build.PkgManager = packageManager

		if err := config.SaveConfig(cfg); err != nil {
			cmd.Printf("Warning: failed to save configuration: %v\n", err)
		} else {
			cmd.Println("Configuration saved to .kirin.toml")
		}
	}

	cmd.Println("✨ Full-stack build completed successfully!")
	return nil
}

// detectPackageManager detects the package manager based on lock files
func detectPackageManager(frontendDir string) (string, error) {
	// Check for lock files in order of preference
	lockFiles := map[string]string{
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
		"package-lock.json": "npm",
	}

	for lockFile, manager := range lockFiles {
		lockPath := filepath.Join(frontendDir, lockFile)
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

// installDependencies installs frontend dependencies
func installDependencies(frontendDir, packageManager string) error {
	var installCmd *exec.Cmd

	switch packageManager {
	case "npm":
		installCmd = exec.Command("npm", "install")
	case "yarn":
		installCmd = exec.Command("yarn", "install")
	case "pnpm":
		installCmd = exec.Command("pnpm", "install")
	default:
		return fmt.Errorf("unsupported package manager: %s", packageManager)
	}

	installCmd.Dir = frontendDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	return installCmd.Run()
}

// buildFrontend runs the frontend build process
func buildFrontend(frontendDir, packageManager string) error {
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

	buildCmd.Dir = frontendDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	return buildCmd.Run()
}

// findMainGo finds the main.go file in main directory or nested subdirectories
func findMainGo() (string, error) {
	if _, err := os.Stat(buildMainFolder); os.IsNotExist(err) {
		return "", fmt.Errorf("%s directory not found", buildMainFolder)
	}

	var mainGoPath string

	// Check for main.go directly in main folder first
	directMainGo := filepath.Join(buildMainFolder, "main.go")
	if _, err := os.Stat(directMainGo); err == nil {
		return directMainGo, nil
	}

	// Look for main.go in nested subdirectories
	err := filepath.Walk(buildMainFolder, func(path string, info os.FileInfo, err error) error {
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
		return "", fmt.Errorf("main.go not found in %s directory or its subdirectories", buildMainFolder)
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
	if buildOutput != "" {
		outputName = buildOutput
	} else {
		// If it's mainFolder/main.go, use current directory name
		if filepath.Base(dir) == buildMainFolder {
			cwd, _ := os.Getwd()
			outputName = filepath.Base(cwd)
		} else {
			// If it's mainFolder/app/main.go, use the app name
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

# With custom directories
kirin build --frontend-folder ui --main-folder app

# With specific package manager
kirin build --pkg-manager pnpm

# With custom output name
kirin build --output myapp
`

func init() {
	buildCmd.Flags().StringVar(&buildOutput, "output", "", "Output binary name (default: derived from directory)")
	buildCmd.Flags().StringVar(&buildFrontendFolder, "frontend-folder", "web", "Frontend directory name")
	buildCmd.Flags().StringVar(&buildMainFolder, "main-folder", "cmd", "Main directory name")
	buildCmd.Flags().StringVar(&buildPkgManager, "pkg-manager", "", "Package manager to use (npm, yarn, pnpm). Auto-detected if not specified")
}
