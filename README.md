# Go Templates

Template caching and paging utilities. Require go 1.22+

## Installation

```shell
go get -u github.com/mawngo/go-tmpls
```
## Usage

```go
package main

import (
	"embed"
	"flag"
	"github.com/mawngo/go-tmpls/html/tmpls"
	"github.com/mawngo/go-tmpls/page"
	"io/fs"
	"net/http"
	"os"
)

//go:embed web/*
var webFS embed.FS

func main() {
	devmode := flag.Bool("dev", false, "Enable dev mode")
	name := flag.String("name", "World", "Your name")
	addr := flag.String("addr", ":8080", "Server address")
	flag.Parse()

	var root fs.FS = webFS
	if *devmode {
		root = os.DirFS(".")
		println("Dev mode enabled")
	}

	// Setup template cache and http.FileServer from root, 
	// which is embedded when dev mode is disabled.
	// You can use StandardTemplateFS to create the TemplateCache only,
	// or NewTemplateCache(fs, options...) to create template cache 
	// if you want to use another directory for template.
	//
	// StandardWebFS setup TemplateCache and http.FileServer 
	// based on golang-standards/project-layout,
	// which read template from web/template and serve static files from web/static.
	cache, static, err := tmpls.StandardWebFS(root,
		// Disable cache in dev mode,
		// so we can see changes without re-run the project.
		tmpls.WithNocache(*devmode),
		// Include all files in components by default. 
		// Those files can be referenced in using base name.
		tmpls.WithGlobs("components/*.gohtml"))
	if err != nil {
		panic(err)
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", static))
	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		// Paging demonstration, just empty data.
		p := page.NewPage[any](make([]any, page.DefaultPageSize),
			page.DefaultPageSize*10,
			page.NewPaginator(req))

		// Execute template with data.
		// This also sets the Content-Type header to text/html; charset=utf-8.
		cache.MustExecute(res,
			page.D{"Name": *name, "Page": p},
			"index.gohtml",
			"layouts/base.gohtml")
	})

	println("Serving at " + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

```

## Template Caching

Cache the template for re-execution without having to parse it again, support template reload for development.

See [examples](/examples/main.go) for setup and integrating template cache.

### Built-in template functions

By default, this library adds some [helpers](/internal/builtin.go) to the template.
To disable all built-in functions use`WithoutBuiltins()`, or `WithoutBuiltins('fn1', 'fn2', ...)` to disable specific
function.

You can add custom funcs using `WithFuncs`.

### Custom cache

By default, this library uses a map to store all parsed templates, thus make them never expire. If you want expiration,
use `WithCache(impl)` to provide your own `Cache[*template.Template]` implementation.

### No cache

When cache is enabled (default), change to the template that has been parsed will not be visible until you rerun the
project (or the cache expired if you use custom cache implementation).

Use `WithNocache(true)` to disable template cache, force template to parse again on each execution.

## Pagination

This library provides a simple pagination implementation for using in template.
See [page](/page) package and the [example](/examples).