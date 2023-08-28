package commands

import (
	"flag"

	repository "ghorgctl/internal/github"
)

func NewPullCommand() *PullCommand {
	lc := &PullCommand{
		fs: flag.NewFlagSet("pull", flag.ContinueOnError),
	}

	return lc
}

type PullCommand struct {
	fs *flag.FlagSet
}

func (c *PullCommand) Name() string {
	return c.fs.Name()
}

func (c *PullCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *PullCommand) Run() error {
	err := repository.Pull()
	if err != nil {
		return err
	}

	return nil
}
