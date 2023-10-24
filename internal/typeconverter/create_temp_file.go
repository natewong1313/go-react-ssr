package typeconverter

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// https://github.com/tkrajina/typescriptify-golang-structs/blob/master/tscriptify/main.go#L139
func createTemporaryFile(structsFilePath, generatedTypesPath, cacheDir string, structNames []string) (string, error) {
	temporaryFilePath := filepath.ToSlash(filepath.Join(cacheDir, "generator.go"))
	file, err := os.Create(temporaryFilePath)
	if err != nil {
		return temporaryFilePath, err
	}
	defer file.Close()

	t := template.Must(template.New("").Parse(TEMPLATE))

	structsArr := make([]string, 0)
	for _, structName := range structNames {
		structName = strings.TrimSpace(structName)
		if len(structName) > 0 {
			structsArr = append(structsArr, "m."+structName)
		}
	}

	var params TemplateParams
	params.Structs = structsArr

	params.ModuleName, err = getModuleName(structsFilePath)
	if err != nil {
		return temporaryFilePath, err
	}
	params.Interface = true
	params.TargetFile = utils.GetFullFilePath(generatedTypesPath)

	err = t.Execute(file, params)
	if err != nil {
		return temporaryFilePath, err
	}

	return temporaryFilePath, nil
}

// getModuleName gets the module name of the props structs file
func getModuleName(propsStructsPath string) (string, error) {
	dir := filepath.ToSlash(filepath.Dir(utils.GetFullFilePath(propsStructsPath)))
	cmd := exec.Command("go", "list")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
