package create

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/natewong1313/go-react-ssr/gossr-cli/cmd"
	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Go SSR project",
	Long:  "Create a new Go SSR project",
	Run:   create,
}

func init() {
	cmd.RootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) {
	checkNodeInstalled()
	fmt.Println("Welcome to the creation wizard!")
	projectDir := prompt_getProjectDirectory(args)
	webFramework := prompt_selectWebFramework()
	stylingPlugin := prompt_selectStylingPlugin()
	packageManager := prompt_packageManager()
	checkPackageManagerInstalled(packageManager)
	projectDirExists := utils.CheckPathExists(projectDir)
	if projectDirExists {
		projectDirEmpty := utils.CheckPathEmpty(projectDir)
		if !projectDirEmpty && !prompt_shouldWipeDirectory() {
			os.Exit(0)
		} else {
			wipeDirectory(projectDir)
		}

	} else {
		if err := os.MkdirAll(projectDir, 0777); err != nil {
			utils.HandleError(err)
		}
	}

	bootstrapper := Bootstrapper{
		PackageManager: packageManager,
		ProjectDir:     projectDir,
		WebFramework:   webFramework,
		StylingPlugin:  stylingPlugin,
	}
	bootstrapper.Start()

}

func checkNodeInstalled() bool {
	cmd := exec.Command("node", "-v")
	err := cmd.Run()
	if err != nil {
		logger.L.Error().Msg("Node.js is not installed. Please install Node and try again.")
		os.Exit(1)
	}
	return true
}
func checkPackageManagerInstalled(packageManager string) bool {
	cmd := exec.Command(packageManager, "-v")
	err := cmd.Run()
	if err != nil {
		logger.L.Error().Msg(packageManager + " is not installed. Please install " + packageManager + " and try again.")
		os.Exit(1)
	}
	return true
}
func wipeDirectory(projectDir string) {
	logger.L.Info().Msg("Wiping directory " + projectDir)
	if err := os.RemoveAll(projectDir); err != nil {
		utils.HandleError(err)
	}
	if err := os.MkdirAll(projectDir, 0777); err != nil {
		utils.HandleError(err)
	}
}
