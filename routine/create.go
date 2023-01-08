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
	routine := &Create{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}

	routine.fs.StringVar(&routine.path, "path", WdOrPanic(), "Path to the desired location of the new project.")
	routine.fs.BoolVar(&routine.verbose, "verbose", false, "Display more detailed information")

	return routine
}

type Create struct {
	fs      *flag.FlagSet
	path    string
	verbose bool
}

func (routine *Create) Name() string {
	return routine.fs.Name()
}

func (routine *Create) Parse(args []string) error {
	return routine.fs.Parse(args)
}

func (routine *Create) Execute() error {
	info, err := os.Stat(routine.path)

	if err != nil && errors.Is(err, os.ErrNotExist) && routine.path != "" {
		mkErr := os.Mkdir(routine.path, DefDirPerm)
		if mkErr != nil {
			// TODO: (BUG) If the directory already exists, leave it alone and move on.
			return fmt.Errorf("failed to create directory: %v", routine.path)
		}

		track.Log(fmt.Sprintf("created: %v", routine.path))
	} else if err != nil && routine.path != "" {
		return fmt.Errorf("failed to access directory, check permission: %v", routine.path)
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("path leads to file, expected path to new location or directory: %v", routine.path)
	}

	configPath := filepath.Join(routine.path, config.YamlLongName)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("found existing configuration file, please rename, move or delete: %v", configPath)
	}

	track.Log(
		fmt.Sprintf("created: %v", configPath),
	)

	toTemplates := filepath.Join(routine.path, "templates")
	toRoutes := filepath.Join(routine.path, "routes")

	dirs := []string{
		toTemplates,
		toRoutes,
		filepath.Join(routine.path, "static"),
		filepath.Join(routine.path, "data"),
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
