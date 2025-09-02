package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"github.com/thetnaingtn/kirin/internal/ui"
)

const (
	longDescription = `🐉  kirin is a CLI tool that scaffolds full-stack gRPC applications with end-to-end type safety.`
)

var rootCmd = &cobra.Command{
	Use:   "kirin",
	Short: "Scaffold full-stack gRPC applications in an interactive way",
	Long:  longDescription,
	RunE:  rootRunE,
}

func init() {
	rootCmd.AddCommand(doctorCmd, createCmd, devCmd, generateCmd, buildCmd)
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions(), fang.WithNotifySignal(os.Interrupt, os.Kill)); err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(ver, d string) {
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s Build at: %s", ver, d))
}

func rootRunE(cmd *cobra.Command, args []string) error {
	prompt := ui.NewPrompt()

	p := tea.NewProgram(prompt, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

var latestVersionRegexp = regexp.MustCompile(`"name":\s*?"v(.*?)"`)

func latestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/thetnaingtn/kirin/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}
