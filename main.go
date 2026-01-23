package main

import (
	"github.com/thetnaingtn/kirin/cmd"
)

var (
	version = "v0.1.0"
)

func main() {
	cmd.SetVersionInfo(version)
	cmd.Execute()
}
