package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetState(t *testing.T) {
	t.Run("valid configuration is unmarshaled", func(t *testing.T) {
		dir, err := CreateTemp(t, YamlLongName)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		toFile := filepath.Join(dir, YamlLongName)
		data := []byte("---\ntitle: test\ndate: 2022-01-01\n---")
		err = os.WriteFile(toFile, data, 0644)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		err = SetState(toFile)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	})
}

func TestAnyExist(t *testing.T) {
	t.Run("'yaml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, YamlLongName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, YamlLongName)

		got, err := SearchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'yml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, YamlShortName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, YamlShortName)

		got, err := SearchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'json' is found", func(t *testing.T) {
		dir, err := CreateTemp(t, JsonName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, JsonName)

		got, err := SearchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("unrecognized configuration is not found", func(t *testing.T) {
		want := ""

		dir, err := CreateTemp(t, "config.txt")
		if err != nil {
			t.Fail()
		}

		got, err := SearchConf(dir)
		if !os.IsNotExist(err) || got != want {
			t.Fail()
		}
	})

	t.Run("no configuration is found", func(t *testing.T) {
		want := ""

		dir, err := CreateTemp(t, "")
		if err != nil {
			t.Fail()
		}

		got, err := SearchConf(dir)
		if !os.IsNotExist(err) || got != want {
			t.Fail()
		}
	})
}
