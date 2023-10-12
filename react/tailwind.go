package react

import (
	"os/exec"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Run the tailwindcss cli to compile the tailwind css file
func compileTailwindCssFile(filePath string) (string, error) {
	// Uses tailwindcss cli to compile the tailwind css file, takes in the global css file as an input and outputs to the file path passed in
	cmd := exec.Command("npx", "tailwindcss", "-i", utils.GetFullFilePath(config.C.GlobalCSSFilePath), "-o", filePath)
	// Set the working directory to the directory of the tailwind config file
	cmd.Dir = filepath.Dir(utils.GetFullFilePath(config.C.TailwindConfigPath))
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return "", nil
}
