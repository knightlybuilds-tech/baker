package main

import (
	"fmt"
	"os"

	"github.com/roguenet/baker/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "baker:", err)
		os.Exit(1)
	}
}
