package main

import (
	"embed"
	"runtime/debug"

	"github.com/kawana77b/tsconfig-template/cmd"
)

// Set by GoReleaser using -ldflags "-X main.version=...".
var version string

//go:embed bases/bases/**
var fs embed.FS

func main() {
	info, _ := debug.ReadBuildInfo()
	cmd.Version = resolveVersion(version, info)
	cmd.TemplateFs = fs
	cmd.Execute()
}

func resolveVersion(linkerVersion string, info *debug.BuildInfo) string {
	if linkerVersion != "" {
		return linkerVersion
	}
	if info != nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
