package main

import (
	"github.com/frenki123/timepilot.nvim/timepilot/internal/timepilot"
    _ "modernc.org/sqlite"
)

func main() {
    app := timepilot.NewApp()
    app.Run()
}
