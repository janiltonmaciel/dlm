package cmd

import (
	"fmt"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

var author = "Janilton Maciel <janilton@gmail.com>"

func CreateApp(version, commit, date string) *cli.App {
	configCli(commit, date)
	configSurvey()
	return createApp(version)
}

func createApp(version string) *cli.App {
	app := cli.NewApp()
	app.Commands = createCommands()
	app.Author = manager.RenderGreen(author)
	app.Version = manager.RenderGreen(version)
	app.Name = manager.RenderGreen("dlm")
	app.HelpName = app.Name
	app.Usage = "Dockerfile language Manager"
	app.UsageText = fmt.Sprintf(`
   %s
   %s
   %s
   %s
   %s
   %s
   %s
   %s`,
		fmt.Sprintf("%-48s Create Dockerfile", manager.RenderGreen("dlm create")),
		fmt.Sprintf("%-48s List versions available for docker %s", manager.RenderGreen("dlm list <language>"), manager.RenderYellow("<language>")),
		fmt.Sprintf("%-48s When listing, show %s version", manager.RenderGreen("  --pre-release"), manager.RenderYellow("pre-release")),
		fmt.Sprintf("%-48s List versions available for docker %s, matching a given %s", manager.RenderGreen("dlm list <language> <version>"), manager.RenderYellow("<language>"), manager.RenderYellow("<version>")),
		fmt.Sprintf("%-48s When listing, show %s version", manager.RenderGreen("  --pre-release"), manager.RenderYellow("pre-release")),
		fmt.Sprintf("%-48s List all supported languages", manager.RenderGreen("dlm languages")),
		fmt.Sprintf("%-48s Print out the installed version of dlm", manager.RenderGreen("dlm --version")),
		fmt.Sprintf("%-48s Show this message", manager.RenderGreen("dlm --help")),
	)

	return app
}

func createCommands() []cli.Command {
	return []cli.Command{
		newCommandCreate(),
		newCommandLanguage(),
		newCommandList(),
	}
}

func configCli(commit, date string) {
	cli.AppHelpTemplate = appHelpTemplate
	cli.VersionPrinter = versionPrinter(commit, date)
	cli.CommandHelpTemplate = commandHelpTemplate
}

func configSurvey() {
	surveyCore.QuestionIcon = "\n?"
	survey.SelectQuestionTemplate = selectQuestionTemplate
	survey.ConfirmQuestionTemplate = confirmQuestionTemplate
}
