package languages

import "github.com/janiltonmaciel/dockerfile-generator/core"

type Node struct{}

func (n Node) Name() string {
	return "Node"
}
func (n Node) Help() string {
	return "https://nodejs.org/en/download/releases/"
}

func (n Node) Default() string {
	return "10.4.0"
}

func (n Node) GetDockerfile(version string) string {
	context := map[string]string{
		"NODE_VERSION": version,
	}

	return core.ParseTemplate(context, n.Name())
}
