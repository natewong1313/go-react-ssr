package type_converter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// https://gist.github.com/LukaGiorgadze/570a89a5c3c6d006120da8c29f6684ee
func getStructNamesFromFile(filePath string) (structs []string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(data), 0)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structs = append(structs, ts.Name.Name)
		}
	}
	return
}