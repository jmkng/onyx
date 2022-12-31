package routine

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmkng/onyx/config"
)

func NewServe() *Serve {
	s := &Serve{
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
	}

	s.fs.StringVar(&s.path, "path", WdOrPanic(), "Path to the project being served.")
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
	err := Setup(s.path)
	if err != nil {
		return err
	}

	configPath, err := config.SearchConf(s.path)
	if err != nil {
		return fmt.Errorf("configuration file onyx.[yaml|yml|json] missing: %v", s.path)
	}

	err = config.Read(configPath)
	if err != nil {
		return fmt.Errorf("configuration file `%v` is malformed", configPath)
	}

	first := s.port
	last := 0

	count := 0
	for {
		if count > 100 {
			last = s.port
			return fmt.Errorf("unable to secure port between %v - %v, please specify with --port`", first, last)
		}

		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
		if err != nil {
			if s.port != 3883 {
				return fmt.Errorf("requested port `%v` is already in use", s.port)
			}

			s.port++
		} else {
			ln.Close()
			break
		}

		count++
	}

	fmt.Printf("serving on http://localhost:%v\n", s.port)

	http.HandleFunc("/", s.handler)

	err = http.ListenAndServe(":"+fmt.Sprint(s.port), nil)
	if err != nil {
		return fmt.Errorf("failed to host server on http://localhost:%v", s.port)
	}

	return nil
}

func (s *Serve) handler(w http.ResponseWriter, req *http.Request) {
	url := req.URL

	var request string

	var output string
	if config.State.Output != "" {
		output = config.State.Output
	} else {
		output = "build"
	}

	ext := filepath.Ext(url.Path)
	if ext != "" {
		request = filepath.Join(s.path, output, url.Path)
	} else {
		request = filepath.Join(s.path, output, url.Path, "index.html")
	}

	_, err := os.Stat(request)
	if err != nil {
		fmt.Fprint(w, "404")
		return
	}

	file, err := os.ReadFile(request)
	if err != nil {
		fmt.Fprint(w, "401")
		return
	}

	w.WriteHeader(http.StatusOK)

	// Set status code, content type, etc
	// write bytes with fmt.Fprintf

	fmt.Fprint(w, string(file))
}
