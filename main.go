package main

import (
	"fmt"
	"os"

	"github.com/jmkng/onyx/track"
)

func main() {
	if err := Init(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	track.Report()
}
