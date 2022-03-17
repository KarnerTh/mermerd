package config

import "github.com/spf13/viper"

const (
	ShowAllConstraintsKey          = "showAllConstraints"
	UseAllTablesKey                = "useAllTables"
	SelectedTablesKey              = "selectedTables"
	SchemaKey                      = "schema"
	ConnectionStringKey            = "connectionString"
	ConnectionStringSuggestionsKey = "connectionStringSuggestions"
	OutputFileNameKey              = "outputFileName"
	EncloseWithMermaidBackticksKey = "encloseWithMermaidBackticks"
)

func ShowAllConstraints() bool {
	return viper.GetBool(ShowAllConstraintsKey)
}

func UseAllTables() bool {
	return viper.GetBool(UseAllTablesKey)
}

func Schema() string {
	return viper.GetString(SchemaKey)
}

func ConnectionString() string {
	return viper.GetString(ConnectionStringKey)
}

func OutputFileName() string {
	return viper.GetString(OutputFileNameKey)
}

func ConnectionStringSuggestions() []string {
	return viper.GetStringSlice(ConnectionStringSuggestionsKey)
}

func SelectedTables() []string {
	return viper.GetStringSlice(SelectedTablesKey)
}

func EncloseWithMermaidBackticks() bool {
	return viper.GetBool(EncloseWithMermaidBackticksKey)
}
