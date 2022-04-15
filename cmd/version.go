package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mermerd",
	Long:  "All software has versions. This is mermerd's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mermerd %s %s\n", viper.Get("version"), viper.Get("commit"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
