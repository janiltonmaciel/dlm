package main

import (
	"os"

	"github.com/janiltonmaciel/dockerfile-gen/cmd"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	app := cmd.CreateApp(version, commit, date)

	err := app.Run(os.Args)
	if err != nil {
		print(err)
	}
}
