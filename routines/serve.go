package routines

type Serve struct {
	Path string
	Port int
}

func (s Serve) Execute() error {
	return nil
}
