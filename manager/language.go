package manager

import (
	"sort"
	"strconv"
	"strings"
)

type (
	Language struct {
		Name       string   `yaml:"language"`
		Alias      string   `yaml:"alias"`
		Url        string   `yaml:"url"`
		Sort       []string `yaml:"sort"`
		sortDistro map[string]int
	}
)

var (
	languageChoices map[string]Language
	languages       []string
)

func GetLanguages() []string {
	return languages
}

func GetLanguage(alias string) *Language {
	if lang, found := languageChoices[strings.ToLower(alias)]; found {
		return &lang
	}
	return nil
}

func setLanguages(langs []Language) {
	languageChoices = make(map[string]Language)
	for i := 0; i < len(langs); i++ {
		lang := langs[i]
		if len(lang.Sort) > 0 {
			lang.sortDistro = make(map[string]int)
			for _, sort := range lang.Sort {
				data := strings.Split(sort, ":")
				key := data[0]
				value, _ := strconv.Atoi(data[1])
				lang.sortDistro[key] = value
			}
		}

		languageChoices[strings.ToLower(lang.Alias)] = lang
		languages = append(languages, lang.Alias)
	}

	sort.Strings(languages)
}
