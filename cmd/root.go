package cmd

import (
	"context"
	"os"

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
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions(), fang.WithNotifySignal(os.Interrupt, os.Kill)); err != nil {
		os.Exit(1)
	}
}

func rootRunE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
