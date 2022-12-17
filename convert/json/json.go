package json

import (
	_json "encoding/json"
)

func Unmarshal(data []byte, out any) error {
	return _json.Unmarshal(data, out)
}
