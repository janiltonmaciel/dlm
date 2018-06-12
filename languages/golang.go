package languages

import "github.com/janiltonmaciel/dockerfile-generator/core"

type Golang struct{}

func (g Golang) Name() string {
	return "Golang"
}

func (g Golang) Help() string {
	return "https://golang.org/doc/devel/release.html"
}

func (g Golang) Default() string {
	return "1.10.3"
}

func (g Golang) GetDockerfile(version string) string {
	context := map[string]string{
		"GOLANG_VERSION": version,
	}

	return core.ParseTemplate(context, g.Name())
}
