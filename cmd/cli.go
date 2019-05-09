package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func VersionPrinter(commit, date string) func(c *cli.Context) {
	return func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "version: %s\n", c.App.Version)
		fmt.Fprintf(c.App.Writer, "author: %s\n", c.App.Author)
		fmt.Fprintf(c.App.Writer, "commit: %s\n", commit)
		fmt.Fprintf(c.App.Writer, "date: %s\n", date)
	}
}

// https://github.com/urfave/cli/blob/master/help.go
// AppHelpTemplate is the text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = `Name:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

Usage:
{{- if .UsageText }}{{.UsageText}}
{{- else }}
	{{.HelpName}}
	{{if .VisibleFlags}}
		[global options]
	{{end}}
	{{if .Commands}}
		command [command options]
	{{end}}
	{{if .ArgsUsage}}
		{{.ArgsUsage}}
	{{else}}
		[arguments...]
	{{end}}
{{- end}}

Examples:
   dfm create                       # Create Dockerfile
   dfm ls node                      # List versions available for docker node
   dfm ls python --pre-release      # List versions available for docker python with pre-release
   dfm ls node 8.15                 # List versions available for docker node, matching version 8.15
   dfm languages                    # List all supported languages


{{- if .Version }}

Version:
   {{ .Version }}
{{- end}}

{{- if len .Authors}}

Author{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}
{{end}}
`
var CommandHelpTemplate = `Name:
   {{.HelpName}} - {{.Usage}}

Usage:
   {{- if .UsageText }}{{ .UsageText}}{{- else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
`

var SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- "  "}}{{- color "cyan"}}[Use arrows to move, enter to select, type to filter{{- if and .Help (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
    {{- $choice}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`
