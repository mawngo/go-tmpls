package main

import (
	"embed"
	"github.com/mawngo/go-tmpls/v2"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed web/*
var webFS embed.FS

func main() {

	var root fs.FS = webFS
	root = os.DirFS(filepath.Join(".", "examples"))
	println(filepath.Abs("."))
	_, err := tmpls.New(root, tmpls.WithExtensions(".gohtml"))
	if err != nil {
		panic(err)
	}
}
