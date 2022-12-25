package routine

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/jmkng/onyx/config"
)

func TestCreateExecute(t *testing.T) {
	t.Run("config file is created when path argument does not already contain a config", func(t *testing.T) {
		dir, err := config.CreateTemp(t, "")
		if err != nil {
			t.Fail()
		}

		routine := Create{
			path: dir,
		}

		err = routine.Execute()
		if err != nil {
			t.Fail()
		}
	})

	t.Run("directory will be created if path argument is a non-existent directory", func(t *testing.T) {
		dir, err := config.CreateTemp(t, "")
		if err != nil {
			t.Fail()
		}

		path := path.Join(dir, "nonExistentDir")

		routine := Create{
			path: path,
		}

		err = routine.Execute()
		if err != nil {
			t.Fail()
		}
	})

	t.Run("existing config file will not be overridden", func(t *testing.T) {
		dir, err := config.CreateTemp(t, "onyx.yaml")
		if err != nil {
			t.Fail()
		}

		routine := Create{
			path: dir,
		}

		err = routine.Execute()
		if err == nil {
			t.Fail()
		}
	})

	t.Run("path to a file cannot be used as path argument", func(t *testing.T) {
		dir, err := config.CreateTemp(t, "onyx.yaml")
		if err != nil {
			t.Fail()
		}

		path := path.Join(dir, "onyx.yaml")
		fmt.Println(path)

		routine := Create{
			path: path,
		}

		err = routine.Execute()
		if err == nil {
			t.Fail()
		}
	})

	t.Run("expected directories exist in a new project", func(t *testing.T) {
		dir, err := config.CreateTemp(t, "")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		routine := Create{
			path: dir,
		}

		err = routine.Execute()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		data := filepath.Join(dir, "data")
		routes := filepath.Join(dir, "routes")
		static := filepath.Join(dir, "static")
		templates := filepath.Join(dir, "templates")
		config := filepath.Join(dir, "onyx.yaml")

		expected := map[string]bool{
			data:      true,
			routes:    true,
			static:    true,
			templates: true,
			config:    false,
		}

		for path, isDir := range expected {
			file, err := os.Stat(path)
			dir := file.IsDir()
			fmt.Println(dir)
			if err != nil || file.IsDir() != isDir {
				t.FailNow()
			}
		}
	})
}
