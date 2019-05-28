package cmd

import (
	"fmt"

	"github.com/janiltonmaciel/dockerfile-gen/manager"
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

var listExamples = fmt.Sprintf(`%-48s %s
   %-48s %s
   %-48s %s
   %-48s %s`,
	manager.RenderGreen("dfm list golang 1.12 --pre-release"), fmt.Sprintf("List versions available for docker %s matching version %s, with %s ", manager.RenderYellow("golang"), manager.RenderYellow("1.12"), manager.RenderYellow("pre-release")),
	manager.RenderGreen("dfm list python 3.7"), fmt.Sprintf("List versions available for docker %s matching version %s", manager.RenderYellow("python"), manager.RenderYellow("3.7")),
	manager.RenderGreen("dfm list ruby"), fmt.Sprintf("List versions available for docker %s", manager.RenderYellow("ruby")),
	manager.RenderGreen("dfm list node 8"), fmt.Sprintf("List versions available for docker %s matching version %s", manager.RenderYellow("node"), manager.RenderYellow("8")),
)

// https://github.com/urfave/cli/blob/master/help.go
// AppHelpTemplate is the text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = fmt.Sprintf(`%s
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

%s
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

%s
   %s
   %s
   %s


{{- if .Version }}

%s
   {{ .Version }}
{{- end}}

{{- if len .Authors}}

%s{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}
{{end}}
`,
	manager.RenderYellow("Name:"),
	manager.RenderYellow("Usage:"),
	manager.RenderYellow("Examples:"),
	fmt.Sprintf("%-48s %s", manager.RenderGreen("dfm create"), "Create Dockerfile"),
	fmt.Sprintf("%-48s %s", manager.RenderGreen("dfm languages"), "List all supported languages"),
	listExamples,
	manager.RenderYellow("Version:"),
	manager.RenderYellow("Author"),
)
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
