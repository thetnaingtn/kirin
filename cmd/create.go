package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	frontend           string
	cloneUrl           = "https://github.com/thetnaingtn/boilerplate"
	supportedFrontends = []string{"React", "Vue", "Svelte"}
)

func normalizeFrontend(frontend string) string {
	switch frontend {
	case "react", "React":
		return "react"
	case "vue", "Vue":
		return "vue"
	case "svelte", "Svelte":
		return "svelte"
	default:
		return ""
	}
}

func init() {
	createCmd.Flags().StringVarP(&frontend, "frontend", "f", "react", "Specify the frontend framework (supported: React, Vue, Svelte)")
}

var createCmd = &cobra.Command{
	Use:     "create <app-name> [module name]",
	Example: createExamples,
	Aliases: []string{"c"},
	Short:   "Create a new full-stack gRPC application",
	RunE:    newRunE,
	Args:    cobra.MinimumNArgs(1),
}

func newRunE(cmd *cobra.Command, args []string) (err error) {
	cmd.Println("Scaffolding a new full-stack gRPC application...")
	start := time.Now()
	appName := args[0]
	modName := appName
	if len(args) > 1 {
		modName = args[1]
	}

	wd, _ := os.Getwd()
	projectPath := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, appName)

	if err = os.Mkdir(projectPath, 0750); err != nil {
		return
	}

	defer func() {
		if err != nil {
			os.RemoveAll(projectPath)
		}
	}()

	var git string
	git, err = exec.LookPath("git")

	if err != nil {
		return fmt.Errorf("git is not installed or not found in PATH")
	}

	frontend = normalizeFrontend(frontend)

	if frontend == "" {
		return fmt.Errorf("unsupported frontend framework: %s (supported: %v)", frontend, supportedFrontends)
	}

	c := exec.Command(git, "clone", "-b", fmt.Sprintf("frontend/%s", frontend), cloneUrl, projectPath)

	if err := c.Run(); err != nil {
		return err
	}

	if err = replace(projectPath, "go.mod", "bolierplate", modName); err != nil {
		return
	}

	if err = replace(projectPath, "*.go", "bolierplate", modName); err != nil {
		return
	}

	cmd.Printf(createSuccessMessage, projectPath, modName, formatTime(time.Since(start)))

	return nil
}

var (
	createExamples = `
kirin create myapp
Generate a new full-stack gRPC application named "myapp" with default frontend (React)

kirin create myapp your.own/module/name
Generate a new full-stack gRPC application with a custom module name

kirin create myapp --frontend vue
Generate a new full-stack gRPC application with the provided frontend framework (supported: react, vue, svelte)
`

	createSuccessMessage = `
Create new application project in %s (module %s)

✨  Done in %s.
`
)
