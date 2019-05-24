package manager

import (
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

func GetContentDockerfile(answersVersions []AnswerVersion, answerDistro string, c *cli.Context) (string, error) {
	distributions, distribution := distributionLanguage(answersVersions, answerDistro, c)
	context := NewContext(distributions, distribution)
	println("\n>>>>>>>>> DOCKERFILE >>>>>>>>>")

	var data string
	var err error
	for index, distro := range distributions {
		isFrom := (index == 0)
		// printDistro(distro)
		if !isFrom {
			data, err = SanitizeDockerfile(distro)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return "", err
			}
		}
		languageBlock := NewLanguageBlock(distro, data, isFrom)
		context.Languages = append(context.Languages, languageBlock)
	}

	return ParseTemplate(context), nil
}

func distributionLanguage(answersVersions []AnswerVersion, answerDistro string, c *cli.Context) ([]Distribution, Distribution) {
	distros, distribution := filterByImage(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return appendDistribution(distros, *distribution)
	}
	println("VAZIO: filterByImage")

	distros, distribution = filterByRelease(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return appendDistribution(distros, *distribution)
	}
	println("VAZIO: filterByRelease")

	distros, distribution = filterByDistro(answersVersions, answerDistro)
	if len(distros) > 0 && distribution != nil {
		return appendDistribution(distros, *distribution)
	}
	println("VAZIO: filterByDistro")

	return nil, Distribution{}
}

func filterByImage(answersVersions []AnswerVersion, answerDistro string) ([]Distribution, *Distribution) {
	distrosImage := filter(answersVersions, answerDistro, intersectionImage)
	if len(distrosImage) <= 0 {
		return []Distribution{}, nil
	}

	distroLight := distributionLight(distrosImage)

	distros := make([]Distribution, 0)
	for _, av := range answersVersions {
		for _, distro := range av.Version.Distributions {
			if distro.Image == distroLight.Image {
				distros = append(distros, distro)
				break
			}
		}
	}

	distros = sortDistributions(distros)
	return distros, &distros[0]
}

func filterByRelease(answersVersions []AnswerVersion, answerDistro string) (distros []Distribution, distribution *Distribution) {

	for _, distro := range GetDistributions() {
		distrosTemp := make([]Distribution, 0)
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
				distribution = distributionHight([]Distribution{*distribution, *distroHight})
				if distribution.UrlDockerfile == distroHight.UrlDockerfile {
					distros = distrosTemp
				}
			}
		}
	}

	return
}

func filterByDistro(answersVersions []AnswerVersion, answerDistro string) (distros []Distribution, distribution *Distribution) {

	distrosTemp := make([]Distribution, 0)
	for _, av := range answersVersions {
		if d := findDistributionHight(av.Version.Distributions, answerDistro); d != nil {
			// fmt.Printf("D* Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
			distrosTemp = append(distrosTemp, *d)
		}

		if len(distrosTemp) == len(answersVersions) {
			// fmt.Printf("distrosTemp %+v\n", distrosTemp)
			d := findDistributionHight(distrosTemp, answerDistro)
			// fmt.Printf("distrosTemp Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
			if distribution == nil || (d.Release > distribution.Release) {
				distros = distrosTemp
				distribution = d
			}
		}
	}
	distros = sortDistributions(distros)
	return
}

func findDistributionLight(distros []Distribution, distro Distribution) *Distribution {
	data := findDistroByNameRelease(distros, distro.Name, distro.Release)

	if len(data) <= 0 {
		return nil
	}
	return distributionLight(data)
}

func findDistributionHight(distros []Distribution, distroName string) *Distribution {
	var distribution Distribution
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

func filter(answersVersions []AnswerVersion, answerDistro string,
	functionFilter func(a, b []Distribution, answerDistro string) (c []Distribution)) []Distribution {

	var distributions []Distribution
	tam := len(answersVersions)
	if tam == 1 {
		distributions = make([]Distribution, 0)
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

func sortDistributions(distributions []Distribution) []Distribution {
	sort.Slice(distributions, func(i, j int) bool {
		// fmt.Printf("I) Language: %s - sort:%d\n", distributions[i].Language.Name, distributions[i].Sort())
		// fmt.Printf("j) Language: %s - sort:%d\n\n", distributions[j].Language.Name, distributions[j].Sort())
		return distributions[i].Sort() > distributions[j].Sort()
	})
	return distributions
}

func appendDistribution(distros []Distribution, distro Distribution) ([]Distribution, Distribution) {
	printDistro(distro)
	distributions := make([]Distribution, 0)
	distributions = append(distributions, distro)
	for _, d := range distros {
		if d.Hash() != distro.Hash() {
			distributions = append(distributions, d)
		}
	}

	return distributions, distro
}

func intersectionImage(a, b []Distribution, answerDistro string) (c []Distribution) {
	m := make(map[string]Distribution)
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

func distributionLight(distros []Distribution) *Distribution {
	var distribution Distribution
	for _, distro := range distros {
		if distribution.Name == "" ||
			(distro.Release >= distribution.Release && distro.Weight <= distribution.Weight) {
			distribution = distro
		}
	}

	return &distribution
}

func distributionHight(distros []Distribution) *Distribution {
	distros = sortDistributions(distros)

	var distribution Distribution
	for _, distro := range distros {
		if distribution.Name == "" ||
			(distro.Release >= distribution.Release && distro.Weight >= distribution.Weight) {
			distribution = distro
		}
	}

	return &distribution
}

func findDistroByNameRelease(distros []Distribution, distroName string, distroRelease float32) []Distribution {
	data := make([]Distribution, 0)
	for _, d := range distros {
		if d.Name == distroName && d.Release == distroRelease {
			data = append(data, d)
		}
	}
	return data
}

func printDistros(distros []Distribution) {
	println()
	for _, distro := range distros {
		fmt.Printf("Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", distro.Image, distro.UrlDockerfile, distro.Release, distro.Weight)
	}
}

func printDistro(d Distribution) {
	println()
	// fmt.Printf("FROM %s\n", distro.Image)
	fmt.Printf("Name:%s - Image: %s - Dockerfile:%s - Release: %f - Peso: %d\n", d.Name, d.Image, d.UrlDockerfile, d.Release, d.Weight)
}

func printVariable(label string, val interface{}) {
	fmt.Printf("\n%s: %v\n", label, val)
}
