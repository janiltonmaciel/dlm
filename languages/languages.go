package languages

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gobuffalo/packr"
)

var box packr.Box

func init() {
	box = packr.NewBox("../templates")
}

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

func ParseTemplate(context map[string]string, lang string) string {
	var content bytes.Buffer
	templateStr := box.String("Dockerfile-" + lang)
	t := template.Must(template.New("t1" + lang).Option("missingkey=zero").Parse(templateStr))
	if err := t.Execute(&content, context); err != nil {
		fmt.Printf("Error no loadTemplate: %s", err)
		return ""
	}
	return content.String()
}
