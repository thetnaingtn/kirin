package main

import (
	"github.com/thetnaingtn/kirin/cmd"
)

var (
	version = "v0.1.0"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, date)
	cmd.Execute()
}
