package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tfDir string
var errorOnUpdate bool

var rootCmd = &cobra.Command{
	Use:   "tfpvc",
	Short: "Terraform Provider Version Check",
	Long: `A utility to check whether Terraform providers configured in a project are up-to-date.

* Reads .terraform.lock.hcl file to get details of providers used in a project
* Looks up registry.terraform.io to get the latest version for each
* If updates are available this is notified on STDOUT
  * Can optionally return code on exit if update is available`,
	Run: func(cmd *cobra.Command, args []string) {
		check()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&tfDir, "tfDir", ".", "Directory with TF Files")
	rootCmd.PersistentFlags().BoolVar(&errorOnUpdate, "errorOnUpdate", false, "Exit with error code if updates are available")
}

func initConfig() {
	viper.AutomaticEnv()
}
