package main

import (
	"embed"

	"github.com/kawana77b/tsconfig-template/cmd"
)

var version string = "0.0.0"

//go:embed bases/bases/**
var fs embed.FS

func main() {
	cmd.Version = version
	cmd.TemplateFs = fs
	cmd.Execute()
}
