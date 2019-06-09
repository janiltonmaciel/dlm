package manager

import (
	"fmt"
	"strings"

	"gopkg.in/gookit/color.v1"
)

var (
	ProjectBaseUrl = "https://raw.githubusercontent.com/janiltonmaciel/dlm/master"

	RenderYellow = color.FgLightYellow.Render
	RenderGreen  = color.FgGreen.Render
	RenderCyan   = color.FgLightCyan.Render
	RenderRed    = color.FgRed.Render
)

func init() {
	loadDistributions()
	loadLanguages()
	loadVersions()
	loadContexts()

	go refresh()
}

func loadDistributions() {
	var d []Distribution
	loadConfig("distributions.yml", &d)
	setDistributions(d)
}

func loadLanguages() {
	var l []Language
	loadConfig("languages.yml", &l)
	setLanguages(l)
}

func loadVersions() {
	for _, langName := range GetLanguages() {
		lang := GetLanguage(langName)
		fileNameVersions := fmt.Sprintf("versions/%s.yml", strings.ToLower(langName))
		var v []Version
		loadConfig(fileNameVersions, &v)
		setVersions(*lang, v)
	}
}

func loadContexts() {
	var ccDistros []ContextConfig
	loadConfig("distributions-context.yml", &ccDistros)
	setDistributionsContext(ccDistros)

	var ccLanguages []ContextConfig
	loadConfig("languages-context.yml", &ccLanguages)
	setLanguagesContext(ccLanguages)
}

func refresh() {
	var url string
	for _, langName := range GetLanguages() {
		lang := GetLanguage(langName)
		url = fmt.Sprintf("%s/config/versions/%s.yml", ProjectBaseUrl, strings.ToLower(langName))

		var v []Version
		loadConfigUrl(url, &v)
		if v != nil {
			setVersions(*lang, v)
		}
	}

	var ccDistros []ContextConfig
	url = fmt.Sprintf("%s/distributions-context.yml", ProjectBaseUrl)
	loadConfigUrl(url, &ccDistros)
	if ccDistros != nil {
		setDistributionsContext(ccDistros)
	}

	var ccLanguages []ContextConfig
	url = fmt.Sprintf("%s/config/languages-context.yml", ProjectBaseUrl)
	loadConfigUrl(url, &ccLanguages)
	if ccLanguages != nil {
		setLanguagesContext(ccLanguages)
	}
}
