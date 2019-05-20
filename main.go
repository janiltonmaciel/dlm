package main

import (
	"os"

	"github.com/janiltonmaciel/dockerfile-gen/cmd"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

var (
	version string
	commit  string
	date    string
	author  = "Janilton Maciel <janilton@gmail.com>"
)

func init() {
	cli.AppHelpTemplate = cmd.AppHelpTemplate
	cli.VersionPrinter = cmd.VersionPrinter(commit, date)
	cli.CommandHelpTemplate = cmd.CommandHelpTemplate
	surveyCore.QuestionIcon = "\n?"
	survey.SelectQuestionTemplate = cmd.SelectQuestionTemplate
}

func main() {
	appCmd := cmd.NewCommandApp()
	app := cli.NewApp()
	app.Name = appCmd.HelpName
	app.HelpName = appCmd.HelpName
	app.Usage = appCmd.Usage
	app.UsageText = appCmd.UsageText
	app.Author = author
	app.Version = version

	app.Commands = []cli.Command{
		cmd.NewCommandCreate(),
		cmd.NewCommandList(),
		cmd.NewCommandLanguage(),
	}

	err := app.Run(os.Args)
	if err != nil {
		print(err)
	}

}
