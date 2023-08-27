package commands

import (
	"flag"

	repository "ghorgctl/internal/github"
)

func NewListCommand() *ListCommand {
	lc := &ListCommand{
		fs: flag.NewFlagSet("list", flag.ContinueOnError),
	}

	lc.fs.StringVar(&lc.org, "org", "", "name of the organization or user that owns the repositories")
	lc.fs.StringVar(&lc.prefix, "prefix", "", "repository prefix to match against")

	return lc
}

type ListCommand struct {
	fs *flag.FlagSet

	org    string
	prefix string
}

func (l *ListCommand) Name() string {
	return l.fs.Name()
}

func (l *ListCommand) Org() string {
	return l.org
}

func (l *ListCommand) Prefix() string {
	return l.prefix
}

func (l *ListCommand) Init(args []string) error {
	return l.fs.Parse(args)
}

func (l *ListCommand) Run() error {
	err := repository.List(l.org, l.prefix)
	if err != nil {
		return err
	}

	return nil
}
