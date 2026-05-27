package main

import (
	"embed"
	"io/fs"
)

//go:embed data
var dataFS embed.FS

func getDataFS() fs.FS {
	sub, err := fs.Sub(dataFS, "data")
	if err != nil {
		panic(err)
	}
	return sub
}
