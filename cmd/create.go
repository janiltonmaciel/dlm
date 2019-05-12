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

	this.distributionLanguage(answersVersions, answerDistro, c)

	// distributions, distribution := this.distributionLanguage(answersVersions, answerDistro, c)

	// context := core.NewContext(distribution.Name)
	// context.From = distribution.Image
	// for _, distro := range distributions {
	// 	data, err := core.SanitizeDockerfile(distro.UrlDockerfile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	block := core.NewBlock("LANG", distro, data)
	// 	context.Blocks = append(context.Blocks, block)
	// }

	// contentDockerfile := core.ParseTemplate(context)
	// if contentDockerfile == "" {
	// 	return nil
	// }
	// this.saveDockerfile(contentDockerfile)

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

func IntersectionRelease(a, b []core.Distribution, answerDistro string) (c []core.Distribution) {
	m := make(map[float32]core.Distribution)
	for _, itemA := range a {
		if strings.ToLower(itemA.Name) != answerDistro {
			continue
		}
		// fmt.Printf("A: Release: %f - image: %s - rep: %s \n", itemA.Release, itemA.Image, itemA.UrlRepository)
		if item, ok := m[itemA.Release]; ok {
			if itemA.Weight < item.Weight {
				m[itemA.Release] = itemA
			}
		} else {
			m[itemA.Release] = itemA
		}
	}
	println()
	for r, v := range m {
		fmt.Printf("A: Release: %f - image: %s - rep: %s \n", r, v.Image, v.UrlRepository)
	}
	println()

	m2 := make(map[float32][]core.Distribution)
	for _, itemB := range b {
		fmt.Printf("B: Release: %f - image: %s - rep: %s \n", itemB.Release, itemB.Image, itemB.UrlRepository)
		if itemA, ok := m[itemB.Release]; ok && strings.ToLower(itemB.Name) == answerDistro {
			if item, ok := m2[itemA.Release]; ok {
				if itemB.Weight < item[0].Weight {
					item[0] = itemB
				}
			} else {
				m2[itemA.Release] = []core.Distribution{itemB, itemA}
			}
		}
	}
	for _, v := range m2 {
		c = append(c, v...)
	}

	return
}

func IntersectionImage(a, b []core.Distribution, answerDistro string) (c []core.Distribution) {
	m := make(map[string]core.Distribution)
	for _, item := range a {
		m[item.Image] = item
	}

	for _, item := range b {
		if _, ok := m[item.Image]; ok && strings.ToLower(item.Name) == answerDistro {
			c = append(c, item)
		}
	}
	return
}

func distributionLight(distros []core.Distribution, answerDistro string) (distribution core.Distribution) {
	for _, distro := range distros {
		if strings.ToLower(distro.Name) == answerDistro {
			if distribution.Name == "" ||
				(distro.Release >= distribution.Release && distro.Weight <= distribution.Weight) {
				distribution = distro
			}
		}
	}
	return
}

func Filter(answersVersions []core.AnswerVersion, answerDistro string,
	functionFilter func(a, b []core.Distribution, answerDistro string) (c []core.Distribution)) []core.Distribution {

	var distributions []core.Distribution
	tam := len(answersVersions)
	if tam == 1 {
		distributions = answersVersions[0].Version.Distributions
	} else {
		distributions = functionFilter(
			answersVersions[0].Version.Distributions,
			answersVersions[1].Version.Distributions,
			answerDistro,
		)
		// printDistros(distributions)
		for i := 2; i < tam; i++ {
			distributions = functionFilter(
				distributions,
				answersVersions[i].Version.Distributions,
				answerDistro,
			)
		}
	}
	return distributions
}

func printDistros(distros []core.Distribution) {
	println()
	for _, distro := range distros {
		fmt.Printf("Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", distro.Image, distro.UrlDockerfile, distro.Release, distro.Weight)
	}
}

func printDistro(distro core.Distribution) {
	println()
	fmt.Printf("FROM %s\n", distro.Image)
}

func FilterByImage(answersVersions []core.AnswerVersion, answerDistro string) (distros []core.Distribution, distribution core.Distribution) {
	distrosImage := Filter(answersVersions, answerDistro, IntersectionImage)
	if len(distrosImage) < 0 {
		return
	}

	distribution = distributionLight(distrosImage, answerDistro)

	distros = make([]core.Distribution, 0)
	for _, av := range answersVersions {
		for _, distro := range av.Version.Distributions {
			if distro.Image == distribution.Image {
				distros = append(distros, distro)
				break
			}
		}
	}
	return
}

func FilterByRelease(answersVersions []core.AnswerVersion, answerDistro string) (distros []core.Distribution, distribution core.Distribution) {
	distrosImage := Filter(answersVersions, answerDistro, IntersectionRelease)
	printDistros(distrosImage)
	if len(distrosImage) < 0 {
		return
	}
	return
}

func (this create) distributionLanguage(answersVersions []core.AnswerVersion, answerDistro string, c *cli.Context) ([]core.Distribution, core.Distribution) {
	distros, distribution := FilterByImage(answersVersions, answerDistro)
	if len(distros) > 0 {
		printDistro(distribution)
		printDistros(distros)
		return distros, distribution
	}
	println("VAZIO: FilterByImage")

	distros, distribution = FilterByRelease(answersVersions, answerDistro)
	if len(distros) > 0 {
		printDistro(distribution)
		printDistros(distros)
		return distros, distribution
	}
	println("VAZIO: FilterByRelease")

	return nil, core.Distribution{}

	// distros := make(map[string][]core.Distribution)
	// for _, av := range answersVersions {
	// 	for _, distro := range av.Version.Distributions {
	// 		if strings.ToLower(distro.Name) == answerDistro {
	// 			if d, found := distros[av.Language.Name]; found {
	// 				distros[av.Language.Name] = append(d, distro)
	// 			} else {
	// 				d := make([]core.Distribution, 0)
	// 				distros[av.Language.Name] = append(d, distro)
	// 			}
	// 		}
	// 	}
	// }

	// for _, distro := range distros {
	// 	if distribution.Name == "" || (distro.Release <= distribution.Release && distro.Weight > distribution.Weight) ||
	// 		(distro.Release > distribution.Release && distro.Weight >= distribution.Weight) {
	// 		distribution = distro
	// 	}
	// }

	// # buildpack-deps:stretch > buildpack-deps:stretch-scm > buildpack-deps:stretch-curl > debian:stretch > debian:stretch-slim

	// var distribution core.Distribution
	// for _, distro := range distros {
	// 	if distribution.Name == "" || (distro.Release <= distribution.Release && distro.Weight > distribution.Weight) ||
	// 		(distro.Release > distribution.Release && distro.Weight >= distribution.Weight) {
	// 		distribution = distro
	// 	}
	// }
	// return distros, distribution
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

		valueDefault := ""
		switch strings.ToLower(language.Name) {
		case "node":
			valueDefault = "12.2.0"
			// valueDefault = "11.4.0"

		case "python":
			valueDefault = "3.7.0a2"
		}

		prompt := &survey.Input{
			Message: fmt.Sprintf("Docker %s version:", color.FgGreen.Render(language.Name)),
			Help:    help,
			Default: valueDefault,
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
