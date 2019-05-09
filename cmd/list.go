package cmd

import (
	"fmt"
	"strings"

	"github.com/janiltonmaciel/dockerfile-gen/core"
	"github.com/urfave/cli"
	"gopkg.in/gookit/color.v1"
)

type list struct {
	Name      string
	Usage     string
	UsageText string
	Flags     []cli.Flag
}

func NewCommandList() list {
	return list{
		Name:  "ls",
		Usage: "List versions available for docker",
		UsageText: `
    dfm ls <language>                 # List versions available for docker
        --pre-release                   When listing, show pre-release version
    dfm ls <language> <version>       # List versions available for docker, matching a given <version>
        --pre-release                   When listing, show pre-release version

Examples:
   dfm ls python --pre-release        # List versions available for docker python with pre-release
   dfm ls node 8.15                   # List versions available for docker node, matching version 8.15
`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "pre-release",
				Usage: "# Show pre-release versions",
			},
		},
	}
}

func (this list) Action(c *cli.Context) error {
	if c.NArg() <= 0 {
		return this.showCommandHelp(c)
	}

	languageInput := strings.TrimSpace(c.Args().Get(0))
	if languageInput == "" {
		return this.showCommandHelp(c)
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
	versions := core.FindVersions(languageInput, withPrerelease, versionInput)

	for _, version := range versions {
		fmt.Printf("%25s   %s\n",
			color.FgGreen.Render(version.Version),
			color.FgDefault.Render(version.DistributionReleases))
	}
	return nil
}

func (this list) showCommandHelp(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, color.FgRed.Render("Incorrect usage!"))
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, c.Command.Name)
}
