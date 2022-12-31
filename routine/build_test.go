package routine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsIgnored(t *testing.T) {
	t.Run("hidden unix files are ignored", func(t *testing.T) {
		files := []string{
			".test",
			"path/to/.test",
		}

		for _, v := range files {
			ignored := isIgnored(v)

			if !ignored {
				t.Fail()
			}
		}
	})
}

func TestIsUnknown(t *testing.T) {
	t.Run(".[html|md|tmpl] are recognized", func(t *testing.T) {
		files := []string{
			"test.html",
			"test.md",
			"test.tmpl",
		}

		for _, v := range files {
			unknown := isUnknown(v)

			if unknown {
				t.Fail()
			}
		}
	})

	t.Run("return true for unrecognized type", func(t *testing.T) {
		unknown := isUnknown("test.mock")

		if !unknown {
			t.Fail()
		}
	})

	t.Run("return true for missing extension", func(t *testing.T) {
		unknown := isUnknown("test")

		if !unknown {
			t.Fail()
		}
	})
}

func TestIsComplex(t *testing.T) {
	t.Run("string with valid delimiters '---' and tabs '\t' is cleaned, and returns true", func(t *testing.T) {
		mock := `---
		title: test
		---`

		complex := isComplex(mock)

		if !complex {
			t.Fail()
		}
	})

	t.Run("string with valid delimiters '---' returns true", func(t *testing.T) {
		mock := "---\ntitle: test\n---"

		complex := isComplex(mock)

		if !complex {
			t.Fail()
		}
	})

	t.Run("one too few delimiters '--' returns false", func(t *testing.T) {
		mock := `--\ntitle: test\n---`

		complex := isComplex(mock)

		if complex {
			t.Fail()
		}
	})

	t.Run("one too many delimiters '----' is not complex", func(t *testing.T) {
		mock := `----
		title: test
		---`

		complex := isComplex(mock)

		if complex {
			t.Fail()
		}
	})
}

func TestPull(t *testing.T) {
	t.Run("YAML is recognized and extracted", func(t *testing.T) {
		mock := "---\ntitle: test\nauthor: test\n---body"

		head, body, err := pull(mock)
		if err != nil {
			t.Logf("pull returned err: %v", err)
			t.FailNow()
		}

		if head != "\ntitle: test\nauthor: test\n" || body != "body" {
			t.Logf("body/head are unexpected values\nhead: %v\nbody: %v", head, body)
			t.Fail()
		}
	})
}

func TestOut(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	t.Run("`index` and `tmpl` extensions on index files return the same", func(t *testing.T) {
		tmpl := filepath.Join("routes", "index.tmpl")
		html := filepath.Join("routes", "index.html")

		paths := []string{
			tmpl,
			html,
		}

		expected := filepath.Join(wd, "build", "index.html")

		for _, v := range paths {
			path, err := out(wd, v)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			result := path
			if result != expected {
				t.Logf("expected %v, got %v", expected, result)
				t.Fail()
			}
		}
	})

	t.Run("relative path to index file", func(t *testing.T) {
		mock := filepath.Join("routes", "index.tmpl")

		path, err := out(wd, mock)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		expected := filepath.Join(wd, "build", "index.html")
		if path != expected {
			t.Errorf("invalid result for relative path to index file `%v`\nreceived: %v", mock, path)
			t.Fail()
		}
	})

	t.Run("relative path to group member", func(t *testing.T) {
		mock := filepath.Join("routes", "posts", "post-one.md")

		path, err := out(wd, mock)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		expected := filepath.Join(wd, "build", "posts", "post-one", "index.html")
		if path != expected {
			t.Errorf("gave `%v` received `%v`", mock, path)
			t.Fail()
		}
	})
}
