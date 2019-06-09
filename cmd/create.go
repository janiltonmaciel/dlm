package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
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

func (this create) action(c *cli.Context) (err error) {

	var answersLanguages []manager.Language
	if answersLanguages, err = this.languagesQuestion(c); err != nil {
		return err
	}

	var answersVersions []manager.AnswerVersion
	if answersVersions, err = this.versionsQuestion(answersLanguages, c); err != nil {
		return err
	}

	var answerDistro string
	if answerDistro, err = this.distributionsQuestion(answersVersions, c); err != nil {
		return err
	}

	var commandLibs string
	if commandLibs, err = this.libsQuestion(answerDistro, c); err != nil {
		return err
	}

	var contentDockerfile string
	if contentDockerfile, err = manager.GetContentDockerfile(commandLibs, answersVersions, answerDistro, c); err != nil {
		return err
	}

	return this.saveDockerfile(contentDockerfile, c)
}

func (this create) languagesQuestion(c *cli.Context) (answersLanguages []manager.Language, err error) {
	prompt := &survey.MultiSelect{
		Message:  "Select the programming languages:",
		Options:  manager.GetLanguages(),
		PageSize: 10,
	}
	var answers []string
	err = survey.AskOne(prompt, &answers, survey.Required)
	if err != nil {
		return
	}

	for _, lang := range answers {
		language := manager.GetLanguage(lang)
		if language == nil {
			msg := fmt.Sprintf("Language not found: %s", manager.RenderYellow(lang))
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
			manager.RenderGreen("dlm list"),
			manager.RenderGreen(strings.ToLower(language.Alias)),
			"List versions available for docker",
			language.Alias,
		)

		valueDefault := manager.GetValueDefault(language)
		prompt := &survey.Input{
			Message: fmt.Sprintf("Docker %s version:", manager.RenderGreen(language.Alias)),
			Help:    help,
			Default: valueDefault,
		}

		version = ""
		q := &survey.Question{
			Prompt:    prompt,
			Validate:  survey.Required,
			Transform: manager.TransformFindVersion(language.Name),
		}
		if err = survey.Ask([]*survey.Question{q}, &version); err != nil {
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
		fmt.Fprintln(c.App.Writer)
		for _, av := range answersVersions {
			names := strings.Join(this.distributionsName(av.Version.Distributions), ", ")
			fmt.Fprintf(c.App.Writer,
				"Language: %-20v - Version: %-20v- Distributions: %s",
				manager.RenderGreen(av.Language.Name),
				manager.RenderGreen(av.Version.Version),
				manager.RenderYellow(names),
			)
			fmt.Fprintln(c.App.Writer)
		}
		fmt.Fprintln(c.App.Writer)
		return answerDistro, cli.NewExitError(manager.RenderRed("X Incompatible distributions"), 1)
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

func (this create) libsQuestion(answerDistro string, c *cli.Context) (string, error) {
	distroContext := manager.GetDistributionContext(answerDistro)
	if len(distroContext.Libs) == 0 {
		return "", nil
	}

	p := &survey.Confirm{
		Message: "Add extra libs:",
		Default: false,
	}

	var addLibs bool
	if err := survey.AskOne(p, &addLibs, nil); err != nil {
		return "", err
	}

	if addLibs {
		prompt := &survey.MultiSelect{
			Message:  "Select libs:",
			Options:  distroContext.Libs,
			PageSize: 10,
		}

		var answers []string
		if err := survey.AskOne(prompt, &answers, survey.Required); err != nil {
			return "", err
		}

		cmd := fmt.Sprintf(distroContext.Command, strings.Join(answers, " "))
		return cmd, nil
	}

	return "", nil
}

func (this create) saveDockerfile(content string, c *cli.Context) error {
	output := "Dockerfile"
	rewrite := false
	if manager.HasDockerfile() {
		p := &survey.Confirm{
			Message: fmt.Sprintf("Rewrite the file `%s`:", output),
			Default: false,
		}
		if err := survey.AskOne(p, &rewrite, nil); err != nil {
			return err
		}
	} else {
		rewrite = true
	}

	fmt.Fprintln(c.App.Writer)
	if rewrite {
		err := ioutil.WriteFile(output, []byte(content), 0644)
		if err == nil {
			fmt.Fprintf(c.App.Writer,
				"%s `%s` \n",
				manager.RenderCyan("⚡️ Successfully generated"),
				manager.RenderGreen(output))
		} else {
			fmt.Fprintf(c.App.Writer,
				"%s `%s` \n",
				manager.RenderRed("X Fail generated"),
				manager.RenderYellow(output))
		}
	}
	fmt.Fprintln(c.App.Writer)

	return nil
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
	msg := fmt.Sprintf(
		"  %s %s %s %s %s\n  %s %s %s\n",
		manager.RenderRed("Docker"),
		manager.RenderRed(lang.Name),
		manager.RenderRed("version"),
		fmt.Sprintf(manager.RenderYellow("'%s'"), version),
		manager.RenderRed("not foud."),
		manager.RenderRed("Try"),
		fmt.Sprintf(manager.RenderYellow("`dlm ls %s`"), strings.ToLower(lang.Name)),
		manager.RenderRed("to browse available versions."),
	)

	return msg
}
