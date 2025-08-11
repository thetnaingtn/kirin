package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	frontend string
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

	cloneUrl := "https://github.com/thetnaingtn/kirin"

	c := exec.Command(git, "clone", cloneUrl, projectPath)

	if err := c.Run(); err != nil {
		return err
	}

	if err := replace(projectPath, "go.mod", "kirin", modName); err != nil {
		return err
	}

	if err := replace(projectPath, "*.go", "kirin", modName); err != nil {
		return err
	}

	cmd.Println(createSuccessMessage)

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
Your new full-stack gRPC application has been created successfully!
`
)

func replace(path, pattern, old, new string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return replaceWalkFn(path, info, pattern, []byte(old), []byte(new))
	})
}

func replaceWalkFn(path string, info os.FileInfo, pattern string, old, new []byte) (err error) {
	var matched bool
	if matched, err = filepath.Match(pattern, info.Name()); err != nil {
		return
	}

	if matched {
		cleanedPath := filepath.Clean(path)

		var oldContent []byte
		if oldContent, err = os.ReadFile(cleanedPath); err != nil {
			return
		}

		if err = os.WriteFile(cleanedPath, bytes.Replace(oldContent, old, new, -1), 0); err != nil {
			return
		}
	}

	return
}
