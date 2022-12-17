package routine

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmkng/onyx/track"
)

func NewBuild() *Build {
	b := &Build{
		fs: flag.NewFlagSet("build", flag.ContinueOnError),
	}

	b.fs.StringVar(&b.path, "path", WdOrPanic(), "Path to the project being built.")

	return b
}

type Build struct {
	fs   *flag.FlagSet
	path string
}

func (b *Build) Name() string {
	return b.fs.Name()
}

func (b *Build) Parse(args []string) error {
	return b.fs.Parse(args)
}

func (b *Build) Execute() error {
	var wg sync.WaitGroup

	filepath.WalkDir(b.path, func(path string, d fs.DirEntry, err error) error {
		base := filepath.Base(path)

		if isIgnored(base) {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if d.IsDir() && b.path != path {
			if isReservedDir(path) {
				return filepath.SkipDir
			}

			wg.Add(1)
			go makeGroup(path, &wg)

			return filepath.SkipDir
		}

		if d.Type().IsRegular() {
			if isReservedFile(path) {
				return nil
			}

			if isUnknownFile(path) {
				// TODO: Log this in "verbose" mode.
				track.Log(fmt.Sprintf("unrecognized: %v", path))
				return nil
			}

			wg.Add(1)
			go makePage(path, &wg)
		}

		return nil
	})

	wg.Wait()
	return nil
}

func makeGroup(path string, wg *sync.WaitGroup) {
	fmt.Println("(DEBUG) making group: " + path)
	wg.Done()
}

func makePage(path string, wg *sync.WaitGroup) {
	fmt.Println("(DEBUG) making page: " + path)
	wg.Done()
}

// isReservedDir will return true if the path leads to a directory where the name is
// considered reserved.
func isReservedDir(path string) bool {
	name := filepath.Base(path)

	switch name {
	case "templates":
		return true
	default:
		return false
	}
}

// isReservedFile will return true if the path leads to a file where the name is
// considered reserved.
func isReservedFile(path string) bool {
	name := filepath.Base(path)

	switch name {
	case "onyx.yaml", "onyx.yml", "onyx.json", "onyx.toml", "onyx.xml":
		return true
	default:
		return false
	}
}

// isUnknownFile will return true if the path leads to a file where the extension
// is not recognized.
func isUnknownFile(path string) bool {
	ext := filepath.Ext(path)

	switch ext {
	case ".html", ".md", ".tmpl":
		return false
	default:
		return true
	}
}

// isIgnored will return true if the path leads to a file that is ignored
// by the system.
func isIgnored(path string) bool {
	result := false

	if strings.HasPrefix(path, ".") {
		result = true
	}

	return result
}
