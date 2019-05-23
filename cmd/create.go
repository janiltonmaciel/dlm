package cmd

import (
	"fmt"
	"io/ioutil"
	"sort"
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

	distributions, distribution := this.distributionLanguage(answersVersions, answerDistro, c)
	context := manager.NewContext(distributions, distribution)
	println("\n>>>>>>>>> DOCKERFILE >>>>>>>>>")
	var data string
	for index, distro := range distributions {
		isFrom := (index == 0)
		printDistro(distro)
		if !isFrom {
			data, err = manager.SanitizeDockerfile(distro)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return err
			}
		}
		languageBlock := manager.NewLanguageBlock(distro, data, isFrom)
		context.Languages = append(context.Languages, languageBlock)
	}

	contentDockerfile := manager.ParseTemplate(context)
	if contentDockerfile == "" {
		println("contentDockerfile VAZIO")
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
	// manager.CheckErr(err)

	return nil
}

func IntersectionImage(a, b []manager.Distribution, answerDistro string) (c []manager.Distribution) {
	m := make(map[string]manager.Distribution)
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

func distributionLight(distros []manager.Distribution) *manager.Distribution {
	var distribution manager.Distribution
	for _, distro := range distros {
		if distribution.Name == "" ||
			(distro.Release >= distribution.Release && distro.Weight <= distribution.Weight) {
			distribution = distro
		}
	}

	return &distribution
}

func distributionHight(distros []manager.Distribution) *manager.Distribution {
	var distribution manager.Distribution
	for _, distro := range distros {
		if distribution.Name == "" ||
			(distro.Release >= distribution.Release && distro.Weight >= distribution.Weight) {
			distribution = distro
		}
	}

	return &distribution
}

func findDistroByNameRelease(distros []manager.Distribution, distroName string, distroRelease float32) []manager.Distribution {
	data := make([]manager.Distribution, 0)
	for _, d := range distros {
		if d.Name == distroName && d.Release == distroRelease {
			data = append(data, d)
		}
	}
	return data
}

func findDistributionLight(distros []manager.Distribution, distro manager.Distribution) *manager.Distribution {
	data := findDistroByNameRelease(distros, distro.Name, distro.Release)

	if len(data) <= 0 {
		return nil
	}
	return distributionLight(data)
}

func findDistribution2(distros []manager.Distribution, distroName string) *manager.Distribution {
	var distribution manager.Distribution
	for _, d := range distros {
		if strings.ToLower(d.Name) == strings.ToLower(distroName) {
			if distribution.Name == "" ||
				(d.Release >= distribution.Release && d.Weight >= distribution.Weight) {
				distribution = d
			}
		}
	}
	return &distribution
}

func Filter(answersVersions []manager.AnswerVersion, answerDistro string,
	functionFilter func(a, b []manager.Distribution, answerDistro string) (c []manager.Distribution)) []manager.Distribution {

	var distributions []manager.Distribution
	tam := len(answersVersions)
	if tam == 1 {
		distributions = make([]manager.Distribution, 0)
		for _, distro := range answersVersions[0].Version.Distributions {
			if strings.ToLower(distro.Name) == strings.ToLower(answerDistro) {
				distributions = append(distributions, distro)
			}
		}
	} else {
		distributions = functionFilter(
			answersVersions[0].Version.Distributions,
			answersVersions[1].Version.Distributions,
			answerDistro,
		)
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

func printDistros(distros []manager.Distribution) {
	println()
	for _, distro := range distros {
		fmt.Printf("Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", distro.Image, distro.UrlDockerfile, distro.Release, distro.Weight)
	}
}

func printDistro(d manager.Distribution) {
	println()
	// fmt.Printf("FROM %s\n", distro.Image)
	fmt.Printf("Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
}

func printVariable(label string, val interface{}) {
	fmt.Printf("\n%s: %v\n", label, val)
}

func FilterByImage(answersVersions []manager.AnswerVersion, answerDistro string) ([]manager.Distribution, *manager.Distribution) {
	distrosImage := Filter(answersVersions, answerDistro, IntersectionImage)
	if len(distrosImage) <= 0 {
		return []manager.Distribution{}, nil
	}

	distroLight := distributionLight(distrosImage)

	distros := make([]manager.Distribution, 0)
	for _, av := range answersVersions {
		for _, distro := range av.Version.Distributions {
			if distro.Image == distroLight.Image {
				distros = append(distros, distro)
				break
			}
		}
	}

	return distros, &distros[0]
}

func FilterByRelease(answersVersions []manager.AnswerVersion, answerDistro string) (distros []manager.Distribution, distribution *manager.Distribution) {

	for _, distro := range manager.GetDistributions() {
		distrosTemp := make([]manager.Distribution, 0)
		if strings.ToLower(distro.Name) == strings.ToLower(answerDistro) {
			for _, av := range answersVersions {
				if d := findDistributionLight(av.Version.Distributions, distro); d != nil {
					distrosTemp = append(distrosTemp, *d)
				}
			}
		}

		if len(distrosTemp) == len(answersVersions) {
			distroHight := distributionHight(distrosTemp)
			if distribution == nil {
				distribution = distroHight
				distros = distrosTemp
			} else {
				distribution = distributionHight([]manager.Distribution{*distribution, *distroHight})
				if distribution.UrlDockerfile == distroHight.UrlDockerfile {
					distros = distrosTemp
				}
			}
		}
	}

	return
}

func FilterByDistro(answersVersions []manager.AnswerVersion, answerDistro string) (distros []manager.Distribution, distribution *manager.Distribution) {

	distrosTemp := make([]manager.Distribution, 0)
	for _, av := range answersVersions {
		if d := findDistribution2(av.Version.Distributions, answerDistro); d != nil {
			// fmt.Printf("D* Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
			distrosTemp = append(distrosTemp, *d)
		}

		if len(distrosTemp) == len(answersVersions) {
			// fmt.Printf("distrosTemp %+v\n", distrosTemp)
			d := findDistribution2(distrosTemp, answerDistro)
			// fmt.Printf("distrosTemp Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
			if distribution == nil || (d.Release > distribution.Release) {
				distros = distrosTemp
				distribution = d
			}
		}
	}

	return
}

func sortDistributions(distros []manager.Distribution, distro manager.Distribution) ([]manager.Distribution, manager.Distribution) {
	sort.Slice(distros, func(i, j int) bool {
		return distros[i].Sort() > distros[j].Sort()
	})

	distributions := make([]manager.Distribution, 0)
	distributions = append(distributions, distro)
	for _, d := range distros {
		if d.Hash() != distro.Hash() {
			distributions = append(distributions, d)
		}
	}

	return distributions, distro
}

func (this create) distributionLanguage(answersVersions []manager.AnswerVersion, answerDistro string, c *cli.Context) ([]manager.Distribution, manager.Distribution) {
	distros, distribution := FilterByImage(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return sortDistributions(distros, *distribution)
	}
	println("VAZIO: FilterByImage")

	distros, distribution = FilterByRelease(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return sortDistributions(distros, *distribution)
	}
	println("VAZIO: FilterByRelease")

	distros, distribution = FilterByDistro(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return sortDistributions(distros, *distribution)
	}
	println("VAZIO: FilterByDistro")

	return nil, manager.Distribution{}

	// # buildpack-deps:stretch > buildpack-deps:stretch-scm > buildpack-deps:stretch-curl > debian:stretch > debian:stretch-slim
	// for _, distro := range distros {
	// 	if distribution.Name == "" || (distro.Release <= distribution.Release && distro.Weight > distribution.Weight) ||
	// 		(distro.Release > distribution.Release && distro.Weight >= distribution.Weight) {
	// 		distribution = distro
	// 	}
	// }
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
		help := fmt.Sprintf(
			"Usage:\n  dfm list %s              # List versions available for docker %s",
			strings.ToLower(language.Alias),
			language.Alias,
		)

		valueDefault := ""
		switch strings.ToLower(language.Name) {
		case "node":
			valueDefault = "12.2.0"
			// valueDefault = "11.4.0"
		case "python":
			valueDefault = "3.8.0a4"
		case "ruby":
			valueDefault = "2.6.0"
		case "golang":
			valueDefault = "1.12.0"
		case "swift":
			valueDefault = "5.0"
		}

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
