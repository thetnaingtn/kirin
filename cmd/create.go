package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thetnaingtn/kirin/internal/kirin"
)

var (
	frontend string
)

func validateFrontend(frontend string) bool {
	normalizeFrontend := strings.ToLower(frontend)
	return slices.Contains([]string{"react", "vue", "svelte"}, normalizeFrontend)
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

func newRunE(cmd *cobra.Command, args []string) error {
	cmd.Println("Scaffolding a new full-stack gRPC application...")
	start := time.Now()
	appName := args[0]
	modName := appName
	if len(args) > 1 {
		modName = args[1]
	}

	if !validateFrontend(frontend) {
		return fmt.Errorf("unsupported frontend framework: %s (supported: React, Vue, Svelte)", frontend)
	}

	frontend = strings.ToLower(frontend)

	if err := kirin.CreateProject(appName, modName, frontend); err != nil {
		return err
	}

	duration := time.Since(start)

	cmd.Printf(createSuccessMessage, appName, modName, kirin.FormatTime(duration))

	return nil
}

var (
	createExamples = `
kirin create myapp
Generate a new full-stack gRPC application named "myapp" with default frontend (React)

kirin create myapp your.own/module/name
Generate a new full-stack gRPC application with a custom module name

kirin create myapp --frontend vue
Generate a new full-stack gRPC application with the provided frontend framework (supported: React, Vue, Svelte)
`

	createSuccessMessage = `
Create new application project named: %s (module %s)

✨  Done in %s.
`
)
