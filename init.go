package main

import (
	"flag"
	"fmt"

	"github.com/jmkng/onyx/routine"
)

// Executable describes a type that surfaces some methods to interface
// with an internal flagset, and has an Execute() method to coordinate
// the execution of the routine.
type Executable interface {
	// Name wraps the internal type's FlagSet and returns the result of
	// calling the FlagSet.Name() method.
	Name() string
	// Parse wraps the internal type's FlagSet and passes the given []string
	// to the FlagSet.Parse([]string) method. The return value will be ErrHelp
	// if -help or -h were set but not defined.
	Parse([]string) error
	// Execute will begin the routine's execution. If a fatal error is encountered
	// during runtime, execution will return early with an error.
	Execute() error
}

// Init will attempt to find a routine with a name that matches the first
// received argument. If one is found, the remaining arguments are passed
// to the routine's internal FlagSet.Parse([]string) method. After the
// arguments are parsed, the routine's .Execute() method is called to perform
// the action. If a fatal error occurs during runtime, the routine returns
// early with the error.
func Init(args []string) error {
	flag.Usage = func() {
		fmt.Println("Onyx commands:")
		fmt.Println("> create")
		fmt.Println("> build")
		fmt.Println("> serve")
		fmt.Println("View command information with `onyx <command> [--help|-h]`")
	}

	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	routines := []Executable{
		routine.NewCreate(),
		routine.NewBuild(),
		routine.NewServe(),
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
