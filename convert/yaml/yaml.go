package yaml

import (
	_yaml "gopkg.in/yaml.v3"
)

func Unmarshal(data []byte, out any) error {
	err := _yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

func Marshal(data any) ([]byte, error) {
	bytes, err := _yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
