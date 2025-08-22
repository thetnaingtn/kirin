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
	"protoc-gen-go",
	"protoc-gen-go-grpc",
}

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"d", "doc"},
	Example: doctorExamples,
	Short:   "Checks if required dependencies are installed on your system.",
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
		warning := fmt.Sprintf(preRequisitesMissing, mayBePlural(len(missingPrerequisites)), strings.Join(missingPrerequisites, ", "))
		cmd.Print(termenv.String(warning).Foreground(termenv.ANSIBrightYellow))
	} else {
		msg := "🚀 All pre-requisites are installed."
		cmd.Println(termenv.String(msg).Foreground(termenv.ANSICyan))
	}

	return nil
}

func mayBePlural(count int) string {
	if count == 1 {
		return "is"
	}
	return "are"
}

var (
	preRequisitesMissing = `
⚠️  Pre-requisites %s missing: %s.

Install via your package manager or download from their official websites.
Here are some links to get you started:
- protoc: https://protobuf.dev/installation/
- buf: https://buf.build/docs/cli/installation/
- git: https://git-scm.com/downloads
- protoc-gen-go: https://protobuf.dev/getting-started/gotutorial/#compiling-protocol-buffers
	`

	doctorExamples = `
# Check if pre-requisites are installed
kirin doctor 

# Or use aliases
kirin doc
kirin d
`
)
