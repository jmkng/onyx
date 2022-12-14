package routines

import (
	"errors"
	"flag"
	"os"
	"path"

	"github.com/jmkng/onyx/track"
)

const (
	errAcc = "cannot access directory (missing permissions?)"
)

func NewCreate() *Create {
	c := &Create{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.path, "path", "", "Path to the desired location of the new project.")

	return c
}

type Create struct {
	fs   *flag.FlagSet
	path string
}

func (c *Create) Name() string {
	return c.fs.Name()
}

func (c *Create) Parse(args []string) error {
	return c.fs.Parse(args)
}

func (c *Create) Execute() error {
	info, stErr := os.Stat(c.path)
	if stErr != nil && errors.Is(stErr, os.ErrNotExist) && c.path != "" {
		mkErr := os.Mkdir(c.path, DefDirPerm)
		if mkErr != nil {
			return errors.New("path does not exist, and unable to create (missing permission?)")
		}
	} else if stErr != nil && c.path != "" {
		return errors.New(errAcc)
	}

	if info != nil && !info.IsDir() {
		return errors.New("path leads to file, expected directory")
	}

	path := path.Join(c.path, "config.yaml")

	if _, err := os.Stat(path); err == nil {
		return errors.New("config file already exists in this directory")
	}

	wrErr := os.WriteFile(path, []byte(yamlConf), DefFilePerm)
	if wrErr != nil {
		return errors.New(errAcc)
	}

	track.Log(
		track.Event("Wrote file: " + path),
	)

	return nil
}

var yamlConf = `---
source: some/place
destination: some/place/build/static
include: other/place
exclude: other/place/private 
preserve: ["static"]
---
`
