package core

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	loadDistributionConfig()
	languageConfig := loadLanguageConfig()
	setLanguages(languageConfig)
	setVersions(languageChoices)
}

// # Replace shell with bash so we can source files
// RUN rm /bin/sh && ln -s /bin/bash /bin/sh

const (
	DistroDebian = "DEBIAN"
	DistroUbuntu = "UBUNTU"
	DistroAlpine = "ALPINE"
)

const (
	LangNode = "NODE"
)

type (
	Config interface {
		Parse(data []byte) error
	}
	Block struct {
		Description string
		Data        string
		Before      []string
		After       []string
	}

	Context struct {
		From         string
		BeforeBlocks []string
		Blocks       []Block
		AfterBlocks  []string
	}

	AnswerVersion struct {
		Language Language
		Version  Version
	}

	Language struct {
		Name string `yaml:"language"`
		Url  string `yaml:"url"`
	}

	Distribution struct {
		Name          string   `yaml:"name"`
		ReleaseName   string   `yaml:"releaseName"`
		Release       float32  `yaml:"release"`
		Image         string   `yaml:"image"`
		Weight        int      `yaml:"weight"`
		Tags          []string `yaml:"tags"`
		UrlRepository string   `yaml:"urlRepository"`
		UrlDockerfile string   `yaml:"urlDockerfile"`
		Language      Language
	}

	Version struct {
		Version              string         `yaml:"version"`
		MajorVersion         string         `yaml:"majorVersion"`
		Prerelease           bool           `yaml:"prerelease"`
		Date                 string         `yaml:"date"`
		DistributionReleases string         `yaml:"distributionReleases"`
		Distributions        []Distribution `yaml:"distributions"`
	}

	LanguageConfig     []Language
	VersionConfig      []Version
	DistributionConfig []Distribution
)

var (
	languageChoices map[string]Language
	languages       []string
	distributions   []Distribution

	DistributionBuild = map[string]map[string][]string{
		DistroDebian: {
			"Before": []string{
				"# DISABLE GPG IPV6",
				"RUN mkdir ~/.gnupg && echo 'disable-ipv6' >> ~/.gnupg/dirmngr.conf",
			},
			"After": []string{
				`RUN apt-get update && apt-get install -y --no-install-recommends bash \
				&& rm -rf /var/lib/apt/lists/*`,
			},
		},

		DistroUbuntu: {
			"Before": []string{},
			"After":  []string{},
		},

		DistroAlpine: {
			"Before": []string{},
			"After":  []string{},
		},
	}

	LanguageBuild = map[string]map[string][]string{
		LangNode: {
			"Before": []string{
				"ENV PATH node_modules/.bin:$PATH",
			},
			"After": []string{},
		},
	}
)

func (lc *LanguageConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, lc)
}

func (vc *VersionConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, vc)
}

func (dc *DistributionConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, dc)
}

func (d Distribution) Description(languageName string) string {
	desc := fmt.Sprintf(`
##### %s ######
# Official Docker Image for %s
# repository: %s
# dockerfile: %s
# image: %s
	`, strings.ToUpper(languageName), languageName, d.UrlRepository, d.UrlDockerfile, d.Image)

	return desc
}

func NewContext(distributionName string) Context {
	distroBuild := getDistributionBuild(distributionName)
	return Context{
		BeforeBlocks: distroBuild["Before"],
		Blocks:       make([]Block, 0),
		AfterBlocks:  distroBuild["After"],
	}
}

func NewBlock(languageName string, distro Distribution, data string) Block {
	langBuild := getLanguageBuild(languageName)
	block := Block{
		Description: distro.Description(languageName),
		Data:        data,
		Before:      langBuild["Before"],
		After:       langBuild["After"],
	}
	return block
}

func getDistributionBuild(distributionName string) map[string][]string {
	distroBuild, found := DistributionBuild[strings.ToUpper(distributionName)]
	if found {
		return distroBuild
	}
	distroBuild = map[string][]string{
		"Before": []string{},
		"After":  []string{},
	}
	return distroBuild
}

func getLanguageBuild(languageName string) map[string][]string {
	langBuild, found := LanguageBuild[strings.ToUpper(languageName)]
	if found {
		return langBuild
	}

	langBuild = map[string][]string{
		"Before": []string{},
		"After":  []string{},
	}
	return langBuild
}

func GetLanguages() []string {
	return languages
}
func GetDistributions() []Distribution {
	return distributions
}

func GetLanguage(name string) *Language {
	if lang, found := languageChoices[strings.ToLower(name)]; found {
		return &lang
	}
	return nil
}

func loadLanguageConfig() LanguageConfig {
	var languageConfig LanguageConfig
	loadConfig("languages.yml", &languageConfig)
	return languageConfig
}

func loadConfig(fileName string, config Config) {
	data, err := FindBytes(fileName)
	if err != nil {
		panic(err)
	}

	if err := config.Parse(data); err != nil {
		panic(err)
	}
}

func loadDistributionConfig() {
	var dc DistributionConfig
	loadConfig("distributions.yml", &dc)
	distributions = dc
}

func setLanguages(config LanguageConfig) {
	languageChoices = make(map[string]Language)
	for _, lang := range config {
		languageChoices[strings.ToLower(lang.Name)] = lang
		languages = append(languages, strings.Title(lang.Name))
	}
}
