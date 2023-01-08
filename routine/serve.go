package routine

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
)

func NewServe() *Serve {
	routine := &Serve{
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
	}

	routine.fs.StringVar(&routine.path, "path", WdOrPanic(), "Path to the project being served.")
	routine.fs.IntVar(&routine.port, "port", 3883, "Port used to host the site.")
	routine.fs.BoolVar(&routine.verbose, "verbose", false, "Display more detailed information")

	return routine
}

type Serve struct {
	fs      *flag.FlagSet
	path    string
	port    int
	verbose bool
}

func (routine *Serve) Name() string {
	return routine.fs.Name()
}

func (routine *Serve) Parse(args []string) error {
	return routine.fs.Parse(args)
}

func (routine *Serve) Execute() error {
	err := Setup(routine.path)
	if err != nil {
		return err
	}

	first := routine.port
	last := 0

	count := 0
	for {
		if count > 100 {
			last = routine.port
			return fmt.Errorf("unable to secure port between %v - %v, please specify with --port`", first, last)
		}

		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", routine.port))
		if err != nil {
			if routine.port != 3883 {
				return fmt.Errorf("requested port `%v` is already in use", routine.port)
			}

			routine.port++
		} else {
			ln.Close()
			break
		}

		count++
	}

	fmt.Printf("serving on http://localhost:%v\n", routine.port)

	http.Handle("/", http.FileServer(
		http.Dir(
			filepath.Join(routine.path, "build"),
		),
	))

	err = http.ListenAndServe(":"+fmt.Sprint(routine.port), nil)
	if err != nil {
		return fmt.Errorf("failed to host server on http://localhost:%v", routine.port)
	}

	return nil
}

// func (routine *Serve) handler(w http.ResponseWriter, req *http.Request) {
// 	url := req.URL
//
// 	var output string
// 	if config.State.Output != "" {
// 		output = config.State.Output
// 	} else {
// 		output = "build"
// 	}
//
// 	var request string
// 	ext := filepath.Ext(url.Path)
// 	if ext != "" {
// 		request = filepath.Join(routine.path, output, url.Path)
// 	} else {
// 		request = filepath.Join(routine.path, output, url.Path, "index.html")
// 	}
//
// 	_, err := os.Stat(request)
// 	if err != nil {
// 		fmt.Fprint(w, "404")
// 		return
// 	}
//
// 	file, err := os.ReadFile(request)
// 	if err != nil {
// 		fmt.Fprint(w, "401")
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
//
// 	fmt.Fprint(w, string(file))
// }
