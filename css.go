package go_ssr

import (
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func (engine *Engine) BuildLayoutCSSFile() error {
	if engine.CachedLayoutCSSFilePath == "" && engine.Config.LayoutCSSFilePath != "" {
		layoutCSSCacheDir, err := utils.GetCSSCacheDir()
		if err != nil {
			return err
		}
		cachedCSSFilePath, err := createCachedCSSFile(layoutCSSCacheDir, engine.Config.LayoutCSSFilePath)
		if err != nil {
			return err
		}
		engine.CachedLayoutCSSFilePath = cachedCSSFilePath
	}
	if engine.Config.TailwindConfigPath != "" {
		engine.Logger.Debug().Msg("Building css file with tailwind")
		return engine.buildCSSWithTailwind()
	}
	return nil
}

func createCachedCSSFile(globalCSSCacheDir, globalCSSFilePath string) (string, error) {
	cachedCSSFilePath := utils.GetFullFilePath(filepath.Join(globalCSSCacheDir, "gossr.css"))
	file, err := os.Create(cachedCSSFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	globalCSSFile, err := os.Open(globalCSSFilePath)
	if err != nil {
		return "", err
	}
	defer globalCSSFile.Close()
	_, err = io.Copy(file, globalCSSFile)
	return cachedCSSFilePath, err
}

func (engine *Engine) buildCSSWithTailwind() error {
	// Uses tailwindcss cli to compile the tailwind css file, takes in the global css file as an input and outputs to the file path passed in
	cmd := exec.Command("npx", "tailwindcss", "-i", engine.Config.LayoutCSSFilePath, "-o", engine.CachedLayoutCSSFilePath)
	// Set the working directory to the directory of the tailwind config file
	cmd.Dir = filepath.Dir(engine.Config.TailwindConfigPath)
	_, err := cmd.CombinedOutput()
	return err
}
