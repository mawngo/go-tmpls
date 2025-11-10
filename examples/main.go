package main

import (
	"embed"
	"flag"
	"github.com/mawngo/go-tmpls/v2"
	"github.com/mawngo/go-tmpls/v2/page"
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
	// or NewTemplateCache(fs, options...) to create a template cache
	// if you want to use another directory for template.
	//
	// StandardWebFS set up TemplateCache and http.FileServer
	// based on golang-standards/project-layout,
	// which read template from web/template and serve static files from web/static.
	cache, static, err := tmpls.NewStandardWebFS(root,
		// Disable cache in dev mode, so we can see changes without re-run the project.
		tmpls.WithNocache(*devmode))
	if err != nil {
		panic(err)
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", static))
	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		// Paging demonstration, just empty data.
		p := page.NewPage[any](
			page.NewPaging(req.URL),
			make([]any, page.DefaultPageSize),
			page.DefaultPageSize*10,
		)

		// Execute template with data.
		// This also sets the Content-Type header to text/html; charset=utf-8.
		cache.MustExecuteTemplate(res,
			"index", page.D{"Name": *name, "Page": p})
	})

	println("Serving at " + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
