package routines

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"
)

const (
	DefFilePerm fs.FileMode = 0644
	DefDirPerm  fs.FileMode = 0755
)

// Routine describes a type with Execute() and Setup() functions. Setup is
// optionally used to configure the environment before the main logic in Execute.
type Routine interface {
	Execute() error
}

// AnyExist searches a directory and returns the first of any of the given
// files that are present. If none are found, an error is returned.
func AnyExist(dir string, files []string) (string, error) {
	var (
		resultPath  string
		resultError error
	)

	for _, v := range files {
		testPath := path.Join(dir, v)

		info, err := os.Stat(testPath)
		if os.IsNotExist(err) {
			resultError = err
			continue
		} else {
			resultPath = path.Join(dir, info.Name())
			resultError = nil
			break
		}

	}

	return resultPath, resultError
}

// EvaluatePath resolves a zero-value string to the current directory.
// If a non-zero string is given, the same value is returned. If the
// current directory is not accessible, an error is returned.
func EvaluatePath(path string) (string, error) {
	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.New("unable to access current directory")
		}

		path = wd
	}

	return path, nil
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

	fullPath := path.Join(dir, file)

	_, err := os.Create(fullPath)
	if err != nil {
		return "", error
	}

	return dir, nil
}
