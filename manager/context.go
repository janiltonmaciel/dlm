package manager

import (
	"fmt"
	"strings"
)

type (
	Block struct {
		Description string
		Data        string
		Before      []string
		After       []string
	}

	Context struct {
		From      Block
		Before    []string
		Languages []Block
		After     []string
	}

	ContextConfig struct {
		Name    string   `yml:"name"`
		Before  []string `yml:"before"`
		After   []string `yml:"after"`
		Command string   `yml:"command"`
		Libs    []string `yml:"libs"`
	}
)

var (
	distributionContext map[string]ContextConfig
	languageContext     map[string]ContextConfig

	DisableGpgIPV6 = `
# DISABLE GPG IPV6
RUN mkdir -p ~/.gnupg && echo 'disable-ipv6' >> ~/.gnupg/dirmngr.conf`
)

func NewContext(commandLibs string, distros []Distribution, distributionFrom Distribution) Context {
	distroContext := GetDistributionContext(distributionFrom.Name)

	before := make([]string, 0)
	if len(distros) > 1 {
		before = append(before, DisableGpgIPV6)
	}
	before = append(before, distroContext.Before...)

	after := make([]string, 0)
	after = append(after, distroContext.After...)
	if commandLibs != "" {
		after = append(after, commandLibs)
	}

	return Context{
		From: Block{
			Description: distributionFrom.Description(),
			Data:        fmt.Sprintf("FROM %s", distributionFrom.ImageRepository),
		},
		Before:    before,
		Languages: make([]Block, 0),
		After:     after,
	}
}

func NewLanguageBlock(distro Distribution, data string, isFrom bool) Block {
	langContext := getLanguageContext(distro.Language.Name)
	description := distro.Description()
	if isFrom {
		description = ""
	}
	block := Block{
		Description: description,
		Data:        data,
		Before:      langContext.Before,
		After:       langContext.After,
	}
	return block
}

func GetDistributionContext(distributionName string) ContextConfig {
	distroContext, found := distributionContext[strings.ToUpper(distributionName)]
	if found {
		return distroContext
	}

	return ContextConfig{
		Before: []string{},
		After:  []string{},
	}
}

func getLanguageContext(languageName string) ContextConfig {
	langContext, found := languageContext[strings.ToUpper(languageName)]
	if found {
		return langContext
	}

	return ContextConfig{
		Before: []string{},
		After:  []string{},
	}
}

func setDistributionsContext(cc []ContextConfig) {
	distributionContext = make(map[string]ContextConfig)
	for _, context := range cc {
		distributionContext[strings.ToUpper(context.Name)] = context
	}
}

func setLanguagesContext(cc []ContextConfig) {
	languageContext = make(map[string]ContextConfig)
	for _, context := range cc {
		languageContext[strings.ToUpper(context.Name)] = context
	}
}
