package cmd

type app struct {
	Name      string
	HelpName  string
	Usage     string
	UsageText string
}

func NewCommandApp() app {
	return app{
		Name:     "dfm",
		HelpName: "dfm",
		Usage:    "Dockerfile Manager",
		UsageText: `
   dfm create                       # Create Dockerfile
   dfm ls <language>                # List versions available for docker
       --pre-release                  When listing, show pre-release version
   dfm ls <language> <version>      # List versions available for docker, matching a given <version>
       --pre-release                  When listing, show pre-release version
   dfm languages                    # List all supported languages
   dfm --version                    # Print out the installed version of dfm
   dfm --help                       # Show this message`,
	}
}
