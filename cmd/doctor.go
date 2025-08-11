package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/muesli/termenv"

	"github.com/spf13/cobra"
)

var prequisites = []string{
	"protoc",
	"buf",
	"git",
}

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"d", "doc"},
	Short:   "To check whether pre-requisites binaries/libraries are installed.",
	RunE:    doctorRunE,
}

func doctorRunE(cmd *cobra.Command, args []string) error {
	missingPrerequisites := []string{}

	for _, p := range prequisites {
		if _, err := exec.LookPath(p); err != nil {
			missingPrerequisites = append(missingPrerequisites, p)
		}
	}

	if len(missingPrerequisites) > 0 {
		warning := fmt.Sprintf("⚠️  These pre-requisites are missing: %s.", strings.Join(missingPrerequisites, ", "))
		cmd.Println(termenv.String(warning).Foreground(termenv.ANSIBrightYellow))
	} else {
		msg := "🚀 All pre-requisites are installed."
		cmd.Println(termenv.String(msg).Foreground(termenv.ANSICyan))
	}

	return nil
}
