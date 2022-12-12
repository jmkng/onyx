package yaml

import (
	_yaml "gopkg.in/yaml.v3"
)

func Unmarshal(text string, out any) error {
	err := _yaml.Unmarshal([]byte(text), out)
	return err
}
