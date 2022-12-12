package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmkng/onyx/routines"
)

func Parse() routines.Routine {
	var (
		Path string
		Port int
	)

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	newCmd.StringVar(&Path, "path", "", "Path to the desired location of the new project.")

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&Path, "path", "", "Path to the project being built.")

	ServeCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	ServeCmd.StringVar(&Path, "path", "", "Path to the project being served.")
	ServeCmd.IntVar(&Port, "port", 3883, "Port used to host the site.")

	flag.Usage = func() {
		fmt.Println("Try `onyx <command> [flags]`")
		fmt.Println("Onyx commands:")
		fmt.Println("- new")
		fmt.Println("- build")
		fmt.Println("- serve")
		fmt.Println("View command flags with `onyx <command> [--help|-h]`")
	}

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return nil
	}

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
		return routines.New{Path: Path}
	case "build":
		buildCmd.Parse(os.Args[2:])
		return routines.Build{Path: Path}
	case "serve":
		ServeCmd.Parse(os.Args[2:])
		return routines.Serve{Path: Path, Port: Port}
	default:
		return nil
	}
}
