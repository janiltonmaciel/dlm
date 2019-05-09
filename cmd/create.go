package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/gookit/color.v1"

	"github.com/janiltonmaciel/dockerfile-gen/core"
)

type create struct {
	Name  string
	Usage string
}

func NewCommandCreate() create {
	return create{
		Name:  "create",
		Usage: "create dockerfile",
	}
}

func (this create) Action(c *cli.Context) error {

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

	distributions, distribution := this.distributionLanguage(answersVersions, answerDistro, c)

	context := core.NewContext(distribution.Name)
	context.From = distribution.Image
	for languageName, distro := range distributions {
		data, err := core.SanitizeDockerfile(distro.UrlDockerfile)
		if err != nil {
			return err
		}
		block := core.NewBlock(languageName, distro, data)
		context.Blocks = append(context.Blocks, block)
	}

	contentDockerfile := core.ParseTemplate(context)
	if contentDockerfile == "" {
		return nil
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
	// core.CheckErr(err)

	return nil
}

func (this create) distributionLanguage(answersVersions []core.AnswerVersion, answerDistro string, c *cli.Context) (map[string]core.Distribution, core.Distribution) {

	distros := make(map[string]core.Distribution)
	for _, av := range answersVersions {
		for _, distro := range av.Version.Distributions {
			if strings.ToLower(distro.Name) == answerDistro {
				if d, found := distros[av.Language.Name]; found {
					if distro.Release > d.Release {
						distros[av.Language.Name] = distro
					}
				} else {
					distros[av.Language.Name] = distro
				}
			}
		}
	}

	// # buildpack-deps:stretch > buildpack-deps:stretch-scm > buildpack-deps:stretch-curl > debian:stretch > debian:stretch-slim

	var distribution core.Distribution
	for _, distro := range distros {
		if distribution.Name == "" || (distro.Release <= distribution.Release && distro.Weight > distribution.Weight) ||
			(distro.Release > distribution.Release && distro.Weight >= distribution.Weight) {
			distribution = distro
		}
	}
	return distros, distribution
}

func (this create) languagesQuestion(c *cli.Context) (answersLanguages []core.Language, err error) {
	prompt := &survey.MultiSelect{
		Message:  "Select the programming languages:",
		Options:  core.GetLanguages(),
		PageSize: 10,
	}
	var answers []string
	err = survey.AskOne(prompt, &answers, survey.Required)

	for _, lang := range answers {
		language := core.GetLanguage(lang)
		if language == nil {
			return answersLanguages, cli.NewExitError("Language not found", 1)
		}
		answersLanguages = append(answersLanguages, *language)
	}

	return
}

func (this create) versionsQuestion(answersLanguages []core.Language, c *cli.Context) (answersVersions []core.AnswerVersion, err error) {
	var version string
	for _, language := range answersLanguages {
		help := fmt.Sprintf(
			"Usage:\n  dfm ls %s   # List versions available for docker %s",
			strings.ToLower(language.Name),
			strings.ToLower(language.Name),
		)
		prompt := &survey.Input{
			Message: fmt.Sprintf("Docker %s version:", color.FgGreen.Render(language.Name)),
			Help:    help,
		}

		version = ""
		err = survey.AskOne(prompt, &version, survey.Required)
		if err != nil {
			return answersVersions, err
		}

		v := core.FindVersion(language.Name, version)
		if v == nil {
			msg := this.messageNotFoundVersion(&language, version)
			fmt.Fprintf(c.App.Writer, msg)
			return answersVersions, cli.NewExitError("", 1)
		}

		answerVersion := core.AnswerVersion{
			Language: language,
			Version:  *v,
		}
		answersVersions = append(answersVersions, answerVersion)
	}

	return answersVersions, nil
}

func (this create) distributionsQuestion(answersVersions []core.AnswerVersion, c *cli.Context) (answerDistro string, err error) {
	distros := this.answersVersionsToDistros(answersVersions)

	promptDistro := &survey.Select{
		Message: "Choose a distribution:",
		Options: distros,
	}
	survey.AskOne(promptDistro, &answerDistro, nil)
	if answerDistro == "" {
		return answerDistro, cli.NewExitError("Choose a distribution!", 1)
	}

	return strings.ToLower(answerDistro), nil
}

func (this create) answersVersionsToDistros(answersVersions []core.AnswerVersion) []string {
	var distributions []core.Distribution
	tam := len(answersVersions)
	if tam == 1 {
		distributions = answersVersions[0].Version.Distributions
	} else {
		distributions = core.Intersection(
			answersVersions[0].Version.Distributions,
			answersVersions[1].Version.Distributions,
		)
		for i := 2; i < tam; i++ {
			distributions = core.Intersection(
				distributions,
				answersVersions[i].Version.Distributions,
			)
		}
	}

	return this.distributionsName(distributions)
}

func (this create) distributionsName(distributions []core.Distribution) (distros []string) {
	exist := make(map[string]bool)
	for _, distro := range distributions {
		if _, ok := exist[distro.Name]; !ok {
			distros = append(distros, distro.Name)
			exist[distro.Name] = true
		}
	}
	return
}

func (this create) messageNotFoundVersion(lang *core.Language, version string) string {
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
	if core.HasDockerfile() {
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
