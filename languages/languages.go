package languages

type Language interface {
	Name() string
	Help() string
	Default() string
	GetDockerfile(version string) string
}

var Choices = []Language{
	Node{},
	Python{},
	Ruby{},
	Golang{},
}
