package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/KarnerTh/mermerd/analyzer"
	"github.com/KarnerTh/mermerd/config"
	"github.com/KarnerTh/mermerd/database"
	"github.com/KarnerTh/mermerd/diagram"
	"github.com/KarnerTh/mermerd/presentation"
)

var runConfig string

var rootCmd = &cobra.Command{
	Use:   "mermerd",
	Short: "Create Mermaid ERD diagrams from existing tables",
	Long:  "Create Mermaid ERD diagrams from existing tables",
	Run: func(cmd *cobra.Command, args []string) {
		presentation.ShowIntro()
		config := config.NewConfig()
		connectorFactory := database.NewConnectorFactory()
		questioner := analyzer.NewQuestioner()
		analyzer := analyzer.NewAnalyzer(config, connectorFactory, questioner)
		diagram := diagram.NewDiagram(config)

		if !config.Debug() {
			logrus.SetOutput(io.Discard)
		}

		result, err := analyzer.Analyze()
		if err != nil {
			logrus.Error(err)
			presentation.ShowError()
			os.Exit(1)
		}

		err = diagram.Create(result)
		if err != nil {
			logrus.Error(err)
			presentation.ShowError()
			os.Exit(1)
		}

		presentation.ShowSuccess(config.OutputFileName())
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
	rootCmd.Flags().Bool(config.UseAllSchemasKey, false, "use all available schemas")
	rootCmd.Flags().Bool(config.DebugKey, false, "show debug logs")
	rootCmd.Flags().Bool(config.OmitConstraintLabelsKey, false, "omit the constraint labels")
	rootCmd.Flags().Bool(config.OmitAttributeKeysKey, false, "omit the attribute keys (PK, FK)")
	rootCmd.Flags().Bool(config.ShowSchemaPrefix, false, "show schema prefix in table name")
	rootCmd.Flags().BoolP(config.EncloseWithMermaidBackticksKey, "e", false, "enclose output with mermaid backticks (needed for e.g. in markdown viewer)")
	rootCmd.Flags().StringP(config.ConnectionStringKey, "c", "", "connection string that should be used")
	rootCmd.Flags().StringP(config.SchemaKey, "s", "", "schema that should be used")
	rootCmd.Flags().StringP(config.OutputFileNameKey, "o", "result.mmd", "output file name")
	rootCmd.Flags().String(config.ShowDescriptionsKey, "", "show 'enumValues' or 'columnComments' in description column")
	rootCmd.Flags().String(config.SchemaPrefixSeparator, ".", "the separator that should be used between schema and table name")
	rootCmd.Flags().StringSlice(config.SelectedTablesKey, []string{""}, "tables to include")

	bindFlagToViper(config.ShowAllConstraintsKey)
	bindFlagToViper(config.UseAllTablesKey)
	bindFlagToViper(config.UseAllSchemasKey)
	bindFlagToViper(config.DebugKey)
	bindFlagToViper(config.OmitConstraintLabelsKey)
	bindFlagToViper(config.OmitAttributeKeysKey)
	bindFlagToViper(config.EncloseWithMermaidBackticksKey)
	bindFlagToViper(config.ConnectionStringKey)
	bindFlagToViper(config.SchemaKey)
	bindFlagToViper(config.OutputFileNameKey)
	bindFlagToViper(config.SelectedTablesKey)
	bindFlagToViper(config.ShowDescriptionsKey)
	bindFlagToViper(config.ShowSchemaPrefix)
	bindFlagToViper(config.SchemaPrefixSeparator)
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

	_ = viper.ReadInConfig()

	// expand all environment variables (https://github.com/spf13/viper/issues/119#issuecomment-417638360)
	for _, k := range viper.AllKeys() {
		value := viper.Get(k)
		if _, ok := value.(string); ok {
			viper.Set(k, os.ExpandEnv(viper.GetString(k)))
		}
	}
}
