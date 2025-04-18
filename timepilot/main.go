package main

import (
	"github.com/frenki123/timepilot.nvim/timepilot/internal/tprpc"
)

func main() {
    srv := tprpc.NewServer()
    srv.ListenAndServe()
}
