package commands

import (
	"flag"
	"fmt"
)

func NewHelpCommand(commands []Runner) *HelpCommand {
	hc := &HelpCommand{
		fs:   flag.NewFlagSet("help", flag.ContinueOnError),
		cmds: commands,
	}

	return hc

}

type HelpCommand struct {
	fs   *flag.FlagSet
	cmds []Runner
}

func (hc *HelpCommand) Name() string {
	return hc.fs.Name()
}

func (hc *HelpCommand) Init(args []string) error {
	return nil
}

func (hc *HelpCommand) Run() error {
	fmt.Println("ghorbot is a tool to help manage large about of repositories in github.")
	fmt.Println("")
	fmt.Println("commands:")
	for _, c := range hc.cmds {
		fmt.Printf("%s\n", c.Name())
	}
	fmt.Println("")
	fmt.Println("global options:")
	fmt.Println("-org: github organization")

	return nil
}
