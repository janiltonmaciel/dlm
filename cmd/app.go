package cmd

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/gookit/color.v1"
)

var author = "Janilton Maciel <janilton@gmail.com>"

func CreateApp(version, commit, date string) *cli.App {
	configCli(commit, date)
	configSurvey()
	return createApp(version)
}

func createApp(version string) *cli.App {
	renderGreen := color.FgGreen.Render

	app := cli.NewApp()
	app.Commands = createCommands()
	app.Author = renderGreen(author)
	app.Version = renderGreen(version)
	app.Name = renderGreen("dfm")
	app.HelpName = app.Name
	app.Usage = "Dockerfile Manager"
	app.UsageText = fmt.Sprintf(`
   %s
   %s
   %s
   %s
   %s
   %s
   %s
   %s`,
		fmt.Sprintf("%-48s Create Dockerfile", renderGreen("dfm create")),
		fmt.Sprintf("%-48s List versions available for docker %s", renderGreen("dfm list <language>"), renderYellow("<language>")),
		fmt.Sprintf("%-48s When listing, show %s version", renderGreen("  --pre-release"), renderYellow("pre-release")),
		fmt.Sprintf("%-48s List versions available for docker %s, matching a given %s", renderGreen("dfm list <language> <version>"), renderYellow("<language>"), renderYellow("<version>")),
		fmt.Sprintf("%-48s When listing, show %s version", renderGreen("  --pre-release"), renderYellow("pre-release")),
		fmt.Sprintf("%-48s List all supported languages", renderGreen("dfm languages")),
		fmt.Sprintf("%-48s Print out the installed version of dfm", renderGreen("dfm --version")),
		fmt.Sprintf("%-48s Show this message", renderGreen("dfm --help")),
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
	cli.AppHelpTemplate = AppHelpTemplate
	cli.VersionPrinter = VersionPrinter(commit, date)
	cli.CommandHelpTemplate = CommandHelpTemplate
}

func configSurvey() {
	surveyCore.QuestionIcon = "\n?"
	survey.SelectQuestionTemplate = SelectQuestionTemplate
}
