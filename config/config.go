package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmkng/onyx/convert/json"
	"github.com/jmkng/onyx/convert/yaml"
)

const (
	YamlName    = "onyx.yaml"
	YamlAltName = "onyx.yml"
	JsonName    = "onyx.json"
)

// State contains the project configuration during runtime.
var State Options

// Names contains all recognized project configuration file names.
var Names = []string{
	YamlName,
	YamlAltName,
	JsonName,
}

// Options describes all of the options that might be found in a
// recognized configuration file.
type Options struct {
	Domains  []string `json:"domains" yaml:"domains"`
	Preserve []string `json:"preserve" yaml:"preserve"`
	Output   string   `json:"output" yaml:"output"`
}

// Read will attempt to unmarshal the configuration file at the given path.
// If a configuration is successfully unmarshaled and stored in config.State,
// error is nil.
func Read(path string) error {
	ext := filepath.Ext(path)

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	switch ext {
	case ".json":
		err = json.Unmarshal(bytes, &State)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(bytes, &State)
	default:
		panic("attempted to unmarshal unrecognized configuration file type")
	}

	if err != nil {
		return err
	}

	return nil
}

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

// CreateTemp will create a temporary directory on disk. If the file parameter
// is not an empty string, a file will also be created within that temporary directory.
// The returned string points to the new temporary directory, not the file.
func CreateTemp(t testing.TB, file string) (string, error) {
	t.Helper()
	error := fmt.Errorf("failed to create temporary directory for test: %v", t.Name())

	dir := t.TempDir()

	if file == "" {
		return dir, nil
	}

	fullPath := filepath.Join(dir, file)

	newFile, err := os.Create(fullPath)
	if err != nil {
		return "", error
	}

	defer func() {
		err := newFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	return dir, nil
}

// SearchConf is a wrapper around config.AnyExist and config.Names. It searches the
// given directory for a recognized project configuration file, and returns a path
// to that file. If none are found, an error is returned.
func SearchConf(dir string) (string, error) {
	result, err := AnyExist(dir, Names[:])
	if err != nil {
		return "", err
	}

	return result, nil
}
