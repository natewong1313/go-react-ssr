package update

import (
	"os/exec"

	"github.com/natewong1313/go-react-ssr/gossr-cli/cmd"
	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the cli to the latest version",
	Long:  "Update the cli to the latest version",
	Run:   update,
}

func init() {
	if CheckNeedsUpdate() {
		cmd.RootCmd.AddCommand(updateCmd)
	}
}

func update(cmd *cobra.Command, args []string) {
	exec.Command("go", "install", "github.com/natewong1313/go-react-ssr/gossr-cli@latest").Run()
	updateVersionFile()
	logger.L.Info().Msg("Updated to latest version!")
}
