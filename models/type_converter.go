package models

import "strings"


const TEMPLATE = `package main

import (
	"fmt"

	m "{{ .ModelsPackage }}"
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
	err := t.ConvertToFile("{{ .TargetFile }}")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("OK")
}`

type TemplateParams struct {
	ModelsPackage string
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