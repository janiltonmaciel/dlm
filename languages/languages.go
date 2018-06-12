package languages

import "strings"

type Language interface {
	Name() string
	Help() string
	Default() string
	GetDockerfile(version string) string
}

var languages = []Language{
	Node{},
	Python{},
	Ruby{},
	Golang{},
}

var languageChoices []string

func init() {
	for _, lang := range languages {
		languageChoices = append(languageChoices, lang.Name())
	}

}

func GetLanguages() []string {
	return languageChoices
}

func GetLanguage(name string) Language {
	for _, lang := range languages {
		if strings.ToLower(lang.Name()) == strings.ToLower(name) {
			return lang
		}
	}
	return nil
}
