package main

import (
	"github.com/daigo/gsecutil/cmd"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {
	cmd.Execute(Version)
}
