package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/janiltonmaciel/dockerfile-generator/core"
	"github.com/janiltonmaciel/dockerfile-generator/languages"
	"gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

func init() {
	surveyCore.ErrorTemplate = `{{color "red"}}{{ ErrorIcon }} {{.Error}}{{color "reset"}}
	`
	survey.MultiSelectQuestionTemplate = `
	{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
	{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
	{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
	{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
	{{- else }}
		{{- "  "}}{{- color "cyan"}}{{- if and .Help (not .ShowHelp)}}[ {{ HelpInputRune }} for more help]{{end}}{{color "reset"}}
	  {{- "\n"}}
	  {{- range $ix, $option := .PageEntries}}
		{{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
		{{- if index $.Checked $option}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
		{{- color "reset"}}
		{{- " "}}{{$option}}{{"\n"}}
	  {{- end}}
	{{- end}}`

	survey.InputQuestionTemplate = `
	{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
	{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
	{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
	{{- if .ShowAnswer}}
	  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
	{{- else }}
	  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
	  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
	{{- end}}`

}

type Cli struct {
	languages []languages.Language
	options   []string
}

func NewCli() Cli {
	cli := Cli{}
	cli.languages = languages.Choices

	for _, lang := range cli.languages {
		cli.options = append(cli.options, lang.Name())
	}

	return cli
}

func (c *Cli) GetOptions() []string {
	return c.options
}

func (c Cli) GetLanguage(name string) languages.Language {
	for _, lang := range c.languages {
		if strings.ToLower(lang.Name()) == strings.ToLower(name) {
			return lang
		}
	}
	return nil
}

func main() {
	cli := NewCli()

	answers := []string{}
	prompt := &survey.MultiSelect{
		Message: "Quais as linguagens de sua imagem:",
		Options: cli.GetOptions(),
		// Help:    "Phone number should include the area code",
	}
	validate := func(val interface{}) error {
		if arr, ok := val.([]string); !ok || len(arr) <= 0 {
			return errors.New("Selecione as linguagen(s)")
		}
		return nil
	}
	survey.AskOne(prompt, &answers, validate)

	var context = make(map[string]string)
	var lang languages.Language
	var version string
	for _, answer := range answers {
		lang = cli.GetLanguage(answer)
		if lang == nil {
			os.Exit(1)
		}

		version = ""
		prompt := &survey.Input{
			Message: fmt.Sprintf("Informe a versão do %s:", lang.Name()),
			Help:    lang.Help(),
			Default: lang.Default(),
		}
		survey.AskOne(prompt, &version, nil)

		context[lang.Name()] = lang.GetDockerfile(version)
	}

	ExtraLibs := "tar make git ca-certificates curl openssh"
	libs := ""
	p := &survey.Input{
		Message: "Add extra libs",
		Help:    "Alpine libs",
		Default: ExtraLibs,
	}
	survey.AskOne(p, &libs, nil)
	context["Libs"] = libs

	dockerfile := core.ParseTemplate(context, "Build")

	output := "Dockerfile"

	rewrite := true
	if _, err := os.Stat(output); err == nil {
		p := &survey.Confirm{
			Message: fmt.Sprintf("Rewrite the file `%s`", output),
			Default: true,
		}
		survey.AskOne(p, &rewrite, nil)
	}

	if rewrite {
		err := ioutil.WriteFile(output, []byte(dockerfile), 0644)
		if err == nil {
			fmt.Printf("> Successfully Generated `%s` \n", output)
		} else {
			fmt.Printf("> Fail Generated `%s` \n", output)
		}
	}
}
