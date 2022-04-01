package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const tfpvcVersion = "LOCAL"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns version data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(tfpvcVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
