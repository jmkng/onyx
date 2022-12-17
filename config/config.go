package config

const (
	YamlName    = "onyx.yaml"
	YamlAltName = "onyx.yml"
	JsonName    = "onyx.json"
	TomlName    = "onyx.toml"
	XmlName     = "onyx.xml"
)

// State contains the config from the project being operated on
// during runtime.
var State Options

// Names contains all possible project configuration file names.
var Names = []string{
	YamlName,
	YamlAltName,
	JsonName,
	TomlName,
	XmlName,
}

// Options describes all of the options that might be found in a
// recognized configuration file.
type Options struct {
	Output   string   `json:"output"`
	Preserve []string `json:"preserve"`
}
