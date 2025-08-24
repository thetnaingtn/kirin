package cmd

import (
	goflag "flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/air-verse/air/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cfgPath string
	flagSet *goflag.FlagSet
	argsMap map[string]runner.TomlInfo
)

func init() {
	flagSet = goflag.NewFlagSet("air", goflag.ContinueOnError)
	if len(os.Args) > 2 {
		flagSet.Parse(os.Args[2:])
	}

	pf := pflag.NewFlagSet("air", pflag.ContinueOnError)

	argsMap = runner.ParseConfigFlag(flagSet)

	flagSet.VisitAll(func(f *goflag.Flag) {
		name := f.Name
		if name == "help" || name == "h" {
			return // keep Cobra’s help, don’t import Air’s
		}
		pf.AddGoFlag(f)
	})

	devCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path air to config file")
	devCmd.Flags().AddFlagSet(pf)
}

var devCmd = &cobra.Command{
	Use:     "dev",
	Short:   "Run the development server with live reloading",
	Example: devExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

	},
}

var (
	devExample = `
# Start the development server with live reloading
kirin dev
# kirin dev use air under the hood so any air flags can be passed here
kirin dev --color.build=red --build.cmd="go build -o ./mytmp ." --tmp_dir=mytmp --build.bin=./mytmp/main
	`
)
