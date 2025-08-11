package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

const (
	longDescription = `⛩️  wano is a CLI tool that helps scaffolding full-stack gRPC applications with end-to-end type safety.`
)

var rootCmd = &cobra.Command{
	Use:  "wano",
	Long: longDescription,
	RunE: rootRunE,
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func rootRunE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
