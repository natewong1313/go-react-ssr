package go_ssr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/buger/jsonparser"
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
	cmd := exec.Command("npx", "tailwindcss", "-i", engine.Config.LayoutCSSFilePath, "-o", engine.CachedLayoutCSSFilePath)
	// if in production, use the standalone tailwind executable instead of node
	if os.Getenv("APP_ENV") == "production" {
		executableName, err := detectTailwindDownloadName()
		if err != nil {
			return err
		}
		executableDir, err := utils.GetTailwindExecutableDir()
		if err != nil {
			return err
		}
		executablePath := filepath.Join(executableDir, executableName)
		// check if the executable has already been installed
		if _, err := os.Stat(executablePath); os.IsNotExist(err) {
			engine.Logger.Debug().Msgf("Downloading tailwind executable to %s", executablePath)
			if err = engine.downloadTailwindExecutable(executableName, executableDir); err != nil {
				return err
			}
		}
		cmd = exec.Command(executablePath, "-i", engine.Config.LayoutCSSFilePath, "-o", engine.CachedLayoutCSSFilePath)
	}
	// Set the working directory to the directory of the tailwind config file
	cmd.Dir = filepath.Dir(engine.Config.TailwindConfigPath)
	_, err := cmd.CombinedOutput()
	return err
}

// downloadTailwindExecutable downloads the tailwind executable from github releases (https://tailwindcss.com/blog/standalone-cli)
func (engine *Engine) downloadTailwindExecutable(executableName string, executableDir string) error {
	version, err := getLatestTailwindVersion()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(executableDir, executableName))
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, executableName))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// getLatestTailwindVersion gets the latest tailwind release version from the github api
func getLatestTailwindVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/tailwindlabs/tailwindcss/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	version, err := jsonparser.GetString(respBody, "name")
	return version, err
}

// detectTailwindDownloadName detects the tailwind executable download name based on the OS and architecture
func detectTailwindDownloadName() (string, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	switch os {
	case "darwin":
		switch arch {
		case "arm64":
			return "tailwindcss-macos-arm64", nil
		case "amd64":
			return "tailwindcss-macos-x64", nil
		}
	case "linux":
		switch arch {
		case "arm64":
			return "tailwindcss-linux-arm64", nil
		case "arm":
			return "tailwindcss-linux-armv7", nil
		case "amd64":
			return "tailwindcss-linux-x64", nil
		}
	case "windows":
		switch arch {
		case "arm64":
			return "tailwindcss-windows-arm64.exe", nil
		case "amd64":
			return "tailwindcss-windows-x64.exe", nil
		}
	}
	return "", errors.New("unsupported OS/Architecture")
}
