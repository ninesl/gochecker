package main

import (
	"context"
	"strings"

	"dagger/goauth/internal/dagger"
)

type Goauth struct {
	Token     *dagger.Secret
	Image     string
	Goprivate string
	// +private
	GoMod *dagger.File
}

func New(
	ws *dagger.Workspace,
	// GitHub PAT with read access to private Go modules
	// +optional
	token *dagger.Secret,
	// Base Go container image (defaults to golang:<version> from go.mod)
	// +optional
	image string,
	// GOPRIVATE glob pattern (e.g. "github.com/org/*")
	// +optional
	goprivate string,
) *Goauth {
	return &Goauth{
		Token:     token,
		Image:     image,
		Goprivate: goprivate,
		GoMod:     ws.File("/go.mod"),
	}
}

// Base returns a Go toolchain container authenticated for private GitHub modules.
func (m *Goauth) Base(ctx context.Context) *dagger.Container {
	image := m.Image
	if image == "" {
		ver := "1.27"
		if m.GoMod != nil {
			if data, err := m.GoMod.Contents(ctx); err == nil {
				for _, line := range strings.Split(data, "\n") {
					fields := strings.Fields(strings.TrimSpace(line))
					if len(fields) >= 2 && fields[0] == "go" && fields[1] != "" {
						ver = fields[1]
						break
					}
				}
			}
		}
		image = "golang:" + ver
	}

	ctr := dag.Container().From(image)
	if m.Goprivate != "" {
		ctr = ctr.WithEnvVariable("GOPRIVATE", m.Goprivate)
	}
	if m.Token == nil {
		return ctr
	}
	askpass := "#!/bin/sh\ncase \"$1\" in\n  *Username*) echo x-access-token ;;\n  *) echo \"$GITHUB_TOKEN\" ;;\nesac\n"
	return ctr.
		WithSecretVariable("GITHUB_TOKEN", m.Token).
		WithNewFile("/usr/local/bin/git-askpass", askpass, dagger.ContainerWithNewFileOpts{Permissions: 0o755}).
		WithEnvVariable("GIT_ASKPASS", "/usr/local/bin/git-askpass").
		WithEnvVariable("GIT_TERMINAL_PROMPT", "0")
}
