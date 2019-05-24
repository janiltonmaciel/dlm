package cmd

import (
	"fmt"
	"strings"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
	"github.com/urfave/cli"
	"gopkg.in/gookit/color.v1"
)

func newCommandList() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List versions available for docker language",
		UsageText: fmt.Sprintf(`
   %s
   %s
   %s
   %s

Examples:
   %s
		`,
			fmt.Sprintf("%-48s List versions available for docker %s", renderGreen("dfm list <language>"), renderYellow("<language>")),
			fmt.Sprintf("%-48s When listing, show %s version", renderGreen("  --pre-release"), renderYellow("pre-release")),
			fmt.Sprintf("%-48s List versions available for docker %s, matching a given %s", renderGreen("dfm list <language> <version>"), renderYellow("<language>"), renderYellow("<version>")),
			fmt.Sprintf("%-48s When listing, show %s version", renderGreen("  --pre-release"), renderYellow("pre-release")),
			listExamples,
		),
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
	if languageInput == "" || strings.ToLower(languageInput) == "help" {
		return showCommandHelp(c)
	}

	language := manager.GetLanguage(languageInput)
	if language == nil {
		return showLanguageCommandHelp(c, languageInput)
	}

	versionInput := strings.TrimSpace(c.Args().Get(1))
	withPrerelease := c.Bool("pre-release")
	printVersions(c, language, withPrerelease, versionInput)

	return nil
}

func showCommandHelp(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, color.FgRed.Render("Incorrect usage!"))
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, c.Command.Name)
}

func showLanguageCommandHelp(c *cli.Context, languageInput string) error {
	msg := fmt.Sprintf("%s %s",
		color.FgRed.Render("Language invalid:"),
		color.FgYellow.Render(languageInput))
	fmt.Fprintln(c.App.Writer, msg)
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, "languages")
}

func printVersions(c *cli.Context, language *manager.Language, withPrerelease bool, versionInput string) {
	versions := manager.FindVersions(language.Name, withPrerelease, versionInput)
	fmt.Fprintf(c.App.Writer, "%s:", color.FgLightYellow.Render(language.Alias))
	fmt.Fprintln(c.App.Writer)

	var versionColor color.Color
	for _, version := range versions {
		versionColor = color.FgGreen
		if version.Prerelease {
			versionColor = color.FgLightCyan
		}

		fmt.Printf("%25s   %s\n",
			versionColor.Render(version.Version),
			color.FgDefault.Render(version.DistributionReleases))
	}
}
