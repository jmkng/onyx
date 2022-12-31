package routine

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/jmkng/onyx/config"
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

// Setup will assert that the given path is a valid Onyx project, or return an error
// that describes what is missing.
func Setup(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access directory: %v", path)
	}

	configPath, err := config.SearchConf(path)
	if err != nil {
		return fmt.Errorf("configuration file onyx.[yaml|yml|json] was not found in directory: %v", path)
	}

	err = config.Read(configPath)
	if err != nil {
		return fmt.Errorf("configuration file `%v`  is malformed", configPath)
	}

	return nil
}
