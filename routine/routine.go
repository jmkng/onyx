package routine

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const (
	DefFilePerm fs.FileMode = 0644
	DefDirPerm  fs.FileMode = 0755
)

// AnyExist searches a directory and returns a path to the first of any of
// the given files. If none are found, an error is returned.
func AnyExist(dir string, files []string) (string, error) {
	var (
		resultPath  string
		resultError error
	)

	for _, v := range files {
		testPath := filepath.Join(dir, v)

		info, err := os.Stat(testPath)
		if os.IsNotExist(err) {
			resultError = err
			continue
		} else {
			resultPath = filepath.Join(dir, info.Name())
			resultError = nil
			break
		}

	}

	return resultPath, resultError
}

// createTemp will create a temporary directory for testing. If the file parameter is not an empty string,
// a file will also be created within that temporary directory. The returned string points to the new temporary
// directory, not the file.
func CreateTemp(t testing.TB, file string) (string, error) {
	t.Helper()
	error := errors.New("could not create directory")

	dir := t.TempDir()

	if file == "" {
		return dir, nil
	}

	fullPath := filepath.Join(dir, file)

	newFile, err := os.Create(fullPath)
	if err != nil {
		return "", error
	}

	newFile.Close()

	return dir, nil
}

// WdOrPanic will return the working directory. If the working directory
// is not available, the program will panic.
func WdOrPanic() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to access working directory")
	}

	return wd
}
