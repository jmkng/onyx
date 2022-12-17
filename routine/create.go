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

const examplePost = `---
date: 2022-12-15
---
content
`

// TODO: Write example template
const exampleTemplate = "example template not implemented"

// TODO: Write example index
const exampleIndex = "example index not implemented"

func NewCreate() *Create {
	c := &Create{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.path, "path", WdOrPanic(), "Path to the desired location of the new project.")
	c.fs.BoolVar(&c.example, "example", false, "Include additional files for newbies")

	return c
}

type Create struct {
	fs      *flag.FlagSet
	path    string
	example bool
}

func (c *Create) Name() string {
	return c.fs.Name()
}

func (c *Create) Parse(args []string) error {
	return c.fs.Parse(args)
}

func (c *Create) Execute() error {
	info, stErr := os.Stat(c.path)

	if stErr != nil && errors.Is(stErr, os.ErrNotExist) && c.path != "" {
		mkErr := os.Mkdir(c.path, DefDirPerm)
		if mkErr != nil {
			return fmt.Errorf("failed to access directory: %v", c.path)
		}

		track.Log(fmt.Sprintf("created: %v", c.path))
	} else if stErr != nil && c.path != "" {
		return fmt.Errorf("failed to access directory: %v", c.path)
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("path leads to file, expected directory: %v", c.path)
	}

	configPath := filepath.Join(c.path, config.YamlName)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration already exists in this directory: %v", configPath)
	}

	marConf, err := yaml.Marshal(config.Options{})
	if err != nil {
		panic("failed to marshal a default configuration")
	}

	wrErr := os.WriteFile(configPath, []byte(marConf), DefFilePerm)
	if wrErr != nil {
		return fmt.Errorf("failed to access file: %v", configPath)
	}

	track.Log(fmt.Sprintf("written: %v", configPath))

	if !c.example {
		return nil
	}

	postsName := "posts"
	templatesName := "templates"

	exampleDirs := []string{
		filepath.Join(c.path, templatesName),
		filepath.Join(c.path, postsName),
	}

	for _, v := range exampleDirs {
		err := os.Mkdir(v, DefDirPerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", v)
		}

		track.Log(fmt.Sprintf("created: %v", v))
	}

	indexPath := filepath.Join(c.path, "index.html")
	postPath := filepath.Join(c.path, postsName, "first-post.md")
	templatePath := filepath.Join(c.path, templatesName, "layout.tmpl")

	exampleFiles := map[string]string{
		exampleIndex:    indexPath,
		examplePost:     postPath,
		exampleTemplate: templatePath,
	}

	for data, path := range exampleFiles {
		file, err := os.Create(path)
		if err != nil {
			track.Log(fmt.Sprintf("failed to create file: %v", path))
			continue
		}

		track.Log(fmt.Sprintf("written: %v", path))

		_, err = file.WriteString(data)
		if err != nil {
			track.Log(fmt.Sprintf("failed to write file: %v", indexPath))
		}

		err = file.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to access file: %v", path))
		}
	}

	return nil
}
