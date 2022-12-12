package routines

import (
	"errors"
	"os"
	"path"

	"github.com/jmkng/onyx/track"
)

type New struct {
	Path string
}

var errAcc = errors.New("cannot access directory (missing permissions?)")

func (n New) Execute() error {
	info, stErr := os.Stat(n.Path)
	if stErr != nil && errors.Is(stErr, os.ErrNotExist) {
		mkErr := os.Mkdir(n.Path, DefDirPerm)
		if mkErr != nil {
			return errors.New("path does not exist, and unable to create (missing permission?)")
		}
	} else if stErr != nil {
		return errAcc
	}

	if info != nil && !info.IsDir() {
		return errors.New("path leads to file, expected directory")
	}

	path := path.Join(n.Path, "config.yaml")

	if _, err := os.Stat(path); err == nil {
		return errors.New("config file already exists in this directory")
	}

	wrErr := os.WriteFile(path, []byte(yamlConf), DefFilePerm)
	if wrErr != nil {
		return errAcc
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
---`
