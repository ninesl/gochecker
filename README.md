# goauth

A Dagger module that returns a Go toolchain container authenticated to GitHub,
so `dagger check` can download private Go modules (`GOPRIVATE`).

It is meant to be used as the `base` container for `dagger.io/go` and
`dagger.io/go/golangci-lint`.

## Install in a consuming project

```sh
dagger module install github.com/ninesl/gochecker/goauth --name goauth
dagger module install dagger.io/go
dagger module install dagger.io/go/golangci-lint   # optional
```

## Wire it up in `dagger.toml`

`[modules.goauth.settings]` supplies the constructor arguments to `New`, and
`base = "dag://goauth/base"` tells the Go/lint modules to use `goauth`'s
`base` function as their container. The link is resolved with those settings
already applied.

```toml
[modules.goauth]
source = "github.com/ninesl/gochecker/goauth"

[modules.goauth.settings]
# Read from the environment at run time; the raw token is never stored.
token = "env:GITHUB_TOKEN"
# Optional overrides:
# image = "golang:1.27.3"
# goprivate = "github.com/ninesl/*"

[modules.go.settings]
base = "dag://goauth/base"

[modules.golangci-lint.settings]
base = "dag://goauth/base"
```

The same settings can be written with the CLI:

```sh
dagger module settings goauth token env:GITHUB_TOKEN
dagger module settings go base dag://goauth/base
dagger module settings golangci-lint base dag://goauth/base
```

## Run

```sh
export GITHUB_TOKEN=ghp_...   # PAT with read access to the private modules
dagger check
```

`token` is optional. Without it, `base` still returns a working Go container,
just without GitHub credentials, so public-only projects work unchanged.

## What `base` does

- Starts from `image` (default `golang:1.27.3`).
- Sets `GOPRIVATE` (default `github.com/ninesl/*`).
- When `token` is set:
  - exposes it as `GITHUB_TOKEN`,
  - installs a `git-askpass` helper that answers `x-access-token` plus the PAT,
  - sets `GIT_ASKPASS` and `GIT_TERMINAL_PROMPT=0`, so `go mod download` can
    clone private repositories non-interactively.
