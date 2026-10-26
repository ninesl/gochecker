package main

import "dagger/goauth/internal/dagger"

type Goauth struct {
	Token     *dagger.Secret
	Image     string
	Goprivate string
}

func New(
	// +optional
	token *dagger.Secret,
	// +optional
	// +default="golang:1.27.3"
	image string,
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

func (m *Goauth) Base() *dagger.Container {
	ctr := dag.Container().From(m.Image).WithEnvVariable("GOPRIVATE", m.Goprivate)
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
