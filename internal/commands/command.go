package commands

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}
