package go_ssr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/buger/jsonparser"
	"github.com/natewong1313/go-react-ssr/internal/cache"
)

// BuildTailwindCSSFile creates a css file with tailwind classes using the tailwind cli
func (engine *Engine) BuildTailwindCSSFile() error {
	engine.Logger.Debug().Msg("Building tailwind css file")

	tailwindOutPath := filepath.Join(cache.TailwindCacheDir, "tailwind.css")
	cmd := exec.Command("npx", "tailwindcss", "-o", tailwindOutPath)
	// if in production, use the standalone tailwind executable instead of node
	if os.Getenv("APP_ENV") == "production" {
		executableName, err := getTailwindExecutableName()
		if err != nil {
			return err
		}
		executablePath := filepath.Join(cache.TailwindCacheDir, executableName)
		// check if the executable has already been installed
		if _, err := os.Stat(executablePath); os.IsNotExist(err) {
			engine.Logger.Debug().Msgf("Downloading tailwind executable to %s", executablePath)
			if err = engine.downloadTailwindExecutable(executableName, executablePath); err != nil {
				return err
			}
		}
		cmd = exec.Command(executablePath, "-o", tailwindOutPath)
	}
	// Set the working directory to the directory of the tailwind config file
	cmd.Dir = filepath.Dir(filepath.Join(engine.Config.FrontendSrcDir, "../tailwind.config.js"))
	if _, err := cmd.CombinedOutput(); err != nil {
		engine.Logger.Err(err).Msg("Failed to build css file with tailwind")
		return err
	}
	return nil
}

// downloadTailwindExecutable downloads the tailwind executable from github releases (https://tailwindcss.com/blog/standalone-cli)
func (engine *Engine) downloadTailwindExecutable(executableName, executablePath string) error {
	version, err := getLatestTailwindVersion()
	if err != nil {
		return err
	}

	file, err := os.Create(executablePath)
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
	return err
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
	return jsonparser.GetString(respBody, "name")
}

// getTailwindExecutableName detects the tailwind executable download name based on the OS and architecture
func getTailwindExecutableName() (string, error) {
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
	return "", fmt.Errorf("unsupported OS/Architecture: %s/%s", os, arch)
}
