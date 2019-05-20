package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/janiltonmaciel/dockerfile-gen/core"
)

func NewCommandLanguage() cli.Command {
	return cli.Command{
		Name:  "languages",
		Usage: "List all supported languages",
		UsageText: `
   dfm languages                    # List all supported languages\n`,
		Action: languageAction,
	}
}

func languageAction(c *cli.Context) error {
	languages := core.GetLanguages()
	fmt.Println("Supported languages:")
	for _, lang := range languages {
		fmt.Printf(" - %s\n", strings.ToLower(lang))
	}
	fmt.Println()
	return nil
}
