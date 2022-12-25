package routine

import (
	"io/fs"
	"os"
)

const (
	DefFilePerm fs.FileMode = 0644
	DefDirPerm  fs.FileMode = 0755
)

// WdOrPanic will return the working directory. If the working directory
// is not available, the program will panic.
func WdOrPanic() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to access working directory")
	}

	return wd
}
