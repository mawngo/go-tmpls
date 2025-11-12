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
//   - {{block "name" pipeline}}
const templateNameRegexStr = `{{\s*(template|define|block)\s+"([^"]+)"(?:\s+[\s\S]*?)?\s*}}`

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
	nostack bool

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
		nostack:       opt.nostack,
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

func (t *Templates) resolve(base Template, stackMap map[string][]string, name string) (Template, error) {
	path, ok := t.nameMap[name]
	if !ok {
		return nil, errors.New("template name not found '" + name + "'")
	}

	shouldBuildStack := false
	if base == nil {
		var err error
		base, err = t.baseFn(name)
		if err != nil {
			return nil, err
		}
		stackMap = make(map[string][]string)
		shouldBuildStack = true
	}

	b, err := fs.ReadFile(t.fs, path)
	if err != nil {
		return nil, err
	}

	content := string(b)
	includedTemplateMatches := t.templateNameRegex.FindAllStringSubmatchIndex(content, -1)
	// Exclude template names that are also having {{define}} and {{block}} block.
	excludedTemplateNames := make(map[string]struct{})
	replacements := make([]string, 0, 10)
	for _, v := range includedTemplateMatches {
		action := content[v[2]:v[3]]
		if action == "template" {
			continue
		}
		templateName := content[v[4]:v[5]]
		if len(templateName) == 0 {
			continue
		}
		excludedTemplateNames[templateName] = struct{}{}
		// Register push stacked template.
		if !t.nostack && action == "define" && strings.HasPrefix(templateName, "@stack:") {
			nameList, ok := stackMap[templateName]
			if !ok {
				nameList = make([]string, 0, 5)
			}
			replacedName := templateName + ":" + name
			stackMap[templateName] = append(nameList, replacedName)

			originalMatch := content[v[0]:v[1]]
			replacedMatch := content[v[0]:v[4]] + replacedName + content[v[5]:v[1]]
			replacements = append(replacements, originalMatch, replacedMatch)
		}
	}

	// Process included templates top down, then finally parse the content.
	for _, v := range includedTemplateMatches {
		name := content[v[4]:v[5]]
		if len(name) == 0 {
			continue
		}
		if _, ok := excludedTemplateNames[name]; ok {
			continue
		}
		// Register stacked template name.
		if !t.nostack && strings.HasPrefix(name, "@stack:") {
			if _, ok := stackMap[name]; !ok {
				stackMap[name] = make([]string, 0, 5)
			}
			continue
		}
		var err error
		base, err = t.resolve(base, stackMap, name)
		if err != nil {
			return nil, err
		}
	}

	// Handle stacked templates.
	if shouldBuildStack {
		for stackName, pushedNames := range stackMap {
			stackContent := ""
			if len(pushedNames) > 0 {
				var sb strings.Builder
				// Build stack template.
				for _, name := range pushedNames {
					sb.WriteString("{{template \"")
					sb.WriteString(name)
					sb.WriteString("\" .}}")
				}
				stackContent = sb.String()
			}
			base, err = base.New(stackName).Parse(stackContent)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(replacements) > 0 {
		replacer := strings.NewReplacer(replacements...)
		content = replacer.Replace(content)
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
	tmpl, err := t.lookup(name)
	if err != nil {
		return nil, err
	}
	return tmpl.Clone()
}

// lookup returns a template by name.
func (t *Templates) lookup(name string) (Template, error) {
	if t.nocache {
		t.mu.Lock()
		defer t.mu.Unlock()
		return t.resolve(nil, nil, name)
	}

	t.mu.RLock()
	if tmpl, ok := t.templateMap[name]; ok {
		defer t.mu.RUnlock()
		return tmpl, nil
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()
	if tmpl, ok := t.templateMap[name]; ok {
		return tmpl, nil
	}

	tmpl, err := t.resolve(nil, nil, name)
	if err != nil {
		return nil, err
	}
	t.templateMap[name] = tmpl
	return tmpl, nil
}

// ExecuteTemplate execute the specified template with the given data.
func (t *Templates) ExecuteTemplate(wr io.Writer, name string, data any) error {
	tmpl, err := t.lookup(name)
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
