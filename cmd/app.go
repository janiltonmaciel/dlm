package cmd

import (
	"github.com/urfave/cli"
)

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dfm"
	app.HelpName = app.Name
	app.Usage = "Dockerfile Manager"
	app.UsageText = `
   dfm create                       # Create Dockerfile
   dfm list <language>              # List versions available for docker
       --pre-release                  When listing, show pre-release version
   dfm list <language> <version>    # List versions available for docker, matching a given <version>
       --pre-release                  When listing, show pre-release version
   dfm languages                    # List all supported languages
   dfm --version                    # Print out the installed version of dfm
   dfm --help                       # Show this message`

	app.Commands = createCommands()
	return app
}

func createCommands() []cli.Command {
	return []cli.Command{
		newCommandCreate(),
		newCommandList(),
		newCommandLanguage(),
	}
}
