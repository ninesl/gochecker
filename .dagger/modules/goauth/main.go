package main

import "dagger/goauth/internal/dagger"

type Goauth struct {
	Token     *dagger.Secret
	Image     string
	Goprivate string
}

func New(
	// GitHub PAT with read access to private Go modules
	token *dagger.Secret,
	// Base Go container image
	// +optional
	// +default="golang:1.27.3"
	image string,
	// GOPRIVATE glob pattern
	// +optional
	// +default="github.com/ninesl/*"
	goprivate string,
) *Goauth {
	return &Goauth{
		Token:     token,
		Image:     image,
		Goprivate: goprivate,
	}
}

// Base returns a Go toolchain container authenticated for private GitHub modules.
func (m *Goauth) Base() *dagger.Container {
	askpass := "#!/bin/sh\ncase \"$1\" in\n  *Username*) echo x-access-token ;;\n  *) echo \"$GITHUB_TOKEN\" ;;\nesac\n"
	return dag.Container().
		From(m.Image).
		WithSecretVariable("GITHUB_TOKEN", m.Token).
		WithNewFile("/usr/local/bin/git-askpass", askpass, dagger.ContainerWithNewFileOpts{
			Permissions: 0o755,
		}).
		WithEnvVariable("GIT_ASKPASS", "/usr/local/bin/git-askpass").
		WithEnvVariable("GIT_TERMINAL_PROMPT", "0").
		WithEnvVariable("GOPRIVATE", m.Goprivate)
}
