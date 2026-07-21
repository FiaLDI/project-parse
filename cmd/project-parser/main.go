package main

import (
	"os"

	"github.com/FiaLDI/project-parse/cmd/project-parser/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
