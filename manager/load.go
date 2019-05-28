package manager

import (
	"fmt"
	"strings"

	"gopkg.in/gookit/color.v1"
)

func init() {
	loadDistributions()
	loadLanguages()
	loadVersions()
	loadContexts()
}

var (
	RenderYellow = color.FgLightYellow.Render
	RenderGreen  = color.FgGreen.Render
	RenderCyan   = color.FgLightCyan.Render
	RenderRed    = color.FgRed.Render
)

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
