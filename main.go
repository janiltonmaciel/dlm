package main

import (
	"flag"
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

var (
	version     string
	commit      string
	date        string
	showVersion bool
)

var logo = `
██████╗  ██████╗  ██████╗██╗  ██╗███████╗██████╗ ███████╗██╗██╗     ███████╗     ██████╗ ███████╗███╗   ██╗
██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗██╔════╝██║██║     ██╔════╝    ██╔════╝ ██╔════╝████╗  ██║
██║  ██║██║   ██║██║     █████╔╝ █████╗  ██████╔╝█████╗  ██║██║     █████╗█████╗██║  ███╗█████╗  ██╔██╗ ██║
██║  ██║██║   ██║██║     ██╔═██╗ ██╔══╝  ██╔══██╗██╔══╝  ██║██║     ██╔══╝╚════╝██║   ██║██╔══╝  ██║╚██╗██║
██████╔╝╚██████╔╝╚██████╗██║  ██╗███████╗██║  ██║██║     ██║███████╗███████╗    ╚██████╔╝███████╗██║ ╚████║
╚═════╝  ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝
`

func initFlags() {
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()
}

func main() {

	initFlags()

	fmt.Println(logo)

	if showVersion {
		printInfo()
		return
	}

	var err error

	answers := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select the programming languages:",
		Options: languages.GetLanguages(),
	}
	err = survey.AskOne(prompt, &answers, survey.Required)
	core.CheckErr(err)

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
			Message: fmt.Sprintf("Release %s:", lang.Name()),
			Help:    lang.Help(),
			Default: lang.Default(),
		}
		err = survey.AskOne(prompt, &version, survey.Required)
		core.CheckErr(err)

		context[lang.Name()] = lang.GetDockerfile(version)
	}

	ExtraLibs := "tar make git ca-certificates curl openssh"
	libs := ""
	p := &survey.Input{
		Message: "Add extra libs",
		Help:    "Alpine libs",
		Default: ExtraLibs,
	}
	err = survey.AskOne(p, &libs, nil)
	core.CheckErr(err)

	context["Libs"] = libs

	contentDockerfile := core.ParseTemplate(context, "Build")
	saveDockerfile(contentDockerfile)
}

func saveDockerfile(content string) {
	output := "Dockerfile"
	rewrite := true
	if core.HasDockerfile() {
		p := &survey.Confirm{
			Message: fmt.Sprintf("Rewrite the file `%s`", output),
			Default: true,
		}
		survey.AskOne(p, &rewrite, nil)
	}

	if rewrite {
		err := ioutil.WriteFile(output, []byte(content), 0644)
		if err == nil {
			fmt.Printf("> Successfully Generated `%s` \n", output)
		} else {
			fmt.Printf("> Fail Generated `%s` \n", output)
		}
	}
}

func printInfo() {
	fmt.Println("Version:", version)
	fmt.Println("Commit:", commit)
	fmt.Println("Date:", date)
}
