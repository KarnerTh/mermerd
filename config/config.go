package config

import (
	"github.com/spf13/viper"
)

const (
	ShowAllConstraintsKey          = "showAllConstraints"
	ShowOnlySelectedColumnsKey     = "showOnlySelectedColumns"
	UseAllTablesKey                = "useAllTables"
	SelectedTablesKey              = "selectedTables"
	SchemaKey                      = "schema"
	ConnectionStringKey            = "connectionString"
	ConnectionStringSuggestionsKey = "connectionStringSuggestions"
	OutputFileNameKey              = "outputFileName"
	EncloseWithMermaidBackticksKey = "encloseWithMermaidBackticks"
)

type config struct{}

type MermerdConfig interface {
	ShowAllConstraints() bool
	ShowOnlySelectedColumns() bool
	UseAllTables() bool
	Schema() string
	ConnectionString() string
	OutputFileName() string
	ConnectionStringSuggestions() []string
	SelectedTables() map[string]Table
	EncloseWithMermaidBackticks() bool
}

func NewConfig() MermerdConfig {
	return config{}
}

func (c config) ShowAllConstraints() bool {
	return viper.GetBool(ShowAllConstraintsKey)
}

func (c config) ShowOnlySelectedColumns() bool {
	return viper.GetBool(ShowOnlySelectedColumnsKey)
}

func (c config) UseAllTables() bool {
	return viper.GetBool(UseAllTablesKey)
}

func (c config) Schema() string {
	return viper.GetString(SchemaKey)
}

func (c config) ConnectionString() string {
	return viper.GetString(ConnectionStringKey)
}

func (c config) OutputFileName() string {
	return viper.GetString(OutputFileNameKey)
}

func (c config) ConnectionStringSuggestions() []string {
	return viper.GetStringSlice(ConnectionStringSuggestionsKey)
}

type Table struct {
	Name    string
	Columns []string
}

func (c config) SelectedTables() map[string]Table {
	var selectedTables map[string]Table
	selectedTables = map[string]Table{}

	var tables []map[string]interface{}
	viper.UnmarshalKey(SelectedTablesKey, &tables)

	for _, table := range tables {
		var selectedColumns []string

		if columns, ok := table["columns"]; ok {
			for _, column := range columns.([]interface{}) {
				selectedColumns = append(selectedColumns, column.(string))
			}
		}

		selectedTable := Table{Name: table["name"].(string), Columns: selectedColumns}
		selectedTables[selectedTable.Name] = selectedTable
	}

	return selectedTables
}

func (c config) EncloseWithMermaidBackticks() bool {
	return viper.GetBool(EncloseWithMermaidBackticksKey)
}
