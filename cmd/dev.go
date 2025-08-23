package cmd

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/air-verse/air/runner"
	"github.com/spf13/cobra"
)

var cfgPath string

var devCmd = &cobra.Command{
	Use:  "dev",
	RunE: devRunE,
}

func init() {
	devCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path to Air config")
}

func devRunE(cmd *cobra.Command, args []string) error {
	c := flag.CommandLine

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	cmdArgs := runner.ParseConfigFlag(c)
	cfg, _ := runner.InitConfig(cfgPath, cmdArgs)

	engine, err := runner.NewEngineWithConfig(cfg, false)
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		engine.Stop()
	}()

	engine.Run()

	return nil
}
