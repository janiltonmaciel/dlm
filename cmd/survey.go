package cmd

var (
	ErrorTemplate = `{{color "red"}}{{ ErrorIcon }} {{.Error}}{{color "reset"}}
`
	MultiSelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
	{{- "  "}}{{- color "cyan"}}{{- if and .Help (not .ShowHelp)}}[ {{ HelpInputRune }} for more help]{{end}}{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
	{{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
	{{- if index $.Checked $option}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
	{{- color "reset"}}
	{{- " "}}{{$option}}{{"\n"}}
  {{- end}}
{{- end}}`

	InputQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
{{- end}}`
)
