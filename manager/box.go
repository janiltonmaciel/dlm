package manager

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gobuffalo/packr"
)

var box packr.Box

func init() {
	box = packr.NewBox("../config")
}

func ParseTemplate(context Context) string {
	var content bytes.Buffer
	templateStr, _ := FindString("Dockerfile-Build")
	t := template.Must(template.New("build").Option("missingkey=zero").Parse(templateStr))
	if err := t.Execute(&content, context); err != nil {
		fmt.Printf("Error no loadTemplate: %s", err)
		return ""
	}
	return content.String()
}

func FindString(fileName string) (string, error) {
	return box.FindString(fileName)
}

func FindBytes(fileName string) ([]byte, error) {
	return box.Find(fileName)
}
