package typeconverter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func Scan() {
	structs, err := GetStructNamesFromFile("./models/props.go")
	if err != nil {
		panic(err)
	}
	for _, s := range structs {
		fmt.Println(s)
	}
}
// https://gist.github.com/LukaGiorgadze/570a89a5c3c6d006120da8c29f6684ee
func GetStructNamesFromFile(filePath string) (structs []string, err error) {
	data, err := ioutil.ReadFile(filePath)
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




