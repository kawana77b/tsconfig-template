# tsconfig-template

![demo](/docs/images/demo.gif)

A tool to output tsconfig templates on the command line.

## About

This tool selects from `tsconfig.json` hosted in [tsconfig/bases](https://github.com/tsconfig/bases) and outputs to a local working directory.

## Install

Binaries are available on the [Release page](https://github.com/kawana77b/tsconfig-template/releases).

## Development

Use Go 1.27.3 or later and GoReleaser 2.19.0. CI reads the Go version from `go.mod`.

```bash
git submodule update --init --recursive
go test ./...
go vet ./...
go build
goreleaser check
# Build release archives locally without publishing
goreleaser release --snapshot --clean --skip=publish
```

## Usage

```bash
# Select a template and create tsconfig.json in the current work directory
tsconfig-template
# Check available tsconfig types
tsconfig-template list
# Show the installed version
tsconfig-template --version
```

Release binaries report the version embedded by GoReleaser. Builds installed with
`go install github.com/kawana77b/tsconfig-template@vX.Y.Z` report the module version.
Local builds use Go's module version when available, otherwise `dev`. To explicitly
embed a Git tag in a local build:

```bash
go build -ldflags "-X main.version=$(git describe --tags --exact-match)"
```

## Note

The correct use of the original repository is to install the target files from npm and specify them as `extends` in the local `tsconfig.json`.
**In other words, note that this is a sub-optimal approach.**
However, sometimes it is quicker to output the files locally without installing them, such as when you want to quickly launch a project in Typescript.
This tool was created because we wanted to embed the original repository file data into a binary so that we could use it right out of the box as a tool.

## Credits

The following are the repositories used as references.
Other licensing matters are placed in the `CREDITS` file.

- [tsconfig/bases](https://github.com/tsconfig/bases) - Copyright (c) Microsoft Corporation. / MIT License
