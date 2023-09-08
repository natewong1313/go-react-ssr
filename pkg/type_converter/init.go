package type_converter

import (
	"fmt"
	"os/exec"
)

func Init() error {
	// Get struct names from file
	structNames, err := getStructNamesFromFile("./api/models/props.go")
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	return nil
}