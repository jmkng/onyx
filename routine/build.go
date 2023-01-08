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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jmkng/onyx/config"
	"github.com/jmkng/onyx/convert/md"
	"github.com/jmkng/onyx/convert/yaml"
	"github.com/jmkng/onyx/track"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NewBuild() *Build {
	routine := &Build{
		fs: flag.NewFlagSet("build", flag.ContinueOnError),
	}

	routine.fs.StringVar(&routine.path, "path", WdOrPanic(), "Path to the project being built.")
	routine.fs.BoolVar(&routine.verbose, "verbose", false, "Display more detailed information")

	return routine
}

type Build struct {
	fs      *flag.FlagSet
	path    string
	verbose bool
}

func (routine *Build) Name() string {
	return routine.fs.Name()
}

func (routine *Build) Parse(args []string) error {
	return routine.fs.Parse(args)
}

func (routine *Build) Execute() error {
	err := Setup(routine.path)
	if err != nil {
		return err
	}

	routes := filepath.Join(routine.path, "routes")

	_, err = os.Stat(routes)
	if err != nil {
		return fmt.Errorf("project has no routes: %v", routine.path)
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
			if routine.verbose || config.State.Verbose {
				track.Log(fmt.Sprintf("skipped unrecognized file: %v", filepath.Base(path)))
			}

			return nil
		}

		resourceCt++

		go func() {
			res, err := newResource(routine.path, path, routine.verbose)

			resourceChan <- resourceEvent{
				res: res,
				err: err,
			}
		}()

		return nil
	})
	if err != nil {
		panic(err)
	}

	var render sync.WaitGroup
	render.Add(resourceCt)

	var injectable injectable
	var renderable []resource

	for i := 0; i < resourceCt; i++ {
		event := <-resourceChan
		if event.err != nil {
			return event.err
		}

		renderable = append(renderable, event.res)

		go func() {
			injectable.absorb(event.res)
			render.Done()
		}()
	}

	render.Wait()

	for _, v := range injectable.Data {
		group, ok := v.([]map[string]any)

		if !ok {
			continue
		}

		for _, v := range group {
			_, exists := v["Date"]
			if !exists || v["Date"] == "" {
				continue
			}

			asString, ok := v["Date"].(string)

			errDate := fmt.Errorf("invalid date provided to resource: %v", v["///Path"])

			if !ok {
				return errDate
			}

			date, err := time.Parse(config.DateFmt, asString)
			if err != nil {
				return errDate
			}

			v["***DATE"] = date
		}

		sort.Slice(group, func(i, j int) bool {
			iDate, ok := group[i]["***DATE"].(time.Time)
			if !ok {
				panic("failed to assert ***DATE as time.Time")
			}

			jDate, ok := group[j]["***DATE"].(time.Time)
			if !ok {
				panic("failed to assert ***DATE as time.Time")
			}

			return iDate.After(jDate)
		})
	}

	renderedChan := make(chan resourceEvent)
	renderedCt := 0

	for i := range renderable {
		renderedCt++

		toTemplates := filepath.Join(routine.path, "templates")
		toBase := filepath.Join(toTemplates, "base.tmpl")

		templates := []string{toBase}

		_, err = os.Stat(toBase)
		if err != nil {
			return fmt.Errorf("missing base template `base.tmpl` in %v", toTemplates)
		}

		var base []string
		err = filepath.WalkDir(toTemplates, func(path string, d fs.DirEntry, err error) error {
			if isIgnored(path) {
				if d.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			if d.IsDir() {
				if path == toTemplates {
					return nil
				}

				return filepath.SkipDir
			}

			if !d.IsDir() && d.Type().IsRegular() {
				pathBase := filepath.Base(path)
				segments := strings.Split(pathBase, ".")

				if len(segments) < 2 {
					return nil
				}

				prefix := strings.HasPrefix(segments[0], "base_")

				if !prefix && segments[1] == "tmpl" {
					return nil
				}

				base = append(base, path)
			}

			return nil
		})
		if err != nil {
			panic(err)
		}

		caser := cases.Title(language.English)

		templates = append(templates, base...)

		go func(res resource, out chan resourceEvent) {
			var request string

			if res.template != "" {
				request = filepath.Join(toTemplates, res.template)
			}

			if request != "" {
				_, err := os.Stat(request)
				if err != nil {
					renderedChan <- resourceEvent{
						res: res,
						err: fmt.Errorf("unable to locate template `%v` requested by: %v", res.template, filepath.Base(res.path)),
					}

					return
				}

				templates = append(templates, request)
			}

			if res.ext == ".tmpl" {
				var buf bytes.Buffer
				prerender := template.New("prerender")
				prerender.Parse(string(res.transformed))
				prerender.Execute(&buf, injectable.Data)
				res.transformed = template.HTML(buf.String())
			}

			context := make(map[string]any)

			context["Content"] = res.transformed

			for k, v := range injectable.Data {
				context[k] = v
			}

			for k, v := range res.data {
				title := caser.String(k)
				context[title] = v
			}

			tmpl, err := template.ParseFiles(templates...)
			if err != nil {
				wrapped := fmt.Errorf("failed to parse templates for resource: %v\n%v", res.path, err)

				renderedChan <- resourceEvent{
					res: res,
					err: wrapped,
				}

				return
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, context)
			if err != nil {
				renderedChan <- resourceEvent{
					res: resource{},
					err: fmt.Errorf("encountered a problem while executing template:\n%v", err.Error()),
				}
				return
			}

			res.rendered = buf.String()

			renderedChan <- resourceEvent{
				res: res,
				err: nil,
			}
		}(renderable[i], renderedChan)
	}

	for i := 0; i < renderedCt; i++ {
		event := <-renderedChan
		if event.err != nil {
			return event.err
		}

		parent := filepath.Dir(event.res.destination)

		_, err = os.Stat(parent)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = os.MkdirAll(parent, DefDirPerm)
				if err != nil {
					return fmt.Errorf("unable to create directory: %v", parent)
				}
			} else {
				panic(err)
			}
		}

		err = os.WriteFile(event.res.destination, []byte(event.res.rendered), DefFilePerm)
		if err != nil {
			return fmt.Errorf("unable to write file: %v", event.res.destination)
		}
	}

	var static []string

	toStatic := filepath.Join(routine.path, "static")
	err = filepath.WalkDir(toStatic, func(path string, d fs.DirEntry, err error) error {
		if isIgnored(path) {
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		if d.IsDir() {
			return nil
		}

		if !d.IsDir() && d.Type().IsRegular() {
			static = append(static, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, v := range static {
		dest, err := out(routine.path, v)
		if err != nil {
			panic(err)
		}

		parent := filepath.Dir(dest)

		err = os.MkdirAll(parent, DefDirPerm)
		if err != nil {
			return fmt.Errorf("unable to create directory: %v", parent)
		}

		bytes, err := os.ReadFile(v)
		if err != nil {
			return fmt.Errorf("unable to read file: %v", v)
		}

		err = os.WriteFile(dest, bytes, DefFilePerm)
		if err != nil {
			return fmt.Errorf("unable to write file: %v", dest)
		}
	}

	return nil
}

// pull will find data within a file and extract it, returning the data and
// content as separate strings. An error is returned if no data exists in the file,
// or the data is malformed in some way.
func pull(data string) (string, string, error) {
	if data[0:3] != "---" {
		return "", "", errors.New("no data in file")
	}

	firstEnd := 3
	secondStart := strings.Index(data[3:], "---")
	secondEnd := secondStart + 3

	head := data[firstEnd:secondEnd]
	body := data[(secondEnd + 3):]

	return head, body, nil
}

// diff will determine the difference between two paths, returning the relative
// part of a path from the first to second.
func diff(root, path string) (string, error) {
	var project string

	// normalize root
	if !filepath.IsAbs(root) {
		result, err := filepath.Abs(root)
		if err != nil {
			return "", err
		}

		project = result
	} else {
		project = root
	}

	if !filepath.IsAbs(path) {
		result, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}

		path = result
	}

	projSegments := strings.Split(project, string(filepath.Separator))
	pathSegments := strings.Split(path, string(filepath.Separator))

	projLen := len(projSegments)
	pathLen := len(pathSegments)

	if projLen > pathLen {
		return "", errors.New("path is not inside of project")
	}

	start := pathLen - (pathLen - projLen)

	var result []string

	for i := start; i < pathLen; i++ {
		result = append(result, pathSegments[i])
	}

	return filepath.Join(result...), nil
}

// out will examine the original path to a file and determine where it should be placed.
func out(root, path string) (string, error) {
	var resource string

	// normalize path to resource
	if !filepath.IsAbs(path) {
		full, err := filepath.Abs(path)
		if err != nil {
			return "", errors.New("unable to determine full path to resource")
		}

		resource = full
	} else {
		resource = path
	}

	var output string

	// normalize output directory
	if config.State.Output != "" {
		output = config.State.Output
	} else {
		output = "build"
	}

	// find relative path from project root to resource
	relative, err := diff(root, resource)
	if err != nil {
		return "", errors.New("unable to determine difference path to resource")
	}

	sep := filepath.Separator
	segments := strings.Split(relative, string(sep))

	// relative path segments are examined
	switch segments[0] {
	case "routes":
		// last segment
		last := segments[len(segments)-1]
		switch last {
		case "index.html", "index.tmpl":
			lastIndex := len(segments) - 1
			allBetween := segments[1:lastIndex]
			joined := filepath.Join(allBetween...)
			return filepath.Join(root, output, joined, "index.html"), nil
		default:
			lastIndex := len(segments) - 1
			last := segments[len(segments)-1]
			allBetween := segments[1:lastIndex]
			split := strings.Split(last, ".")

			if len(split) < 2 {
				return "", fmt.Errorf("resource is missing a valid extension: %v", path)
			}

			joined := filepath.Join(allBetween...)
			return filepath.Join(root, output, joined, split[0], "index.html"), nil
		}
	case "static":
		return filepath.Join(root, output, relative), nil
	default:
		panic(
			fmt.Sprintf("received call to determine output path for unexpected section in project: %v", segments[0]),
		)
	}
}

// isIgnored will return true if the path leads to a file that is ignored.
func isIgnored(path string) bool {
	result := false

	base := filepath.Base(path)

	if strings.HasPrefix(base, ".") {
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

// injectable is a data structure used to hold the key/value pairs from all
// resources in a project.
type injectable struct {
	Data map[string]any
	mu   sync.Mutex
}

// absorb will add a resource's data to an injectable context that can be used
// to render templates, and it is thread safe.
func (i *injectable) absorb(res resource) {
	if res.group == "" {
		return
	}

	caser := cases.Title(language.English)

	member := make(map[string]any)

	// hidden keys
	//
	// certain keys prefixed with "***" are intended to be used to
	// provide additional information when handling errors.

	member["***Path"] = res.path
	member["Content"] = res.transformed
	member["Date"] = res.date
	member["Link"] = res.link

	for k, v := range res.data {
		keyTitle := caser.String(k)
		member[keyTitle] = v
	}

	groupTitle := caser.String(res.group)

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

// newResource creates and initializes a new Resource from the file at filePath.
func newResource(root, file string, verbose bool) (resource, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return resource{}, err
	}

	ext := filepath.Ext(file)
	asStr := string(raw)

	dest, err := out(root, file)
	if err != nil {
		return resource{}, err
	}

	res := resource{
		path:        file,
		destination: dest,
		ext:         ext,
		raw:         asStr,
	}

	rel, err := diff(root, file)
	if err != nil {
		return resource{}, fmt.Errorf("unable to determine relative path to resource: %v", file)
	}

	var output string

	if config.State.Output != "" {
		output = config.State.Output
	} else {
		output = "build"
	}

	link, err := diff(filepath.Join(root, output), res.destination)
	if err != nil {
		return resource{}, err
	}

	res.link = ("/" + link)

	sep := string(filepath.Separator)
	segments := strings.Split(rel, sep)

	if len(segments) > 2 {
		parent := filepath.Dir(file)
		group := filepath.Base(parent)
		res.group = group
	}

	rawData := ""
	rawBody := asStr

	complex := isComplex(asStr)

	if complex {
		rawData, rawBody, err = pull(asStr)
		if err != nil {
			return resource{}, fmt.Errorf("unable to pull file: %v", file)
		}

		// TODO: (FEATURE) YAML is assumed, maybe allow JSON or TOML too.
		yaml.Unmarshal(
			[]byte(rawData), &res.data,
		)

		if template, ok := res.data["template"]; ok {
			res.template = template
			delete(res.data, "template")
		}

		if date, ok := res.data["date"]; ok {
			res.date = date
			delete(res.data, "date")
		}
	} else {
		if verbose || config.State.Verbose {
			track.Log(fmt.Sprintf("found no metadata in file: %v", filepath.Base(file)))
		}
	}

	var convertedBody string
	switch ext {
	case ".md":
		var buf bytes.Buffer
		err = md.Unmarshal([]byte(rawBody), &buf)
		convertedBody = buf.String()
	case ".html", ".tmpl":
		convertedBody = rawBody
	default:
		panic("converted unexpected file type")
	}

	if err != nil {
		return resource{}, err
	}

	res.transformed = template.HTML(convertedBody)

	return res, nil
}

// resource represents a file that is being processed as part of a project.
type resource struct {
	path        string
	destination string
	link        string
	ext         string
	raw         string
	transformed template.HTML
	rendered    string
	group       string
	template    string
	date        string
	data        map[string]string
}

// resourceEvent is a struct passed through channels that may contain a resource.
type resourceEvent struct {
	res resource
	err error
}
