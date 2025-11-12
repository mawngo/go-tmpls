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
	Execute(wr io.Writer, data any) error
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
