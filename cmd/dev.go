package cmd

import (
	goflag "flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/air-verse/air/runner"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var cfgPath string

var devCmd = &cobra.Command{
	Use:                "dev",
	Short:              "Run the development server with live reloading",
	RunE:               devRunE,
	DisableFlagParsing: true,
}

func devRunE(cmd *cobra.Command, _ []string) error {
	flagSet := goflag.CommandLine

	flag.CommandLine.AddGoFlagSet(flagSet)
	cmd.Flags().AddFlagSet(flag.CommandLine)

	argsMap := runner.ParseConfigFlag(flagSet)

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		return err
	}

	cfg, err := runner.InitConfig(cfgPath, argsMap)
	if err != nil {
		return err
	}

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
