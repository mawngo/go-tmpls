package tmpls

import (
	"github.com/mawngo/go-tmpls/cache"
	"github.com/mawngo/go-tmpls/internal"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"text/template"
)

// TemplateCache template caching.
type TemplateCache struct {
	fs      fs.FS
	cache   cache.Cache[*template.Template]
	base    *template.Template
	mu      sync.RWMutex
	nocache bool
}

// TemplateCacheOption is the option for configuring TemplateCache.
type TemplateCacheOption func(*templateCacheOptions)

type templateCacheOptions struct {
	cache           cache.Cache[*template.Template]
	nocache         bool
	excludes        []string
	disableBuiltins bool
	funcs           template.FuncMap
	globs           []string
}

// WithCache configure the underlying cache implementation.
func WithCache(cache cache.Cache[*template.Template]) TemplateCacheOption {
	return func(options *templateCacheOptions) {
		options.cache = cache
	}
}

// WithGlobs set the base template globs.
func WithGlobs(globs ...string) TemplateCacheOption {
	return func(options *templateCacheOptions) {
		options.globs = globs
	}
}

// WithFuncs set the base template functions.
func WithFuncs(funcs template.FuncMap) TemplateCacheOption {
	return func(options *templateCacheOptions) {
		options.funcs = funcs
	}
}

// WithNocache disable the template cache, making the cache reparse the template on each call.
func WithNocache(nocache bool) TemplateCacheOption {
	return func(options *templateCacheOptions) {
		options.nocache = nocache
	}
}

// WithoutBuiltins exclude specific built-in functions.
// if no function name is passed, all built-in functions will be excluded.
func WithoutBuiltins(funcNames ...string) TemplateCacheOption {
	return func(options *templateCacheOptions) {
		options.excludes = funcNames
		if len(funcNames) == 0 {
			options.disableBuiltins = true
		}
	}
}

func NewTemplateCache(fs fs.FS, options ...TemplateCacheOption) (*TemplateCache, error) {
	opt := templateCacheOptions{
		nocache: false,
		cache:   make(cache.MapCache[*template.Template]),
	}
	for _, option := range options {
		option(&opt)
	}
	base := template.New("")
	if !opt.disableBuiltins {
		if funcs := internal.NewBuiltinFuncMap(opt.excludes...); len(funcs) > 0 {
			base = base.Funcs(internal.NewBuiltinFuncMap(opt.excludes...))
		}
	}
	if len(opt.funcs) > 0 {
		base = base.Funcs(opt.funcs)
	}
	if len(opt.globs) > 0 {
		var err error
		base, err = base.ParseFS(fs, opt.globs...)
		if err != nil {
			return nil, err
		}
	}

	return &TemplateCache{
		fs:      fs,
		cache:   opt.cache,
		base:    base,
		nocache: opt.nocache,
	}, nil
}

func MustNewTemplateCache(fs fs.FS, options ...TemplateCacheOption) *TemplateCache {
	c, err := NewTemplateCache(fs, options...)
	if err != nil {
		panic(err)
	}
	return c
}

func (t *TemplateCache) MustParse(file string, globs ...string) *template.Template {
	return template.Must(t.Parse(file, globs...))
}

// Parse will parse the template and cache it.
func (t *TemplateCache) Parse(file string, globs ...string) (*template.Template, error) {
	if t.nocache {
		return t.parse(file, globs...)
	}
	name := file
	if len(globs) > 0 {
		name += ":" + strings.Join(globs, ",")
	}

	t.mu.RLock()
	if tmpl, ok := t.cache.Get(name); ok {
		defer t.mu.RUnlock()
		return tmpl, nil
	}
	t.mu.RUnlock()

	if t.mu.TryLock() {
		defer t.mu.Unlock()
		tmpl, err := t.parse(file, globs...)
		if err != nil {
			return nil, err
		}
		t.cache.Set(name, tmpl)
		return tmpl, nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	// Some other thread should have parsed the template already.
	if tmpl, ok := t.cache.Get(name); ok {
		return tmpl, nil
	}

	// This is weird, should never happen.
	tmpl, err := t.parse(file, globs...)
	if err != nil {
		return nil, err
	}
	t.cache.Set(name, tmpl)
	return tmpl, nil
}

func (t *TemplateCache) parse(file string, globs ...string) (*template.Template, error) {
	clone, err := t.base.Clone()
	if err != nil {
		return nil, err
	}
	if len(globs) == 0 {
		return clone.ParseFS(t.fs, file)
	}
	return template.Must(clone.New(file).ParseFS(t.fs, file)).ParseFS(t.fs, globs...)
}

func (t *TemplateCache) MustClone(file string, globs ...string) *template.Template {
	return template.Must(t.Clone(file, globs...))
}

func (t *TemplateCache) Clone(file string, globs ...string) (*template.Template, error) {
	tmpl, err := t.Parse(file, globs...)
	if err != nil {
		return nil, err
	}
	return tmpl.Clone()
}

func (t *TemplateCache) MustExecute(wr io.Writer, data any, file string, globs ...string) {
	err := t.Execute(wr, data, file, globs...)
	if err != nil {
		panic(err)
	}
}

func (t *TemplateCache) Execute(wr io.Writer, data any, file string, globs ...string) error {
	tmpl, err := t.Parse(file, globs...)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(wr, file, data)
}

func (t *TemplateCache) MustExecuteTemplate(wr io.Writer, name string, data any, file string, globs ...string) {
	err := t.ExecuteTemplate(wr, name, data, file, globs...)
	if err != nil {
		panic(err)
	}
}

func (t *TemplateCache) ExecuteTemplate(wr io.Writer, name string, data any, file string, globs ...string) error {
	tmpl, err := t.Parse(file, globs...)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(wr, name, data)
}

// StandardWebFS setup TemplateCache and http.FileServer based on golang-standards/project-layout,
// which read template from web/template and serve static files static from web/static.
func StandardWebFS(cwd fs.FS, options ...TemplateCacheOption) (*TemplateCache, http.Handler, error) {
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
	templateCache, err := NewTemplateCache(templateFs, options...)
	if err != nil {
		return nil, nil, err
	}
	return templateCache, http.FileServer(http.FS(staticFs)), nil
}

// StandardTemplateFS returns TemplateCache based on golang-standards/project-layout,
// which read template from web/template.
func StandardTemplateFS(cwd fs.FS, options ...TemplateCacheOption) (*TemplateCache, error) {
	fileSystem, err := fs.Sub(cwd, "web")
	if err != nil {
		return nil, err
	}

	templateFs, err := fs.Sub(fileSystem, "template")
	if err != nil {
		return nil, err
	}
	return NewTemplateCache(templateFs, options...)
}
