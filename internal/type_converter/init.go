package type_converter

import (
	"os/exec"

	"github.com/natewong1313/go-react-ssr/pkg/config"
	_ "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func Init() error {
	// Get struct names from file
	structNames, err := getStructNamesFromFile(config.C.PropsStructsPath)
	if err != nil {
		return err
	}
	// Create a folder for the temporary generator files
	folderPath, err := createCacheFolder()
	if err != nil {
		return err
	}
	// Create the generator file
	temporaryFilePath, err := createTemporaryFile(folderPath, structNames)
	if err != nil {
		return err
	}

	// Run the file
	cmd := exec.Command("go", "run", temporaryFilePath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
