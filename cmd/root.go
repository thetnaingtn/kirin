package cmd

import (
	"context"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

const (
	longDescription = `🐉  kirin is a CLI tool that helps scaffolding full-stack gRPC applications with end-to-end type safety.`
)

var rootCmd = &cobra.Command{
	Use:  "kirin",
	Long: longDescription,
	RunE: rootRunE,
}

func init() {
	rootCmd.AddCommand(doctorCmd, createCmd)
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions()); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func rootRunE(cmd *cobra.Command, args []string) error {
	prompt := NewPrompt()

	p := tea.NewProgram(prompt, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
