# tsconfig-template

![demo](/docs/images/demo.gif)

A tool to output tsconfig templates on the command line.

## About

This tool selects from `tsconfig.json` hosted in [tsconfig/bases](https://github.com/tsconfig/bases) and outputs to a local working directory.

## Install

Binaries are available on the [Release page](https://github.com/kawana77b/tsconfig-template/releases).

## Usage

```bash
# Select a template and create tsconfig.json in the current work directory
tsconfig-template
# Check available tsconfig types
tsconfig-template list
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
