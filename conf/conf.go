package conf

// State contains the config from the project being operated on
// during runtime.
var State Options

// Names contains all possible project configuration file names.
var Names = []string{
	"config.yml",
	"config.yaml",
	"config.json",
	"config.toml",
	"config.xml",
}

// Options describes all of the options that might be found in a
// recognized configuration file.
type Options struct {
	// X           map[string]string `json:"-"`
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Include     string   `json:"include"`
	Exclude     string   `json:"exclude"`
	Preserve    []string `json:"preserve"`
}
