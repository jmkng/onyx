package routine

import (
	"testing"

	"github.com/jmkng/onyx/config"
)

func TestIsIgnored(t *testing.T) {
	t.Run("unix hidden files are ignored", func(t *testing.T) {
		ignored := isIgnored(".test")

		if !ignored {
			t.Fail()
		}
	})
}

func TestIsReservedDir(t *testing.T) {
	t.Run("templates is reserved", func(t *testing.T) {
		reserved := isReservedDir("templates")

		if !reserved {
			t.Fail()
		}
	})
}

func TestIsReservedFile(t *testing.T) {
	t.Run("configuration files are reserved", func(t *testing.T) {
		for _, v := range config.Names {
			if reserved := isReservedFile(v); !reserved {
				t.Fail()
			}
		}
	})
}

func TestIsUnknownFile(t *testing.T) {
	t.Run(".csv (unrecognized) file is rejected", func(t *testing.T) {
		if recognized := isUnknownFile("test.csv"); !recognized {
			t.Fail()
		}
	})

	t.Run("recognized files are accepted", func(t *testing.T) {
		recognized := []string{
			"test.html",
			"test.md",
			"test.tmpl",
		}

		for _, v := range recognized {
			if recognized := isUnknownFile(v); recognized {
				t.Fail()
			}
		}
	})
}
