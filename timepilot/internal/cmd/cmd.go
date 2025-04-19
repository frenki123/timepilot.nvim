package cmd

import (
	"fmt"

	"github.com/frenki123/timepilot.nvim/timepilot/internal/timepilot"
)

var version = "0.0.1"

func Run(args []string) {
	if len(args) == 0 {
		timepilot.NewApp().Run()
		return
	}
	if len(args) == 1 {
		cmd := args[0]
		if cmd == "serve" {
			timepilot.NewApp().Run()
			return
		}
		if cmd == "version" {
			fmt.Printf("Timepilot v%s\n", version)
			return
		}
		if cmd == "help" {
			showHelp("")
			return
		}
		showHelp(fmt.Sprintf("Unknown argument '%s'", cmd))
		return
	}
	showHelp("To many arguments")
}

func showHelp(errMsg string) {
	if errMsg != "" {
		fmt.Println("ERROR:", errMsg)
	}
	fmt.Println("To run the command 'timepilot' or 'timepilot serve' to start the deamon.")
	fmt.Println("To check the version run 'timepilot version'")
	fmt.Println("To see this msg as help run 'timepilot help'")
}
