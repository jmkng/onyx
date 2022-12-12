package routines

import (
	"os"
	"path"
	"testing"

	"github.com/jmkng/onyx/conf"
)

func TestEvaluatePath(t *testing.T) {
	t.Run("Non-zero value string is returned.", func(t *testing.T) {
		want := "example/path"

		got, err := EvaluatePath(want)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("Zero value string is returned as the current directory.", func(t *testing.T) {
		cd, err := os.Getwd()
		if err != nil {
			t.Fail()
		}

		want := cd

		got, err := EvaluatePath("")
		if err != nil || got != want {
			t.Fail()
		}
	})
}

func TestAnyExist(t *testing.T) {
	t.Run("config.yaml is found", func(t *testing.T) {
		dir, err := CreateTemp(t, "config.yaml")
		if err != nil {
			t.Fail()
		}

		want := path.Join(dir, "config.yaml")

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("config.yml is found", func(t *testing.T) {
		dir, err := CreateTemp(t, "config.yml")
		if err != nil {
			t.Fail()
		}

		want := path.Join(dir, "config.yml")

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("config.json is found", func(t *testing.T) {
		dir, err := CreateTemp(t, "config.json")
		if err != nil {
			t.Fail()
		}

		want := path.Join(dir, "config.json")

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("config.xml is found", func(t *testing.T) {
		dir, err := CreateTemp(t, "config.xml")
		if err != nil {
			t.Fail()
		}

		want := path.Join(dir, "config.xml")

		got, err := searchConf(dir)
		if err != nil || got != want {
			t.Fail()
		}
	})

	t.Run("config.toml is found", func(t *testing.T) {
		dir, err := CreateTemp(t, "config.toml")
		if err != nil {
			t.Fail()
		}

		want := path.Join(dir, "config.toml")

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

// searchConf will use the AnyExist function to compare the given dir with an array of
// recognized configuration file names.
func searchConf(dir string) (string, error) {
	result, err := AnyExist(dir, conf.Names[:])
	if err != nil {
		return "", err
	}

	return result, nil
}
