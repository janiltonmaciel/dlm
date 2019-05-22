package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/janiltonmaciel/dockerfile-gen/core"
	"gopkg.in/gookit/color.v1"
)

func newCommandLanguage() cli.Command {
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
	fmt.Fprintln(c.App.Writer, color.FgLightYellow.Render("Supported languages:"))
	for _, lang := range languages {
		fmt.Fprintf(c.App.Writer, "    %s", color.FgGreen.Render(strings.ToLower(lang)))
		fmt.Fprintln(c.App.Writer)
	}
	fmt.Println()
	return nil
}
