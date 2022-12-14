package main

import (
	"flag"
	"fmt"

	"github.com/jmkng/onyx/routines"
)

type Executable interface {
	Name() string
	Parse([]string) error
	Execute() error
}

func Init(args []string) error {
	flag.Usage = func() {
		fmt.Println("Onyx commands:")
		fmt.Println("- create")
		fmt.Println("- build")
		fmt.Println("- serve")
		fmt.Println("View command flags with `onyx <command> [--help|-h]`")
	}

	flag.Parse()

	if len(args) < 1 {
		flag.Usage()
		return nil
	}

	routines := []Executable{
		routines.NewCreate(),
		routines.NewBuild(),
		routines.NewServe(),
	}

	for _, rt := range routines {
		if rt.Name() == args[0] {
			err := rt.Parse(args[1:])
			if err != nil {
				return err
			}

			return rt.Execute()
		}
	}

	return fmt.Errorf("unknown command: %s", args[0])
}
