package typeconverter

import (
	"fmt"
	"os/exec"
)

func Init() error {
	structNames, err := getStructNamesFromFile("./models/props.go")
	if err != nil {
		return err
	}
	
	folderPath, err := createCacheFolder()
	if err != nil {
		return err
	}

	temporaryFilePath, err := createTemporaryFile(folderPath, structNames)
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "run", temporaryFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
	// fmt.Println(string(output))
	return nil
}