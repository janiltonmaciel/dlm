package cmd

import (
	"fmt"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
	"github.com/urfave/cli"
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
	languages := manager.GetLanguages()
	fmt.Fprintln(c.App.Writer, manager.RenderYellow("Supported languages:"))
	for _, lang := range languages {
		fmt.Fprintf(c.App.Writer, "    %s", manager.RenderGreen(lang))
		fmt.Fprintln(c.App.Writer)
	}
	fmt.Fprintln(c.App.Writer)
	return nil
}
