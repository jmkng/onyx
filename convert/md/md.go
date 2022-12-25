package md

import (
	"io"

	_md "github.com/yuin/goldmark"
)

func Unmarshal(data []byte, out io.Writer) error {
	err := _md.Convert(data, out)
	if err != nil {
		return err
	}

	return nil
}
