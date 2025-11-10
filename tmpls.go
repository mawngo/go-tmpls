package tmpls

import (
	"github.com/mawngo/go-tmpls/v2/internal"
	"html/template"
	"io"
	"io/fs"
	fspath "path"
	"strings"
	"sync"
)

// Templates collection of cached and preprocessed templates.
type Templates struct {
	fs         fs.FS
	extensions map[string]struct{}
	separator  string
	baseFn     func() (Template, error)
	onExecute  func(w io.Writer, t Template, name string, data any) error

	cached Template
	mu     sync.Mutex
}

func New(fs fs.FS, options ...TemplatesOption) (*Templates, error) {
	opt := templatesOptions{
		pathSeparator: ".",
		extensions: map[string]struct{}{
			".html":   {},
			".gohtml": {},
			".gotxt":  {},
		},
		initFn: func() Template {
			return htmlTemplate{
				Template: template.New(""),
			}
		},
	}

	for _, option := range options {
		option(&opt)
	}

	t := &Templates{
		fs:         fs,
		extensions: opt.extensions,
		separator:  opt.pathSeparator,
		onExecute:  opt.onExecute,
	}

	t.baseFn = func() (Template, error) {
		base := opt.initFn()
		if !opt.disableBuiltins {
			if funcs := internal.NewBuiltinFuncMap(opt.excludeFuncs...); len(funcs) > 0 {
				base = base.Funcs(funcs)
			}
		}
		if len(opt.funcs) > 0 {
			base = base.Funcs(opt.funcs)
		}
		return t.scan(fs, base)
	}

	// Catch potential errors early.
	base, err := t.baseFn()
	if err != nil {
		return nil, err
	}

	if !opt.nocache {
		t.cached = base
		t.baseFn = func() (Template, error) {
			return t.cached, nil
		}
	}
	return t, nil
}

func (t *Templates) scan(dir fs.FS, base Template) (Template, error) {
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
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
		name := strings.Join(strings.Split(path, "/"), t.separator)
		name = strings.TrimSuffix(name, ext)

		b, err := fs.ReadFile(dir, path)
		if err != nil {
			return err
		}
		base, err = base.New(name).Parse(string(b))
		return err
	})
	return base, err
}

// Base return cached template (cloned).
func (t *Templates) Base() (Template, error) {
	base, err := t.baseFn()
	if err != nil {
		return nil, err
	}
	return base.Clone()
}

// Register parses the text and adds the resulting template to the template collection.
func (t *Templates) Register(name string, text string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Cache is enabled.
	if t.cached != nil {
		base, err := t.cached.New(name).Parse(text)
		if err != nil {
			return err
		}
		t.cached = base
		return nil
	}

	// Cache is disabled.
	// Patch the base function to create a new template.
	// Not very effective, but usually the cache is only disabled in the dev environment.
	t.baseFn = func() (Template, error) {
		base, err := t.baseFn()
		if err != nil {
			return base, err
		}
		return base.New(name).Parse(text)
	}
	return nil
}

// ExecuteTemplate execute the specified template with the given data.
func (t *Templates) ExecuteTemplate(wr io.Writer, name string, data any) error {
	tmpl, err := t.Base()
	if err != nil {
		return err
	}
	if t.onExecute != nil {
		if err := t.onExecute(wr, tmpl, name, data); err != nil {
			return err
		}
	}
	return tmpl.ExecuteTemplate(wr, name, data)
}
