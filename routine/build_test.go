package routine

import (
	"testing"
)

func TestIsIgnored(t *testing.T) {
	t.Run("hidden unix files are ignored", func(t *testing.T) {
		ignored := isIgnored(".test")

		if !ignored {
			t.Fail()
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

func TestExtract(t *testing.T) {
	t.Run("YAML is recognized and extracted", func(t *testing.T) {
		mock := "---\ntitle: test\nauthor: test\n---body"

		head, body, err := extract(mock)
		if err != nil {
			t.Logf("extract returned err: %v", err)
			t.FailNow()
		}

		if head != "\ntitle: test\nauthor: test\n" || body != "body" {
			t.Logf("body/head are unexpected values\nhead: %v\nbody: %v", head, body)
			t.Fail()
		}
	})
}

func TestDestination(t *testing.T) {
	t.Run("`index` and `tmpl` extensions on index files return the same", func(t *testing.T) {
		paths := []string{
			"project/routes/index.tmpl",
			"project/routes/index.html",
		}

		for _, v := range paths {
			result := destination(v)
			if result != "project/build/index.html" {
				t.Fail()
			}
		}
	})

	t.Run("relative path to index file", func(t *testing.T) {
		mock := "project/routes/index.tmpl"

		path := destination(mock)

		if path != "project/build/index.html" {
			t.Errorf("invalid result for relative path to index file `%v`\nreceived: %v", mock, path)
			t.Fail()
		}
	})

	t.Run("relative path to group member", func(t *testing.T) {
		mock := "project/routes/posts/post-one.md"

		path := destination(mock)

		if path != "project/build/post-one/index.html" {
			t.Errorf("gave `%v` received `%v`", mock, path)
			t.Fail()
		}
	})
}
