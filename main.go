package main

import (
	"fmt"
	"os"

	commands "ghorgctl/internal/commands"
)

func main() {
	cmds := []commands.Runner{
		commands.NewCloneCommand(),
		commands.NewListCommand(),
		commands.NewPullCommand(),
	}

	if len(os.Args) <= 1 {
		help(cmds)
	}

	cmd := os.Args[1]

	if cmd == "help" {
		help(cmds)
	}

	for _, c := range cmds {
		if c.Name() == cmd {
			c.Init(os.Args[2:])
			err := c.Run()
			if err != nil {
				fmt.Printf("error executing command: %s\n", c)
				fmt.Println(err.Error())
			}

			os.Exit(0)
		}
	}

	fmt.Println("unknown command. use help for supported commands.")
	os.Exit(1)
}

func help(cmds []commands.Runner) {
	help := commands.NewHelpCommand(cmds)
	help.Run()

	os.Exit(0)
}
