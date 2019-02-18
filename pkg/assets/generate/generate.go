package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("install-wizard/dist/install-wizard")

	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     "pkg/assets/static.go",
		PackageName:  "assets",
		VariableName: "Assets",
		BuildTags:    "!dev",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
