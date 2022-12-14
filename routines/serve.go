package routines

import (
	"flag"
	"fmt"
)

func NewServe() *Serve {
	s := &Serve{
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
	}

	s.fs.StringVar(&s.path, "path", "", "Path to the project being served.")
	s.fs.IntVar(&s.port, "port", 3883, "Port used to host the site.")

	return s
}

type Serve struct {
	fs   *flag.FlagSet
	path string
	port int
}

func (s *Serve) Name() string {
	return s.fs.Name()
}
func (s *Serve) Parse(args []string) error {
	return s.fs.Parse(args)
}

func (s *Serve) Execute() error {
	fmt.Println("Executed serve routine")
	return nil
}
