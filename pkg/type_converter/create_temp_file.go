package typeconverter

import (
	"gossr/config"
	"gossr/models"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"
)

func createCacheFolder() (string, error) {
	osCacheDir, _ := os.UserCacheDir()
	cacheFolderPath := filepath.Join(osCacheDir, "gossr")
	os.RemoveAll(cacheFolderPath)
	err := os.MkdirAll(cacheFolderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return cacheFolderPath, nil
}
// https://github.com/tkrajina/typescriptify-golang-structs/blob/master/tscriptify/main.go#L139
func createTemporaryFile(folderPath string, structNames []string) (string, error) {
	temporaryFilePath := filepath.Join(folderPath, "generator.go")
	file, err := os.Create(temporaryFilePath)
	if err != nil {
		return temporaryFilePath, err
	}
	defer file.Close()

	t := template.Must(template.New("").Parse(models.TEMPLATE))

	structsArr := make([]string, 0)
	for _, structName := range structNames {
		structName = strings.TrimSpace(structName)
		if len(structName) > 0 {
			structsArr = append(structsArr, "m."+structName)
		}
	}

	var params models.TemplateParams
	params.Structs = structsArr
	
	params.ModelsPackage = getModelsPackageName()
	params.Interface = true
	params.TargetFile = config.Config.Web.GeneratedTypesPath

	err = t.Execute(file, params)
	if err != nil {
		return temporaryFilePath, err
	}
	
	return temporaryFilePath, nil
}

func getModelsPackageName() string {
	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo == nil {
		return ""
	}
	return buildInfo.Main.Path + "/models"
}