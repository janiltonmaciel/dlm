package cmd

import (
	"fmt"
	"strings"

	"github.com/janiltonmaciel/dockerfile-gen/core"
	"github.com/urfave/cli"
	"gopkg.in/gookit/color.v1"
)

func NewCommandList() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List versions available for docker language",
		UsageText: `
    dfm list <language>                 # List versions available for docker
        --pre-release                   When listing, show pre-release version
    dfm list <language> <version>       # List versions available for docker, matching a given <version>
        --pre-release                   When listing, show pre-release version

Examples:
   dfm list golang --pre-release        # List versions available for docker golang with pre-release
   dfm list python 3.7                  # List versions available for docker python, matching version 3.7
   dfm list python 3 --pre-release      # List versions available for docker python, matching version 3
   dfm list node 8                      # List versions available for docker node, matching version 8.15
`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "pre-release",
				Usage: "# Show pre-release versions",
			},
		},
		Action: listAction,
	}
}

func listAction(c *cli.Context) error {
	if c.NArg() <= 0 {
		return showCommandHelp(c)
	}

	languageInput := strings.TrimSpace(c.Args().Get(0))
	if languageInput == "" {
		return showCommandHelp(c)
	}

	language := core.GetLanguage(languageInput)
	if language == nil {
		msg := fmt.Sprintf("%s %s",
			color.FgRed.Render("Language invalid:"),
			color.FgYellow.Render(languageInput))
		fmt.Fprintln(c.App.Writer, msg)
		fmt.Fprintln(c.App.Writer)
		return cli.ShowCommandHelp(c, "languages")
	}

	versionInput := strings.TrimSpace(c.Args().Get(1))
	withPrerelease := c.Bool("pre-release")
	versions := core.FindVersions(language.Name, withPrerelease, versionInput)

	for _, version := range versions {
		fmt.Printf("%25s   %s\n",
			color.FgGreen.Render(version.Version),
			color.FgDefault.Render(version.DistributionReleases))
	}
	return nil
}

func showCommandHelp(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, color.FgRed.Render("Incorrect usage!"))
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, c.Command.Name)
}
