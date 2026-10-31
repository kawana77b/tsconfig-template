---
name: update-project-dependencies
description: Update this project's tsconfig/bases submodule, Go dependencies and toolchain, GoReleaser, GitHub Actions, and development container. Use for upstream synchronization and build environment maintenance, not application behavior fixes.
---

# Update Project Dependencies

Work from the repository root. Apply only the requested updates: a template refresh does not imply a Go upgrade. Preserve any explicit Go compatibility requirement. Determine current stable releases from official sources at execution time rather than reusing version numbers from this skill or previous work.

## Project structure

- `bases/` is a Git submodule pointing to `https://github.com/tsconfig/bases.git`.
- `main.go` embeds `bases/bases/**`. `internal/bases_repo/bases_repo.go` lists JSON files and reads their raw bytes; `recommended.json` is the default template. There is no separate data generation step: rebuilding embeds the updated templates, including JSON with comments.
- `go.mod` and `go.sum` define the Go toolchain requirement and dependencies. Both workflows use `go-version-file: go.mod`.
- `.goreleaser.yaml` defines cross-platform release builds. `.github/workflows/goreleaser.yaml` publishes draft releases on version tags. `.github/workflows/build_test.yaml` validates pull requests and main updates, including a Linux snapshot release build.
- `.devcontainer/devcontainer.json` defines the Go image and GoReleaser Feature. `.devcontainer/devcontainer-lock.json` pins the Feature implementation, not the installed GoReleaser version.
- `.github/dependabot.yml` covers Go modules, GitHub Actions, and Dev Containers. `Makefile` contains submodule initialization and license generation commands.

## Inspect before updating

Check the working branch, root working tree, and submodule working tree. Record the current submodule commit and tool versions. Preserve unrelated changes and do not overwrite local submodule work.

Read the files listed above and `README.md` before editing. Consult official Go downloads, module release information, GitHub Action releases, GoReleaser migration documentation, and Dev Container image manifests for the components being updated. Prefer stable releases; inspect breaking changes before crossing major versions.

## Refresh upstream templates

Initialize the recorded checkout if needed, then check its status before advancing it:

```sh
git submodule update --init --recursive bases
git -C bases status --short
git -C bases rev-parse HEAD
git submodule update --init --remote bases
git -C bases log -1 --format="%h %cs %s"
git diff --submodule=short -- bases
```

`make update` only initializes the pinned checkout; it does not advance to the upstream tip. Verify the upstream default branch instead of assuming its name. Compare the recorded old commit with the new commit inside `bases`, including changes under its `bases/` directory. Summarize additions, removals, and modified templates.

The intended repository change is the submodule pointer. Do not edit, commit, or push changes in the upstream repository. Template content updates normally need no Go code changes; inspect the embedding and enumeration paths if upstream layout changes.

## Update Go and modules

Use the selected Go toolchain for dependency resolution and verification, not merely whatever `go` is on PATH. Check `go version`; `GOTOOLCHAIN` can select a downloaded toolchain without replacing the user's global installation.

```sh
go get -u ./...
go mod tidy
go list -m -u all
```

Set the intended `go` requirement in `go.mod` and remove or update an obsolete `toolchain` directive deliberately. On PowerShell, quote version-bearing arguments to `go mod edit`, such as `'-go=1.27.3'`, to avoid native argument parsing surprises; use the selected release rather than that example version.

Review `go.mod` and `go.sum`. A dependency already at its latest supported release may remain unchanged. Do not switch import paths or replace libraries solely to make every dependency show a newer number. Reconcile upgrades that require a newer Go version with the user's compatibility requirement.

Regenerate `CREDITS` when the Go toolchain or dependencies change, using the existing `make credits` recipe or its equivalent. On Windows, preserve UTF-8 text when capturing `gocredits -skip-missing .` output. Review removed notices as well as additions.

## Keep build environments consistent

- Update GoReleaser's CLI version in both workflows and the Dev Container Feature together. Check the actual local executable's version before validating; an old executable on PATH is not evidence that the new configuration works.
- Migrate `.goreleaser.yaml` to the selected release's schema. For v2, use `version: 2` and plural `formats` in archives and format overrides. Preserve existing archive names, target platforms, and draft release behavior unless the task calls for changing them.
- The `guiyomh` GoReleaser Feature expects a version **without a leading `v`**, such as `2.19.0`. The GitHub Action accepts a tagged version such as `v2.19.0`. Verify the Feature installer if its implementation changes.
- Verify the exact Dev Container image tag exists. Update its Go series consistently with `go.mod`; do not infer image tag prefixes from the previous image. Refresh the Feature lock through Dev Container tooling when changing the Feature implementation rather than hand-editing digests.
- Update GitHub Actions using their official releases and check runtime/runner requirements. Preserve `submodules: recursive` in both checkouts, full history for release metadata, CI read permissions, and release-job write permissions.
- Keep CI's Go version sourced from `go.mod`, its three-OS validation, and its snapshot build without publishing. Maintain Dependabot coverage and update development requirements in `README.md` when versions or commands change.

## Verify the result

Run the applicable checks with the selected tool versions:

```sh
go test ./...
go vet ./...
go mod verify
goreleaser check
goreleaser release --snapshot --clean --skip=publish
git diff --check
```

The snapshot command validates cross-platform binaries and archive creation without publishing. Validate changed workflow files with `actionlint` when available.

After a template or toolchain update, run the built CLI's `list` command and compare its output with the actual JSON filenames under `bases/bases`; derive the count instead of hardcoding it. Run `--default` in a temporary directory and compare `tsconfig.json` byte-for-byte with `bases/bases/recommended.json`. This command writes to the current directory, so do not run it against an existing project configuration.

After any Dev Container change, perform a real `devcontainer build --workspace-folder .` using the available Dev Containers CLI, and inspect the Go/GoReleaser versions in the resulting image. An image manifest lookup alone does not validate Feature installation. If Docker or another required tool is unavailable, state which checks could not run.

## Finish

Review the diff for unintended module, license, lockfile, or template changes. Report old and new versions/commits, template changes, completed checks, and remaining validation limits. Distinguish local checks from workflows actually executed on GitHub.

When asked to commit, use English Conventional Commits messages and keep template synchronization separate from environment updates when both are requested. Examples: `chore: update tsconfig bases to latest upstream` and `chore: update Go toolchain, dependencies, and CI`. Keep application bug fixes outside these maintenance commits. Commit, push, tag, and publish only within the user's requested scope.
