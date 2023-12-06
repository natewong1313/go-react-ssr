package create

import (
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
)

func prompt_getProjectDirectory(args []string) string {
	display := func(path string) string {
		logger.L.Info().Msg("Creating project at " + path)
		return path
	}
	if len(args) > 0 {
		projectDir, err := filepath.Abs(args[0])
		projectDir = filepath.ToSlash(projectDir)
		if err != nil {
			utils.HandleError(err)
		}
		return display(projectDir)
	}

	prompt := promptui.Prompt{
		Label: "Enter the path of your project (leave blank to use current directory)",
	}

	result, err := prompt.Run()
	projectDir, err := filepath.Abs(result)
	if err != nil {
		utils.HandleError(err)
	}
	projectDir = filepath.ToSlash(projectDir)
	return display(projectDir)
}

func prompt_selectWebFramework() string {
	prompt := promptui.Select{
		Label: "Select a web framework to use",
		Items: []string{"Gin", "Fiber", "Echo"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.HandleError(err)
	}
	return result
}

func prompt_selectStylingPlugin() string {
	prompt := promptui.Select{
		Label: "Select a styling plugin to use",
		Items: []string{"None", "Tailwind", "Material UI"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.HandleError(err)
	}
	return result
}

func prompt_isUsingTypescript() bool {
	prompt := promptui.Prompt{
		Label:   "Use Typescript? (y/n)",
		Default: "y",
	}

	result, err := prompt.Run()
	if err != nil {
		utils.HandleError(err)
	}
	return strings.ToLower(result) == "y"
}
func prompt_packageManager() string {
	prompt := promptui.Select{
		Label: "Select a package manager to use",
		Items: []string{"npm", "yarn", "pnpm", "bun"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.HandleError(err)
	}
	return result
}
func prompt_shouldWipeDirectory() bool {
	prompt := promptui.Prompt{
		Label:   "Directory is not empty. Continue? (this will wipe the directory) (y/n)",
		Default: "n",
	}

	result, err := prompt.Run()
	if err != nil {
		utils.HandleError(err)
	}
	return strings.ToLower(result) == "y"
}
