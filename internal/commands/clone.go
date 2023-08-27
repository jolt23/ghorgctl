package commands

import (
	"flag"

	repository "ghorgctl/internal/github"
)

func NewCloneCommand() *CloneCommand {
	lc := &CloneCommand{
		fs: flag.NewFlagSet("clone", flag.ContinueOnError),
	}

	lc.fs.StringVar(&lc.org, "org", "", "name of the organization or user that owns the repositories")
	lc.fs.StringVar(&lc.prefix, "prefix", "", "repository prefix to match against")

	return lc
}

type CloneCommand struct {
	fs *flag.FlagSet

	org    string
	prefix string
}

func (c *CloneCommand) Name() string {
	return c.fs.Name()
}

func (c *CloneCommand) Org() string {
	return c.org
}

func (c *CloneCommand) Prefix() string {
	return c.prefix
}

func (c *CloneCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *CloneCommand) Run() error {
	err := repository.Clone(c.org, c.prefix)
	if err != nil {
		return err
	}

	return nil
}
