package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	frontend string
	cloneUrl = "https://github.com/thetnaingtn/kirin"
)

func init() {
	createCmd.Flags().StringVarP(&frontend, "frontend", "f", "react", "Specify the frontend framework (supported: react, vue, svelte)")
}

var createCmd = &cobra.Command{
	Use:     "create <app-name> [module name]",
	Example: createExamples,
	Aliases: []string{"c"},
	Short:   "Create a new full-stack gRPC application",
	RunE:    newRunE,
	Args:    cobra.MinimumNArgs(1),
}

func newRunE(cmd *cobra.Command, args []string) error {
	cmd.Println("Scaffolding a new full-stack gRPC application...")
	start := time.Now()
	appName := args[0]
	modName := appName
	if len(args) > 1 {
		modName = args[1]
	}

	wd, _ := os.Getwd()
	projectPath := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, appName)

	if err := os.Mkdir(projectPath, 0750); err != nil {
		return err
	}

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	c := exec.Command(git, "clone", "-b", fmt.Sprintf("frontend/%s", frontend), cloneUrl, projectPath)

	if err := c.Run(); err != nil {
		return err
	}

	if err := replace(projectPath, "go.mod", "kirin", modName); err != nil {
		return err
	}

	if err := replace(projectPath, "*.go", "kirin", modName); err != nil {
		return err
	}

	cmd.Printf(createSuccessMessage, projectPath, modName, formatTime(time.Since(start)))

	return nil
}

var (
	createExamples = `
wano create myapp
Generate a new full-stack gRPC application named "myapp" with default frontend (React)

wano create myapp your.own/module/name
Generate a new full-stack gRPC application with a custom module name

wano create myapp --frontend vue
Generate a new full-stack gRPC application with the provided frontend framework (supported: react, vue, svelte)
`

	createSuccessMessage = `
Create new application project in %s (module %s)

✨  Done in %s.
`
)
