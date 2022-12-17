package routine

import (
	"fmt"
	"path"
	"testing"
)

func TestCreateExecute(t *testing.T) {
	t.Run("Config file is created when path argument does not already contain a config.", func(t *testing.T) {
		dir, err := CreateTemp(t, "")
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

	t.Run("Directory will be created if path argument is a non-existent directory.", func(t *testing.T) {
		dir, err := CreateTemp(t, "")
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

	t.Run("Existing config file will not be overridden.", func(t *testing.T) {
		dir, err := CreateTemp(t, "onyx.yaml")
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

	t.Run("Path to a file cannot be used as path argument.", func(t *testing.T) {
		dir, err := CreateTemp(t, "onyx.yaml")
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
}
