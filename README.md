# goauth

Provides a GitHub-authenticated Go base `Container` (`dag://goauth/base`) for `dagger.io/go`, `dagger.io/go/golangci-lint`, and `dagger.io/go/staticcheck` so `dagger check` can fetch private Go modules.

It reads the Go version directly from the consuming workspace's `go.mod` (`golang:<version>`) unless `image` is overridden.

## Usage in a Go repo (`dagger.toml`)

```toml
[modules.go]
source = "dagger.io/go"

[modules.go.settings]
base = "dag://goauth/base"

[modules.golangci-lint]
source = "dagger.io/go/golangci-lint"

[modules.golangci-lint.settings]
base = "dag://goauth/base"

[modules.goauth]
source = "github.com/ninesl/gochecker/goauth"

[modules.goauth.settings]
goprivate = "github.com/ninesl/*"
token = "cmd:gh auth token" # or "env:GITHUB_TOKEN" / "file:/path/to/pat"
```

Then run:

```sh
dagger check
```
