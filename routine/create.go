package routine

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmkng/onyx/config"
	"github.com/jmkng/onyx/convert/yaml"
	"github.com/jmkng/onyx/track"
)

func NewCreate() *Create {
	c := &Create{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.path, "path", WdOrPanic(), "Path to the desired location of the new project.")

	return c
}

type Create struct {
	fs   *flag.FlagSet
	path string
}

func (c *Create) Name() string {
	return c.fs.Name()
}

func (c *Create) Parse(args []string) error {
	return c.fs.Parse(args)
}

func (c *Create) Execute() error {
	info, err := os.Stat(c.path)

	if err != nil && errors.Is(err, os.ErrNotExist) && c.path != "" {
		mkErr := os.Mkdir(c.path, DefDirPerm)
		if mkErr != nil {
			return fmt.Errorf("failed to create directory: %v", c.path)
		}

		track.Log(fmt.Sprintf("created: %v", c.path))
	} else if err != nil && c.path != "" {
		return fmt.Errorf("failed to access directory, check permission: %v", c.path)
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("path leads to file, expected path to new location or directory: %v", c.path)
	}

	configPath := filepath.Join(c.path, config.YamlName)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("found existing configuration file, please rename, move or delete: %v", configPath)
	}

	track.Log(
		fmt.Sprintf("created: %v", configPath),
	)

	toTemplates := filepath.Join(c.path, "templates")
	toRoutes := filepath.Join(c.path, "routes")

	dirs := []string{
		toTemplates,
		toRoutes,
		filepath.Join(c.path, "static"),
		filepath.Join(c.path, "data"),
	}

	for _, v := range dirs {
		err := os.Mkdir(v, DefDirPerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", v)
		}

		track.Log(fmt.Sprintf("created: %v", v))
	}

	base := filepath.Join(toTemplates, "base.tmpl")
	index := filepath.Join(toRoutes, "index.tmpl")

	conf, err := yaml.Marshal(config.Options{})
	if err != nil {
		panic("failed to marshal a default configuration")
	}

	err = os.WriteFile(configPath, []byte(conf), DefFilePerm)
	if err != nil {
		return fmt.Errorf("failed to access file: %v", configPath)
	}

	files := map[string]string{
		base:  "",
		index: "",
	}

	for path, data := range files {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", path)
		}

		defer func() {
			err = file.Close()
			if err != nil {
				panic(fmt.Sprintf("failed to access file: %v", path))
			}
		}()

		_, err = file.WriteString(data)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", path)
		}

		track.Log(fmt.Sprintf("created: %v", path))
	}

	return nil
}
