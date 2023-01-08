package routine

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/jmkng/onyx/config"
)

const (
	// Default permission set for new files.
	DefFilePerm fs.FileMode = 0644
	// Default permission set for new directories.
	DefDirPerm fs.FileMode = 0755
)

// WdOrPanic will return the working directory, or panic if os.Getwd() fails.
func WdOrPanic() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to access working directory")
	}

	return wd
}

// Setup will assert that the given path is a valid Onyx project by searching for a
// recognized configuration file. If one is found, it will be read and stored in
// config.State. An error is returned if no configuration file exists or the file
// is malformed.
func Setup(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access directory: %v", path)
	}

	configPath, err := config.SearchConf(path)
	if err != nil {
		return fmt.Errorf("configuration file onyx.[yaml|yml|json] was not found in directory: %v", path)
	}

	err = config.SetState(configPath)
	if err != nil {
		return fmt.Errorf("configuration file `%v` is malformed", configPath)
	}

	return nil
}

// IsVerbose will return true if the given bool is true or config.State.Verbose is true.
func IsVerbose(arg bool) bool {
	return arg || config.State.Verbose
}
