package tmpls

import (
	"errors"
	"fmt"
	"github.com/mawngo/go-tmpls/v2/internal"
	htmltemplate "html/template"
	"io"
	"io/fs"
	"net/http"
	fspath "path"
	"regexp"
	"strings"
	"sync"
)

// templateNameRegex regex for matching template name in
//   - {{template "name"}}
//   - {{template "name" pipeline}}
//   - {{define "name"}}
const templateNameRegexStr = `{{\s*(template|define)\s+"([^"]+)"(?:\s+[\s\S]*?)?\s*}}`

// Templates collection of cached and preprocessed templates.
type Templates struct {
	fs            fs.FS
	extensions    map[string]struct{}
	prefixMap     map[string]string
	separator     string
	onExecute     OnTemplateExecuteFn
	preloadFilter func(name string, path string) bool

	baseFn func(name string) (Template, error)
	// Map of parsed template by name.
	templateMap map[string]Template
	// Map of processed template name to template paths.
	nameMap map[string]string
	mu      sync.RWMutex
	nocache bool

	templateNameRegex *regexp.Regexp
}

// New create a new [Templates] instance.
// On creation, all templates in the specified file system will be parsed.
//
// See [TemplatesOption] for more configurations.
func New(fs fs.FS, options ...TemplatesOption) (*Templates, error) {
	opt := templatesOptions{
		pathSeparator: ".",
		extensions: map[string]struct{}{
			".html":   {},
			".gohtml": {},
			".gotxt":  {},
		},
		initFn: func(name string) Template {
			return htmlTemplate{
				Template: htmltemplate.New(name),
			}
		},
		preloadFilter: func(name string, _ string) bool {
			return name[0] != '_'
		},
	}

	for _, option := range options {
		option(&opt)
	}

	t := &Templates{
		fs:            fs,
		nocache:       opt.nocache,
		extensions:    opt.extensions,
		prefixMap:     opt.prefixMap,
		separator:     opt.pathSeparator,
		onExecute:     opt.onExecute,
		preloadFilter: opt.preloadFilter,

		nameMap:           make(map[string]string),
		templateMap:       make(map[string]Template),
		templateNameRegex: regexp.MustCompile(templateNameRegexStr),
	}

	t.baseFn = func(name string) (Template, error) {
		base := opt.initFn(name)
		if !opt.disableBuiltins {
			if funcs := internal.NewBuiltinFuncMap(opt.excludeFuncs...); len(funcs) > 0 {
				base = base.Funcs(funcs)
			}
		}
		if len(opt.funcs) > 0 {
			base = base.Funcs(opt.funcs)
		}
		return base, nil
	}

	if err := t.scanNames(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Templates) scanNames() error {
	return fs.WalkDir(t.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		path = fspath.Clean(path)
		ext := fspath.Ext(path)
		if len(t.extensions) > 0 {
			if _, ok := t.extensions[ext]; !ok {
				return nil
			}
		}

		name := path
		for prefix, replace := range t.prefixMap {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = replace + name[len(prefix):]
			break
		}
		name = strings.Join(strings.Split(name, "/"), t.separator)
		name = strings.TrimSuffix(name, ext)

		if prevPath, ok := t.nameMap[name]; ok {
			return fmt.Errorf(`template name conflict: "%s" (files %s and %s)`, name, prevPath, path)
		}
		t.nameMap[name] = path
		return err
	})
}

func (t *Templates) resolve(base Template, name string) (Template, error) {
	path, ok := t.nameMap[name]
	if !ok {
		return nil, errors.New("template name not found '" + name + "'")
	}

	if base == nil {
		var err error
		base, err = t.baseFn(name)
		if err != nil {
			return nil, err
		}
	}

	b, err := fs.ReadFile(t.fs, path)
	if err != nil {
		return nil, err
	}

	content := string(b)
	includedTemplateMatches := t.templateNameRegex.FindAllStringSubmatch(content, -1)
	// Exclude template names that are also having {{define}} block.
	excludedTemplateNames := make(map[string]struct{})
	for _, v := range includedTemplateMatches {
		if v[1] == "template" {
			continue
		}
		name := v[2]
		if len(name) == 0 {
			continue
		}
		excludedTemplateNames[name] = struct{}{}
	}

	// Process included templates top down, then finally parse the content.
	for _, v := range includedTemplateMatches {
		name := v[2]
		if len(name) == 0 {
			continue
		}
		if _, ok := excludedTemplateNames[name]; ok {
			continue
		}
		var err error
		base, err = t.resolve(base, name)
		if err != nil {
			return nil, err
		}
	}
	return base.New(name).Parse(content)
}

// Preload parse all scanned templates.
// You can configure which templates to preload by [WithPreloadFilter].
// By default, any template whose resolved name starts with an underscore (_) will be ignored.
func (t *Templates) Preload() ([]Template, error) {
	res := make([]Template, 0, 10)
	for name, path := range t.nameMap {
		if t.preloadFilter != nil {
			if !t.preloadFilter(name, path) {
				continue
			}
		}
		temp, err := t.Lookup(name)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

// Lookup returns a cloned template by name.
// If the template does not exist, it returns nil.
func (t *Templates) Lookup(name string) (Template, error) {
	if t.nocache {
		t.mu.Lock()
		defer t.mu.Unlock()
		return t.resolve(nil, name)
	}

	t.mu.RLock()
	if tmpl, ok := t.templateMap[name]; ok {
		defer t.mu.RUnlock()
		return tmpl.Clone()
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()
	if tmpl, ok := t.templateMap[name]; ok {
		return tmpl.Clone()
	}

	tmpl, err := t.resolve(nil, name)
	if err != nil {
		return nil, err
	}
	t.templateMap[name] = tmpl
	return tmpl.Clone()
}

// ExecuteTemplate execute the specified template with the given data.
func (t *Templates) ExecuteTemplate(wr io.Writer, name string, data any) error {
	tmpl, err := t.Lookup(name)
	if err != nil {
		return err
	}
	if t.onExecute != nil {
		if err := t.onExecute(tmpl, wr, data); err != nil {
			return err
		}
	}
	return tmpl.Execute(wr, data)
}

// MustExecuteTemplate execute the specified template with the given data and panic if any error occurs.
func (t *Templates) MustExecuteTemplate(wr io.Writer, name string, data any) {
	err := t.ExecuteTemplate(wr, name, data)
	if err != nil {
		panic(err)
	}
}

// NewStandardWebFS sets up [Templates] and [http.FileServer] based on golang-standards/project-layout,
// which read templates from web/template and serve static files from web/static.
func NewStandardWebFS(cwd fs.FS, options ...TemplatesOption) (*Templates, http.Handler, error) {
	fileSystem, err := fs.Sub(cwd, "web")
	if err != nil {
		return nil, nil, err
	}
	staticFs, err := fs.Sub(fileSystem, "static")
	if err != nil {
		return nil, nil, err
	}

	templateFs, err := fs.Sub(fileSystem, "template")
	if err != nil {
		return nil, nil, err
	}
	templateCache, err := New(templateFs, options...)
	if err != nil {
		return nil, nil, err
	}
	return templateCache, http.FileServer(http.FS(staticFs)), nil
}

// NewStandardTemplateFS returns [Templates] based on golang-standards/project-layout,
// which read templates from web/template.
func NewStandardTemplateFS(cwd fs.FS, options ...TemplatesOption) (*Templates, error) {
	fileSystem, err := fs.Sub(cwd, "web")
	if err != nil {
		return nil, err
	}

	templateFs, err := fs.Sub(fileSystem, "template")
	if err != nil {
		return nil, err
	}
	return New(templateFs, options...)
}
