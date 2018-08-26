package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("assets")

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
