package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/gookit/color.v1"
)

type create struct{}

func newCommandCreate() cli.Command {
	c := create{}
	return cli.Command{
		Name:   "create",
		Usage:  "create dockerfile",
		Action: c.action,
	}
}

func (this create) action(c *cli.Context) error {

	answersLanguages, err := this.languagesQuestion(c)
	if err != nil {
		return err
	}

	answersVersions, err := this.versionsQuestion(answersLanguages, c)
	if err != nil {
		return err
	}

	answerDistro, err := this.distributionsQuestion(answersVersions, c)
	if err != nil {
		return err
	}

	contentDockerfile, err := manager.GetContentDockerfile(answersVersions, answerDistro, c)
	if err != nil {
		return err
	}

	this.saveDockerfile(contentDockerfile)

	// ExtraLibs := "tar make git ca-certificates curl openssh"
	// libs := ""
	// p := &survey.Input{
	// 	Message: "Add extra libs",
	// 	Help:    "Alpine libs",
	// 	Default: ExtraLibs,
	// }
	// err = survey.AskOne(p, &libs, nil)
	// manager.CheckErr(err)

	return nil
}

func (this create) languagesQuestion(c *cli.Context) (answersLanguages []manager.Language, err error) {
	prompt := &survey.MultiSelect{
		Message:  "Select the programming languages:",
		Options:  manager.GetLanguages(),
		PageSize: 10,
	}
	var answers []string
	err = survey.AskOne(prompt, &answers, survey.Required)

	for _, lang := range answers {
		language := manager.GetLanguage(lang)
		if language == nil {
			msg := fmt.Sprintf("Language not found: %s", color.FgLightYellow.Render(lang))
			return answersLanguages, cli.NewExitError(msg, 1)
		}
		answersLanguages = append(answersLanguages, *language)
	}

	return
}

func (this create) versionsQuestion(answersLanguages []manager.Language, c *cli.Context) (answersVersions []manager.AnswerVersion, err error) {
	var version string
	for _, language := range answersLanguages {
		help := fmt.Sprintf("Usage:\n %s %-25s %s %s",
			renderGreen("dfm list"),
			renderGreen(strings.ToLower(language.Alias)),
			"List versions available for docker",
			language.Alias,
		)

		valueDefault := manager.GetValueDefault(language)
		prompt := &survey.Input{
			Message: fmt.Sprintf("Docker %s version:", color.FgGreen.Render(language.Alias)),
			Help:    help,
			Default: valueDefault,
		}

		version = ""
		err = survey.AskOne(prompt, &version, survey.Required)
		if err != nil {
			return answersVersions, err
		}

		v := manager.FindVersion(language.Name, version)
		if v == nil {
			msg := this.messageNotFoundVersion(&language, version)
			fmt.Fprintf(c.App.Writer, msg)
			return answersVersions, cli.NewExitError("", 1)
		}

		answerVersion := manager.AnswerVersion{
			Language: language,
			Version:  *v,
		}
		answersVersions = append(answersVersions, answerVersion)
	}

	return answersVersions, nil
}

func (this create) distributionsQuestion(answersVersions []manager.AnswerVersion, c *cli.Context) (answerDistro string, err error) {
	distros := this.answersVersionsToDistros(answersVersions)
	if len(distros) < 1 {
		fmt.Println()
		for _, av := range answersVersions {
			names := strings.Join(this.distributionsName(av.Version.Distributions), ", ")
			fmt.Printf(
				"Language: %-20v - Version: %-20v- Distributions: %s",
				color.FgGreen.Render(av.Language.Name),
				color.FgGreen.Render(av.Version.Version),
				color.FgLightYellow.Render(names),
			)
			fmt.Println()
		}
		fmt.Println()
		return answerDistro, cli.NewExitError(color.FgRed.Render("Versões de distributions incompativeis!"), 1)
	}

	promptDistro := &survey.Select{
		Message: "Choose a distribution:",
		Options: distros,
	}

	err = survey.AskOne(promptDistro, &answerDistro, nil)
	if err != nil || answerDistro == "" {
		return answerDistro, cli.NewExitError("Choose a distribution!", 1)
	}

	return strings.ToLower(answerDistro), nil
}

func (this create) answersVersionsToDistros(answersVersions []manager.AnswerVersion) []string {
	var distributions []manager.Distribution
	tam := len(answersVersions)
	if tam == 1 {
		distributions = answersVersions[0].Version.Distributions
	} else {
		distributions = manager.Intersection(
			answersVersions[0].Version.Distributions,
			answersVersions[1].Version.Distributions,
		)
		for i := 2; i < tam; i++ {
			distributions = manager.Intersection(
				distributions,
				answersVersions[i].Version.Distributions,
			)
		}
	}

	return this.distributionsName(distributions)
}

func (this create) distributionsName(distributions []manager.Distribution) (distros []string) {
	exist := make(map[string]bool)
	for _, distro := range distributions {
		if _, ok := exist[distro.Name]; !ok {
			distros = append(distros, distro.Name)
			exist[distro.Name] = true
		}
	}
	return
}

func (this create) messageNotFoundVersion(lang *manager.Language, version string) string {
	colorRed := color.FgRed.Render
	colorYellow := color.FgLightYellow.Render
	msg := fmt.Sprintf(
		"  %s %s %s %s %s\n  %s %s %s\n",
		colorRed("Docker"),
		colorRed(lang.Name),
		colorRed("version"),
		fmt.Sprintf(colorYellow("'%s'"), version),
		colorRed("not foud."),
		colorRed("Try"),
		fmt.Sprintf(colorYellow("`dfm ls %s`"), lang.Name),
		colorRed("to browse available versions."),
	)

	return msg
}

func (this create) saveDockerfile(content string) {
	output := "Dockerfile"
	rewrite := false
	if manager.HasDockerfile() {
		p := &survey.Confirm{
			Message: fmt.Sprintf("Rewrite the file `%s`", output),
			Default: true,
		}
		if err := survey.AskOne(p, &rewrite, nil); err != nil {
			return
		}
	} else {
		rewrite = true
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
