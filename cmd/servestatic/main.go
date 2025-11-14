package main

import (
	"flag"
	"fmt"
	"github.com/mawngo/go-tmpls/v2"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	devmode := flag.Bool("dev", false, "Enable dev mode")
	addr := flag.String("addr", ":8080", "Server address")
	// Static directory, relative to the source directory.
	// To reference to files inside the static directory, use the configured name, in this case "static".
	static := flag.String("static", "static", "Static directory")
	// Template directory, relative to the source directory.
	template := flag.String("template", ".", "Template directory")
	extensions := flag.String("text", ".gohtml,.html", "Template file extensions")
	// Prefix mapping, see [tmpls.WithPrefixMap]. Each pair is separated by a coma,
	// each key and value pair is separated by a colon.
	// Example: -prefixmap "_partials/:_" => maps _partials/ to _.
	prefix := flag.String("prefixmap", "", "Prefix mapping")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		println("Arguments required: source directory")
		return
	}
	if len(args) > 1 {
		println("Too many arguments")
		return
	}

	root := os.DirFS(args[0])

	*template = path.Clean(*template)
	templateRoot, err := fs.Sub(root, *template)
	if err != nil {
		println("Invalid template directory", *template, err.Error())
		return
	}

	var staticRoot fs.FS
	if *static != "" {
		*static = path.Clean(*static)
		if *static == "" || *static == "." || strings.HasPrefix(*static, "/") {
			println("Invalid static directory", *static)
		}
		staticRoot, err = fs.Sub(root, *static)
		if err != nil {
			println("Invalid static directory", *static, err.Error())
			return
		}
	}

	var prefixes []string
	if *prefix != "" {
		rawPrefixes := strings.Split(*prefix, ",")
		prefixes = make([]string, 0, len(rawPrefixes)*2)
		for _, rawPrefix := range rawPrefixes {
			kv := strings.SplitN(rawPrefix, ":", 2)
			if len(kv) < 2 {
				println("Invalid prefix mapping: missing value: ", rawPrefix)
			}
			prefixes = append(prefixes, kv[0], kv[1])
		}
	}

	templates, err := tmpls.New(templateRoot,
		tmpls.WithNocache(*devmode),
		tmpls.WithExtensions(strings.Split(*extensions, ",")...),
		tmpls.WithPrefixMap(prefixes...),
		tmpls.WithPreloadMatcher(func(name string, path string) bool {
			if name[0] == '_' {
				return false
			}
			if strings.HasPrefix(path, *static) {
				return false
			}
			return true
		}),
		tmpls.WithOnExecute(func(_ tmpls.Template, w io.Writer, _ any) error {
			if rwr, ok := w.(http.ResponseWriter); ok {
				if rwr.Header().Get("Content-Type") == "" {
					rwr.Header().Set("Content-Type", "text/html; charset=utf-8")
				}
			}
			return nil
		}),
		tmpls.WithPrefixMap(prefixes...),
	)

	if err != nil {
		println("Error initializing", err.Error())
		return
	}

	preloaded, err := templates.Preload()
	if err != nil {
		println("Error preloading templates", err.Error())
		return
	}
	for _, template := range preloaded {
		name := template.Name()
		templatePath := templates.LookupPath(name)
		if name == "" || templatePath == "" {
			continue
		}
		templatePath = strings.TrimSuffix(templatePath, path.Ext(templatePath))

		// Handling root view.
		if templatePath == "index" {
			continue
		}

		templatePath = strings.TrimSuffix(templatePath, "/index")
		fmt.Printf("View [%s] => GET /%s\n", name, templatePath)
		http.HandleFunc("GET /"+templatePath, func(res http.ResponseWriter, req *http.Request) {
			templates.MustExecuteTemplate(res, name, map[string]any{
				"Req": req,
			})
		})
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" || req.URL.Path == "" {
			println("index")
			err := templates.ExecuteTemplate(res, "index", map[string]any{
				"Req": req,
			})
			if err == nil {
				return
			}
		}
		err := templates.ExecuteTemplate(res, "404", map[string]any{
			"Req": req,
		})
		if err == nil {
			return
		}
		http.NotFound(res, req)
	})

	if staticRoot != nil {
		staticPath := strings.TrimSuffix(*static, "/")
		fmt.Printf("Static [%s] => GET /%s/*\n", staticPath, staticPath)
		staticPath = "/" + staticPath + "/"
		http.Handle("GET "+staticPath, http.StripPrefix(staticPath, http.FileServer(http.FS(staticRoot))))
	}

	println("Serving at " + *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
