package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"mermerd/analyzer"
	"mermerd/config"
	"mermerd/diagram"
	"mermerd/util"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runConfig string

var rootCmd = &cobra.Command{
	Use:   "mermerd",
	Short: "Create Mermaid ERD diagrams from existing tables",
	Long:  "Create Mermaid ERD diagrams from existing tables",
	Run: func(cmd *cobra.Command, args []string) {
		util.ShowIntro()
		result, err := analyzer.Analyze()
		if err != nil {
			fmt.Println(err.Error())
			util.ShowError()
			os.Exit(1)
		}

		err = diagram.Create(result)
		if err != nil {
			fmt.Println(err.Error())
			util.ShowError()
			os.Exit(1)
		}

		util.ShowSuccess()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&runConfig, "runConfig", "", "run configuration (replaces global configuration)")
	rootCmd.Flags().Bool(config.ShowAllConstraintsKey, false, "show all constraints, even though the table of the resulting constraint was not selected")
	rootCmd.Flags().Bool(config.UseAllTablesKey, false, "use all available tables")
	rootCmd.Flags().StringP(config.ConnectionStringKey, "c", "", "connection string that should be used")
	rootCmd.Flags().StringP(config.SchemaKey, "s", "", "schema that should be used")
	rootCmd.Flags().StringP(config.OutputFileNameKey, "o", "result.mmd", "output file name")

	bindFlagToViper(config.ShowAllConstraintsKey)
	bindFlagToViper(config.UseAllTablesKey)
	bindFlagToViper(config.ConnectionStringKey)
	bindFlagToViper(config.SchemaKey)
	bindFlagToViper(config.OutputFileNameKey)
}

func bindFlagToViper(key string) {
	_ = viper.BindPFlag(key, rootCmd.Flags().Lookup(key))
}

func initConfig() {
	if runConfig != "" {
		color.Blue(fmt.Sprintf("Using run configuration (from %s)", runConfig))
		viper.SetConfigFile(runConfig)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mermerd")
	}

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
