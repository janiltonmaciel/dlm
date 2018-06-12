package languages

import "strings"

type Python struct {
}

func (n Python) Name() string {
	return "Python"
}

func (n Python) Help() string {
	return "https://www.python.org/downloads/"
}

func (n Python) Default() string {
	return "3.6.5"
}

func (p Python) GetDockerfile(version string) string {
	context := map[string]string{
		"PYTHON_VERSION": version,
	}

	pythonMajor := "3"
	if strings.HasPrefix(version, "2") {
		pythonMajor = "2"
	}

	return ParseTemplate(context, p.Name()+pythonMajor)
}
