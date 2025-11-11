package tmpls

import (
	html "html/template"
	"io"
	text "text/template"
)

var _ Template = (*htmlTemplate)(nil)
var _ Template = (*textTemplate)(nil)

type FuncMap map[string]any

// Template simple interface for the go std template library.
type Template interface {
	// Unwrap return the underlying template implementation.
	Unwrap() any

	// Name See [html/template.Template.Name].
	Name() string
	// Templates See [html/template.Template.Templates].
	Templates() []Template
	// Lookup See [html/template.Template.Lookup].
	Lookup(name string) Template
	// Funcs See [html/template.Template.Funcs].
	Funcs(funcs FuncMap) Template
	// Clone See [html/template.Template.Clone].
	Clone() (Template, error)
	// New See [html/template.Template.New].
	New(name string) Template

	// Parse See [html/template.Template.Parse].
	Parse(text string) (Template, error)
	// ExecuteTemplate See [html/template.Template.ExecuteTemplate].
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

// htmlTemplate simple wrapper for [html/template.Template].
type htmlTemplate struct {
	*html.Template
}

func (t htmlTemplate) Clone() (Template, error) {
	var err error
	t.Template, err = t.Template.Clone()
	return t, err
}

func (t htmlTemplate) New(name string) Template {
	t.Template = t.Template.New(name)
	return t
}

func (t htmlTemplate) Parse(text string) (Template, error) {
	var err error
	t.Template, err = t.Template.Parse(text)
	return t, err
}

func (t htmlTemplate) Funcs(funcs FuncMap) Template {
	t.Template = t.Template.Funcs(html.FuncMap(funcs))
	return t
}

func (t htmlTemplate) Unwrap() any {
	return t.Template
}

func (t htmlTemplate) Templates() []Template {
	templates := t.Template.Templates()
	res := make([]Template, 0, len(templates))
	for _, v := range templates {
		res = append(res, htmlTemplate{v})
	}
	return res
}

func (t htmlTemplate) Lookup(name string) Template {
	return htmlTemplate{t.Template.Lookup(name)}
}

// textTemplate simple wrapper for [text/template.Template].
type textTemplate struct {
	*text.Template
}

func (t textTemplate) Clone() (Template, error) {
	var err error
	t.Template, err = t.Template.Clone()
	return t, err
}

func (t textTemplate) New(name string) Template {
	t.Template = t.Template.New(name)
	return t
}

func (t textTemplate) Parse(text string) (Template, error) {
	var err error
	t.Template, err = t.Template.Parse(text)
	return t, err
}

func (t textTemplate) Funcs(funcs FuncMap) Template {
	t.Template = t.Template.Funcs(html.FuncMap(funcs))
	return t
}

func (t textTemplate) Unwrap() any {
	return t.Template
}

func (t textTemplate) Templates() []Template {
	templates := t.Template.Templates()
	res := make([]Template, 0, len(templates))
	for _, v := range templates {
		res = append(res, textTemplate{v})
	}
	return res
}

func (t textTemplate) Lookup(name string) Template {
	return textTemplate{t.Template.Lookup(name)}
}
