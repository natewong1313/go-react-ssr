package cmd

import (
	"os"

	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gossr-cli",
	Short: "This application helps you get a go-react-ssr powered app up and running in no time.",
	Long: `This application helps you get a go-react-ssr powered app up and running in no time.
	Complete documentation is available at https://github.com/natewong1313/go-react-ssr`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logger.Init()
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
