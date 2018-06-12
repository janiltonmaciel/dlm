package core

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

func ParseTemplate(context map[string]string, lang string) string {
	var content bytes.Buffer
	templateStr := GetTemplate("Dockerfile-" + lang)
	t := template.Must(template.New("t1" + lang).Option("missingkey=zero").Parse(templateStr))
	if err := t.Execute(&content, context); err != nil {
		fmt.Printf("Error no loadTemplate: %s", err)
		return ""
	}
	return content.String()
}

func GetTemplate(templateName string) string {
	return box.String(templateName)
}
