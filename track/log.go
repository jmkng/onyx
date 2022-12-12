package track

import (
	"fmt"

	"github.com/jmkng/mute"
)

var logs []string

var logger = mute.Init(
	mute.Route{
		Memory: &logs,
		Format: mute.Text,
	},
)

// Log is a wrapper for Logger.Send().
func Log(e ...mute.Event) {
	logger.Send(e...)
}

// Event creates and returns a mute.Event.
func Event(msg string) mute.Event {
	return mute.Event{
		Message: msg,
	}
}

// Report prints all available log messages.
func Report() {
	for _, v := range logs {
		fmt.Println(v)
	}
}
