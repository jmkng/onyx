package routines

import (
	"flag"
	"fmt"
)

func NewBuild() *Build {
	b := &Build{
		fs: flag.NewFlagSet("build", flag.ContinueOnError),
	}

	b.fs.StringVar(&b.path, "path", "", "Path to the project being built.")

	return b
}

type Build struct {
	fs   *flag.FlagSet
	path string
}

func (b *Build) Name() string {
	return b.fs.Name()
}

func (b *Build) Parse(args []string) error {
	return b.fs.Parse(args)
}

func (b *Build) Execute() error {
	fmt.Println("Executed build routine")
	return nil
}
