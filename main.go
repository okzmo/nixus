package main

import (
	"fmt"
	"os"

	"github.com/okzmo/nixus/cmd/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Unknown command. Usage: nixus <command>")
		os.Exit(1)
	}

	commands.Execute(os.Args[1])
}
