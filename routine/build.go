package routine

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmkng/onyx/config"
	"github.com/jmkng/onyx/convert/md"
	"github.com/jmkng/onyx/convert/yaml"
	"github.com/jmkng/onyx/track"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrNoData error = errors.New("no data in file")
)

func NewBuild() *Build {
	b := &Build{
		fs: flag.NewFlagSet("build", flag.ContinueOnError),
	}

	b.fs.StringVar(&b.path, "path", WdOrPanic(), "Path to the project being built.")

	return b
}

type Build struct {
	fs   *flag.FlagSet
	path string
}

func (b *Build) Name() string {
	return b.fs.Name()
}

func (b *Build) Parse(args []string) error {
	return b.fs.Parse(args)
}

func (b *Build) Execute() error {
	_, err := os.Stat(b.path)
	if err != nil {
		return fmt.Errorf("cannot access directory: %v", b.path)
	}

	configPath, err := config.SearchConf(b.path)
	if err != nil {
		return fmt.Errorf("configuration file onyx.[yaml|yml|json] missing: %v", b.path)
	}

	err = config.Read(configPath)
	if err != nil {
		return fmt.Errorf("configuration file `%v`  is malformed", configPath)
	}

	routes := filepath.Join(b.path, "routes")

	// TODO: (FEATURE) If config specifies a "domains" key, look for those folders instead
	// of a "routes" directory.
	_, err = os.Stat(routes)
	if err != nil {
		return fmt.Errorf("project has no routes: %v", b.path)
	}

	resourceChan := make(chan resourceEvent)

	resourceCt := 0
	err = filepath.WalkDir(routes, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if isIgnored(path) {
				return filepath.SkipDir
			}

			return nil
		}

		if !d.IsDir() && d.Type().IsRegular() && isIgnored(path) {
			return nil
		}

		if isUnknown(path) {
			// TODO: (BUG) Log in verbose
			track.Log(fmt.Sprintf("skipped unknown file: %v", path))
			return nil
		}

		resourceCt++

		go func() {
			res, err := NewResource(path)

			resourceChan <- resourceEvent{
				Res: res,
				Err: err,
			}
		}()

		return nil
	})
	if err != nil {
		panic(err)
	}

	var render sync.WaitGroup
	render.Add(resourceCt)

	var injectable Injectable
	var renderable []Resource

	for i := 0; i < resourceCt; i++ {
		event := <-resourceChan
		if event.Err != nil {
			return event.Err
		}

		renderable = append(renderable, event.Res)

		go func() {
			injectable.Absorb(event.Res)
			render.Done()
		}()
	}

	render.Wait() // TODO: (BUG) Write render logic

	// TODO: (BUG) Maybe I should stat these first?
	toTemplates := filepath.Join(b.path, "templates")
	toLayout := filepath.Join(toTemplates, "layout.tmpl")

	// err = filepath.WalkDir(toCommon, func(path string, d fs.DirEntry, err error) error {
	// 	if d.IsDir() {
	// 		return nil
	// 	}

	// 	ext := filepath.Ext(path)

	// 	if !d.Type().IsRegular() || ext != ".tmpl" {
	// 		// TODO: (FEATURE) Log in verbose.
	// 		track.Log(fmt.Sprintf("found template file of unrecognized type: %v", path))
	// 		return nil
	// 	}

	// 	common = append(common, path)
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }

	renderedChan := make(chan resourceEvent)

	renderedCt := 0
	renderedCt++

	func(res Resource, out chan resourceEvent) {
		fmt.Println(res)
		var request string

		// TODO: (BUG) Add extension if user just provides base name of file.

		if res.Template != "" {
			request = filepath.Join(toTemplates, res.Template)
		} else {
			if res.Group != "" {
				request = filepath.Join(toTemplates, res.Group)
			}
		}

		if request != "" {
			_, err := os.Stat(request)
			if err != nil {
				renderedChan <- resourceEvent{
					Res: res,
					Err: fmt.Errorf("unable to locate template `%v` requested by: %v", res.Template, res.Path),
				}

				return
			}
		}

		caser := cases.Title(language.English)

		context := make(map[string]any)

		context["Content"] = res.Transformed

		for k, v := range injectable.Data {
			context[k] = v
		}

		if len(res.Data) > 0 {
			for k, v := range res.Data {
				title := caser.String(k)
				context[title] = v
			}
		}

		var list []string

		if request != "" {
			list = []string{toLayout, request}
		} else {
			list = []string{toLayout}
		}

		tmpl, err := template.ParseFiles(list...)
		// template, err := template.ParseFiles(list...)
		if err != nil {
			panic(err)
			// renderedChan <- resourceEvent{
			// 	Res: res,
			// 	Err: fmt.Errorf("failed to parse templates for resource: %v", res.Path),
			// }
		}

		bonk := tmpl.DefinedTemplates()
		fmt.Println(bonk)

		var buf bytes.Buffer

		err = tmpl.ExecuteTemplate(&buf, "layout.tmpl", context)
		if err != nil {
			panic(err)
		}

		res.Rendered = buf.String()

		renderedChan <- resourceEvent{
			Res: res,
			Err: nil,
		}

	}(renderable[1], renderedChan)

	// for i := 0; i < rend  eredCt; i++ {
	<-renderedChan
	// event := <-renderedChan
	// fmt.Println(event.Res.Rendered)
	// fmt.Println("\n\n\n\n\n")
	// }

	return nil
}

