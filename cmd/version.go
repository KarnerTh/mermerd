package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mermerd",
	Long:  "All software has versions. This is mermerd's",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO from git info?
		fmt.Println("mermerd v0.0.3")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
