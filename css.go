package go_ssr

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// BuildLayoutCSSFile builds the layout css file if it exists
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

// createCachedCSSFile creates a cached css file from the layout css file
func createCachedCSSFile(layoutCSSCacheDir, layoutCSSFilePath string) (string, error) {
	cachedCSSFilePath := utils.GetFullFilePath(filepath.Join(layoutCSSCacheDir, "gossr.css"))
	file, err := os.Create(cachedCSSFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	globalCSSFile, err := os.Open(layoutCSSFilePath)
	if err != nil {
		return "", err
	}
	defer globalCSSFile.Close()
	_, err = io.Copy(file, globalCSSFile)
	return cachedCSSFilePath, err
}

// buildCSSWithTailwind builds the css file with tailwind cli
func (engine *Engine) buildCSSWithTailwind() error {
	command := []string{"npx"}
	if engine.Config.JSRuntime == "bun" {
		command = []string{"bunx"}
	}
	if engine.Config.JSRuntime == "yarn" {
		command = []string{"yarn", "exec"}
	}
	if engine.Config.JSRuntime == "pnpm" {
		command = []string{"pnpm", "exec"}
	}
	command = append(command, "tailwindcss", "-i", engine.Config.LayoutCSSFilePath, "-o", engine.CachedLayoutCSSFilePath)
	cmd := exec.Command(command[0], command[1:]...)
	// Set the working directory to the directory of the tailwind config file
	cmd.Dir = filepath.Dir(engine.Config.TailwindConfigPath)
	_, err := cmd.CombinedOutput()
	return err
}