// NewResource creates and initializes a new Resource from the file at filePath.
func NewResource(filePath string) (Resource, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return Resource{}, err
	}

	ext := filepath.Ext(filePath)
	asStr := string(raw)

	// TODO: (BUG) Calculate a path for the file.
	destination(filePath)

	res := Resource{
		Path:        filePath,
		Destination: destination(filePath),
		Ext:         ext,
		Raw:         asStr,
	}

	sep := string(filepath.Separator)
	segments := strings.Split(filePath, sep)

	if len(segments) > 3 {
		parent := filepath.Dir(filePath)
		group := filepath.Base(parent)
		res.Group = group
	}

	if isComplex(asStr) {
		rawData, rawBody, err := extract(asStr)
		if err != nil {
			return Resource{}, fmt.Errorf("unable to carve file: %v", filePath)
		}

		// TODO: (FEATURE) Detect type of metadata. We assume YAML here.
		yaml.Unmarshal(
			[]byte(rawData), &res.Data,
		)

		var buf bytes.Buffer
		switch ext {
		case ".md":
			err = md.Unmarshal([]byte(rawBody), &buf)
		}

		if err != nil {
			return Resource{}, err
		}

		res.Transformed = buf.String()

		// TODO: Promote values
		if template, ok := res.Data["template"]; ok {
			res.Template = template
			delete(res.Data, "template")
		}

		if date, ok := res.Data["date"]; ok {
			res.Date = date
			delete(res.Data, "date")
		}
	} else {
		res.Transformed = asStr
	}

	// TODO: (BUG) Finish initializing resource.
	return res, nil
}

type Resource struct {
	Path        string
	Destination string
	Ext         string
	Raw         string
	Transformed string
	Rendered    string
	Group       string
	Template    string
	Date        string
	Data        map[string]string
}

type Injectable struct {
	Data map[string]any
	// Data map[string][]map[string]any
	mu sync.Mutex
}

// Absorb will add a resource's data to an injectable context that can be used
// to render templates, and it is thread safe.
func (i *Injectable) Absorb(res Resource) {
	if res.Group == "" {
		return
	}

	caser := cases.Title(language.English)

	member := make(map[string]any)

	member["Content"] = res.Transformed
	member["Date"] = res.Date

	for k, v := range res.Data {
		keyTitle := caser.String(k)
		member[keyTitle] = v
	}

	groupTitle := caser.String(res.Group)

	i.mu.Lock()
	defer i.mu.Unlock()

	if i.Data == nil {
		i.Data = make(map[string]any)
	}

	if _, ok := i.Data[groupTitle]; !ok {
		i.Data[groupTitle] = []map[string]any{}
	}

	i.Data[groupTitle] = append(i.Data[groupTitle].([]map[string]any), member)
}

// ResourceEvent is a struct passed through channels that may contain a Resource.
type resourceEvent struct {
	Res Resource
	Err error
}

// isIgnored will return true if the path leads to a file that is ignored.
func isIgnored(path string) bool {
	result := false

	if strings.HasPrefix(path, ".") {
		result = true
	}

	return result
}

// isUnknown will return true if the path leads to a file where the extension
// is not recognized or does not contain a file extension.
func isUnknown(path string) bool {
	ext := filepath.Ext(path)

	switch ext {
	case ".html", ".md", ".tmpl":
		return false
	default:
		return true
	}
}

// isComplex will return true if the given string begins with a recognized delimiter
// to indicate that the file contains some data that needs to be extracted.
func isComplex(data string) bool {
	partitions := strings.Split(
		data,
		"\n",
	)

	found := 0

	// This should iterate over the lines in argument 'data' until the expected delimiters are found.
	for _, v := range partitions {
		if found == 2 {
			break
		}

		v = strings.ReplaceAll(v, "\t", "")

		if v == "---" {
			found++
		}
	}

	return found == 2
}

// extract will find data within a file and extract it, returning the data and
// content as separate strings. An error is returned if no data exists in the file,
// or the data is malformed in some way.
func extract(data string) (string, string, error) {
	if data[0:3] != "---" {
		return "", "", ErrNoData
	}

	firstEnd := 3
	secondStart := strings.Index(data[3:], "---")
	secondEnd := secondStart + 3

	head := data[firstEnd:secondEnd]
	body := data[(secondEnd + 3):]

	return head, body, nil
}

// destination will return an output path for a resource. The returned path
// is based on the name and original location of the resource.
func destination(path string) string {
	var result string

	sep := string(filepath.Separator)
	segments := strings.Split(path, sep)

	first := segments[0]
	file := segments[len(segments)-1]

	var root string

	if config.State.Output != "" {
		root = config.State.Output
	} else {
		if filepath.IsAbs(path) {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			root = filepath.Join(wd, "build")
		} else {
			root = filepath.Join(first, "build")
		}
	}

	fileSplit := strings.Split(file, ".")
	fileNoExt := fileSplit[0]
	if fileNoExt != "index" {
		result = filepath.Join(root, fileNoExt, "index.html")
	} else {
		result = filepath.Join(root, "index.html")
	}

	return result
}
