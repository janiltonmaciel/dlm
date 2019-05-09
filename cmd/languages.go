package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/janiltonmaciel/dockerfile-gen/core"
)

type language struct {
	Name      string
	Usage     string
	UsageText string
}

func NewCommandLanguage() language {
	return language{
		Name:  "languages",
		Usage: "List all supported languages",
		UsageText: `
   dfm languages                    # List all supported languages
`,
	}
}

func (this language) Action(c *cli.Context) error {
	languages := core.GetLanguages()
	fmt.Println("Supported languages:")
	for _, lang := range languages {
		fmt.Printf(" - %s\n", strings.ToLower(lang))
	}
	fmt.Println()
	return nil
}
