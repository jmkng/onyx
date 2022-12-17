package routine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jmkng/onyx/config"
)

func TestAnyExist(t *testing.T) {
	t.Run("'yaml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, config.YamlName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, config.YamlName)

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'yml' config is found", func(t *testing.T) {
		dir, err := CreateTemp(t, config.YamlAltName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, config.YamlAltName)

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'json' is found", func(t *testing.T) {
		dir, err := CreateTemp(t, config.JsonName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, config.JsonName)

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'xml' is found", func(t *testing.T) {
		dir, err := CreateTemp(t, config.XmlName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, config.XmlName)

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("'toml' is found", func(t *testing.T) {
		dir, err := CreateTemp(t, config.TomlName)
		if err != nil {
			t.Fail()
		}

		want := filepath.Join(dir, config.TomlName)

		got, err := searchConf(dir)
		if os.IsNotExist(err) || got != want {
			t.Fail()
		}
	})

	t.Run("unrecognized configuration is not found", func(t *testing.T) {
		want := ""

		dir, err := CreateTemp(t, "config.txt")
		if err != nil {
			t.Fail()
		}

		got, err := searchConf(dir)
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

		got, err := searchConf(dir)
		if !os.IsNotExist(err) || got != want {
			t.Fail()
		}
	})
}

// searchConf will search the given directory for a recognized configuration file.
func searchConf(dir string) (string, error) {
	result, err := AnyExist(dir, config.Names[:])
	if err != nil {
		return "", err
	}

	return result, nil
}
