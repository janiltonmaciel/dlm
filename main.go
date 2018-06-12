package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/janiltonmaciel/dockerfile-gen/core"
	"github.com/janiltonmaciel/dockerfile-gen/languages"

	"gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

func init() {
	surveyCore.ErrorTemplate = core.ErrorTemplate
	survey.MultiSelectQuestionTemplate = core.MultiSelectQuestionTemplate
	survey.InputQuestionTemplate = core.InputQuestionTemplate
}

func main() {
	answers := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select the programming languages:",
		Options: languages.GetLanguages(),
	}
	survey.AskOne(prompt, &answers, survey.Required)

	var context = make(map[string]string)
	var lang languages.Language
	var version string
	for _, answer := range answers {
		lang = languages.GetLanguage(answer)
		if lang == nil {
			os.Exit(1)
		}

		version = ""
		prompt := &survey.Input{
			Message: fmt.Sprintf("Informe a versão do %s:", lang.Name()),
			Help:    lang.Help(),
			Default: lang.Default(),
		}
		survey.AskOne(prompt, &version, survey.Required)

		context[lang.Name()] = lang.GetDockerfile(version)
	}

	ExtraLibs := "tar make git ca-certificates curl openssh"
	libs := ""
	p := &survey.Input{
		Message: "Add extra libs",
		Help:    "Alpine libs",
		Default: ExtraLibs,
	}
	err := survey.AskOne(p, &libs, nil)
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
