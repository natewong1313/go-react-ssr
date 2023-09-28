package type_converter

import "strings"

const TEMPLATE = `package main

import (
	m "{{ .ModuleName }}"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	t := typescriptify.New()
	t.CreateInterface = {{ .Interface }}
	t.BackupDir = ""
{{ range .Structs }}	t.Add({{ . }}{})
{{ end }}
{{ range .CustomImports }}	t.AddImport("{{ . }}")
{{ end }}
	err := t.ConvertToFile(` + "`{{ .TargetFile }}`" + `)
	if err != nil {
		panic(err.Error())
	}
}`

type TemplateParams struct {
	ModuleName    string
	TargetFile    string
	Structs       []string
	InitParams    map[string]interface{}
	CustomImports arrayImports
	Interface     bool
	Verbose       bool
}

type arrayImports []string

func (i *arrayImports) String() string {
	return "// custom imports:\n\n" + strings.Join(*i, "\n")
}

func (i *arrayImports) Set(value string) error {
	*i = append(*i, value)
	return nil
}
