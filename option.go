package tmpls

import (
	"io"
	texttemplate "text/template"
)

// TemplatesOption is the option for configuring [Templates].
type TemplatesOption func(*templatesOptions)

type OnTemplateExecuteFn func(t Template, w io.Writer, data any) error

type templatesOptions struct {
	nocache       bool
	pathSeparator string
	initFn        func() Template
	extensions    map[string]struct{}
	prefixMap     map[string]string

	funcs           FuncMap
	excludeFuncs    []string
	disableBuiltins bool

	preloadFilter func(name string, path string) bool
	onExecute     OnTemplateExecuteFn
}

// WithExtensions configure included template extensions.
// If no extension is passed, all extensions will be included.
//
// By default, the following extensions are included:
//   - .html
//   - .gohtml
//   - .gotxt
func WithExtensions(extensions ...string) TemplatesOption {
	return func(options *templatesOptions) {
		options.extensions = make(map[string]struct{}, len(extensions))
		for _, ext := range extensions {
			options.extensions[ext] = struct{}{}
		}
	}
}

// WithSeparator set the path separator used in the template name.
// Default using dot ('.') as the path separator.
func WithSeparator(separator string) TemplatesOption {
	return func(options *templatesOptions) {
		options.pathSeparator = separator
	}
}

// WithPrefixMap configure prefix mapping for the template name.
//
// For example, WithPrefixMap("components/", "_") will result in all templates in the "components" directory
// will be named as _(name) instead of components.(name).
//
// The prefix always uses / for separating paths.
func WithPrefixMap(keyValues ...string) TemplatesOption {
	return func(options *templatesOptions) {
		pairCnt := len(keyValues) / 2
		options.prefixMap = make(map[string]string, pairCnt)
		for i := 0; i < pairCnt; i++ {
			options.prefixMap[keyValues[i*2]] = keyValues[i*2+1]
		}
	}
}

// WithNocache disable the template cache, making the cache reparse the template on each call.
func WithNocache(nocache bool) TemplatesOption {
	return func(options *templatesOptions) {
		options.nocache = nocache
	}
}

// WithPreloadFilter set a filter function to filter templates that will be preloaded.
// By default, any template whose resolved name starts with an underscore (_) will be ignored.
func WithPreloadFilter(filter func(name string, path string) bool) TemplatesOption {
	return func(options *templatesOptions) {
		options.preloadFilter = filter
	}
}

// WithTextMode replace the underlying implementation with text/template.
func WithTextMode() TemplatesOption {
	return func(options *templatesOptions) {
		options.initFn = func() Template {
			return textTemplate{
				Template: texttemplate.New(""),
			}
		}
	}
}

// WithFuncs set the cached template functions.
func WithFuncs(funcs FuncMap) TemplatesOption {
	return func(options *templatesOptions) {
		options.funcs = funcs
	}
}

// WithoutBuiltinFuncs exclude built-in functions.
// if no function name is passed, all built-in functions will be excluded.
func WithoutBuiltinFuncs(funcNames ...string) TemplatesOption {
	return func(options *templatesOptions) {
		options.excludeFuncs = funcNames
		if len(funcNames) == 0 {
			options.disableBuiltins = true
		}
	}
}

// WithOnExecute set a callback function that runs before the template is executed.
func WithOnExecute(callback OnTemplateExecuteFn) TemplatesOption {
	return func(options *templatesOptions) {
		options.onExecute = callback
	}
}
