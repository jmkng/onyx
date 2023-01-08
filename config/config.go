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
	// Expected configuration file name with no extension.
	BaseName = "onyx"
	// Expected configuration file name with long YAML extension.
	YamlLongName = BaseName + ".yaml"
	// Expected configuration file name with short YAML extension.
	YamlShortName = BaseName + ".yml"
	// Expected configuration file name with json extension.
	JsonName = BaseName + ".json"
	// Expected date format.
	DateFmt = "2006-01-02"
)

// State contains an instance of Options derived from a project configuration
// file during runtime.
var State Options

// Names contains all recognized project configuration file names.
var Names = []string{
	YamlLongName,
	YamlShortName,
	JsonName,
}

// Options describes all of the values that might be found in a project's
// configuration file.
type Options struct {
	// Output controls the location that the resulting files are placed.
	// This directory is the
	Output string `json:"output" yaml:"output"`
	// Verbose determines how much information is logged.
	Verbose bool `json:"verbose" yaml:"verbose"`
	// These aren't supported yet, so better comment them out for now.
	// Domains  []string `json:"domains" yaml:"domains"`
	// Preserve []string `json:"preserve" yaml:"preserve"`
}

// SetState will read the file at the given path and unmarshal the contents into
// config.State. An error is returned if the file is malformed or cannot be read.
func SetState(path string) error {
	ext := filepath.Ext(path)

	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %v", filepath.Base(path))
	}

	switch ext {
	case ".json":
		err = json.Unmarshal(bytes, &State)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(bytes, &State)
	default:
		return fmt.Errorf("attempted to unmarshal unrecognized configuration file type: %v", filepath.Base(path))
	}

	if err != nil {
		return err
	}

	return nil
}

// AnyExist searches a directory and returns a path to the first of any of
// the given files. An error is returned if none are found.
func AnyExist(dir string, files []string) (string, error) {
	var (
		resultPath  string
		resultError error
	)

	for _, v := range files {
		testPath := filepath.Join(dir, v)

		info, err := os.Stat(testPath)
		if err != nil && os.IsNotExist(err) {
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

	dir := t.TempDir()

	if file == "" {
		return dir, nil
	}

	fullPath := filepath.Join(dir, file)

	parent := filepath.Dir(fullPath)

	err := os.MkdirAll(parent, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create parent directory for temporary file: %v\n%v", file, err)
	}

	newFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for test: %v\n%v", t.Name(), err)
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
