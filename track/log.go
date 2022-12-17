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

// Log will create an Event{} for every given string, and pass each Event{} to Send().
func Log(e ...string) {
	for _, v := range e {
		logger.Send(mute.Event{
			Message: v,
		})
	}
}

// Report prints all available log messages.
func Report() {
	for _, v := range logs {
		fmt.Println(v)
	}
}
