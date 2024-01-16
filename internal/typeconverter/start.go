package typeconverter

import (
	"os/exec"

	"github.com/natewong1313/go-react-ssr/internal/cache"

	_ "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

// Start starts the type converter
// It gets the name of structs in PropsStructsPath and generates a temporary file to run the type converter
func Start(structsFilePath, generatedTypesPath string) error {
	// Get struct names from file
	structNames, err := getStructNamesFromFile(structsFilePath)
	if err != nil {
		return err
	}
	// Create the generator file
	temporaryFilePath, err := createTemporaryFile(structsFilePath, generatedTypesPath, cache.TypeConverterCacheDir, structNames)
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
