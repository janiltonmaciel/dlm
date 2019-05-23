package manager

import (
	"fmt"
	"sort"
	"strings"
)

type (
	Distribution struct {
		Name            string   `yaml:"name"`
		ReleaseName     string   `yaml:"releaseName"`
		Release         float32  `yaml:"release"`
		Image           string   `yaml:"image"`
		Weight          int      `yaml:"weight"`
		Tags            []string `yaml:"tags"`
		UrlRepository   string   `yaml:"urlRepository"`
		UrlDockerfile   string   `yaml:"urlDockerfile"`
		ImageRepository string   `yaml:"imageRepository"`
		Language        Language
	}
)

const (
	DistributionDebian = "DEBIAN"
	DistributionUbuntu = "UBUNTU"
	DistributionALpine = "ALPINE"
)

var distributions []Distribution

func (d Distribution) Description() string {
	desc := fmt.Sprintf(`##### %s ######
# Official Docker Image for %s
# repository: %s
# dockerfile: %s
# image: %s
# tag: %s`,
		strings.ToUpper(d.Language.Name),
		d.Language.Name,
		d.UrlRepository,
		d.UrlDockerfile,
		d.Image,
		d.ImageRepository)

	return desc
}

func (d Distribution) Sort() int {
	return d.Language.SortDistro(d.Name)
}

func (d Distribution) Hash() string {
	return fmt.Sprintf("%s-%s-%s-%s", d.Language.Name, d.Image, d.ImageRepository, d.UrlDockerfile)
}

func (l Language) SortDistro(distributionName string) int {
	return l.sortDistro[strings.ToLower(distributionName)]
}

func GetDistributions() []Distribution {
	return distributions
}

func Intersection(a, b []Distribution) (c []Distribution) {
	m := make(map[string]Distribution)
	for _, item := range a {
		m[item.Name] = item
	}

	for _, item := range b {
		if _, ok := m[item.Name]; ok {
			c = append(c, item)
		}
	}
	return
}

func SanitizeDockerfile(distribution Distribution) (string, error) {
	data, err := GetUrl(distribution.UrlDockerfile)
	if err != nil {
		return "", err
	}

	sanitizeAll := GetLanguageSanitizeDockerfile(distribution.Language.Name)
	newData := make([]string, 0)
	for _, line := range data {
		for _, sanitize := range sanitizeAll {
			if result := sanitize.Pattern.MatchString(line); result {
				line = sanitize.Replace(distribution)
				break
			}
		}
		newData = append(newData, line)
	}
	return strings.Join(newData, "\n"), nil
}

func setDistributions(distros []Distribution) {
	sort.Slice(distros, func(i, j int) bool {
		return distros[i].Sort() < distros[j].Sort()
	})
	distributions = distros
}
