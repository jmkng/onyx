package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnyExist(t *testing.T) {
	t.Run("'yaml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, YamlName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, YamlName)

		got, err := SearchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'yml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, YamlAltName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, YamlAltName)

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
