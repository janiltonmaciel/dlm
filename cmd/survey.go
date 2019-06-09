package cmd

var (
	selectQuestionTemplate = `
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

	confirmQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- color "green+hb"}}{{if .Default}}(Yes/no) {{else}}(yes/No) {{end}}{{color "reset"}}
{{- end}}`
)
