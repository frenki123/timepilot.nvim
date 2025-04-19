package main

import (
	"os"

	"github.com/frenki123/timepilot.nvim/timepilot/internal/cmd"
	_ "modernc.org/sqlite"
)

func main() {
	cmd.Run(os.Args[1:])
}
