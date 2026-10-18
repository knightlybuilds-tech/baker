package main

import (
	"fmt"
	"os"

	"github.com/knightlybuilds-tech/baker/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "baker:", err)
		os.Exit(1)
	}
}
