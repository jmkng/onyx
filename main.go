package main

import (
	"fmt"
	"os"

	"github.com/jmkng/onyx/track"
)

func main() {
	routine := Parse()

	if routine != nil {
		err := routine.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	track.Report()
}
