package type_converter

import (
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"os/exec"

	"github.com/natewong1313/go-react-ssr/config"
	_ "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

// Init starts the type converter
// It gets the name of structs in PropsStructsPath and generates a temporary file to run the type converter
func Init() error {
	// Get struct names from file
	structNames, err := getStructNamesFromFile(config.C.PropsStructsPath)
	if err != nil {
		return err
	}
	// Create a folder for the temporary generator files
	cacheDir, err := utils.GetTypeConverterCacheDir()
	if err != nil {
		return err
	}
	// Create the generator file
	temporaryFilePath, err := createTemporaryFile(cacheDir, structNames)
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
