package config

import "github.com/spf13/viper"

const (
	ConnectionStringKey            = "connectionString"
	ConnectionStringSuggestionsKey = "connectionStringSuggestions"
	DebugKey                       = "debug"
	EncloseWithMermaidBackticksKey = "encloseWithMermaidBackticks"
	IgnoreTables                   = "ignoreTables"
	OmitAttributeKeysKey           = "omitAttributeKeys"
	OmitConstraintLabelsKey        = "omitConstraintLabels"
	OutputFileNameKey              = "outputFileName"
	OutputMode                     = "outputMode"
	RelationshipLabelsKey          = "relationshipLabels"
	SchemaKey                      = "schema"
	SchemaPrefixSeparator          = "schemaPrefixSeparator"
	SelectedTablesKey              = "selectedTables"
	ShowAllConstraintsKey          = "showAllConstraints"
	ShowDescriptionsKey            = "showDescriptions"
	ShowSchemaPrefix               = "showSchemaPrefix"
	UseAllSchemasKey               = "useAllSchemas"
	UseAllTablesKey                = "useAllTables"
)

type config struct{}

type MermerdConfig interface {
	ConnectionString() string
	ConnectionStringSuggestions() []string
	Debug() bool
	EncloseWithMermaidBackticks() bool
	IgnoreTables() []string
	OmitAttributeKeys() bool
	OmitConstraintLabels() bool
	OutputFileName() string
	OutputMode() OutputModeType
	RelationshipLabels() []RelationshipLabel
	SchemaPrefixSeparator() string
	Schemas() []string
	SelectedTables() []string
	ShowAllConstraints() bool
	ShowDescriptions() []string
	ShowSchemaPrefix() bool
	UseAllSchemas() bool
	UseAllTables() bool
}

func NewConfig() MermerdConfig {
	return config{}
}

func (c config) ShowAllConstraints() bool {
	return viper.GetBool(ShowAllConstraintsKey)
}

func (c config) UseAllTables() bool {
	return viper.GetBool(UseAllTablesKey)
}

func (c config) IgnoreTables() []string {
	return viper.GetStringSlice(IgnoreTables)
}

func (c config) Schemas() []string {
	return viper.GetStringSlice(SchemaKey)
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

func (c config) SelectedTables() []string {
	return viper.GetStringSlice(SelectedTablesKey)
}

func (c config) RelationshipLabels() []RelationshipLabel {
	labels := viper.GetStringSlice(RelationshipLabelsKey)
	return parseLabels(labels)
}

func (c config) EncloseWithMermaidBackticks() bool {
	return viper.GetBool(EncloseWithMermaidBackticksKey)
}

func (c config) Debug() bool {
	return viper.GetBool(DebugKey)
}

func (c config) OmitConstraintLabels() bool {
	return viper.GetBool(OmitConstraintLabelsKey)
}

func (c config) OmitAttributeKeys() bool {
	return viper.GetBool(OmitAttributeKeysKey)
}

func (c config) ShowDescriptions() []string {
	return viper.GetStringSlice(ShowDescriptionsKey)
}

func (c config) UseAllSchemas() bool {
	return viper.GetBool(UseAllSchemasKey)
}

func (c config) ShowSchemaPrefix() bool {
	return viper.GetBool(ShowSchemaPrefix)
}

func (c config) SchemaPrefixSeparator() string {
	return viper.GetString(SchemaPrefixSeparator)
}

func (c config) OutputMode() OutputModeType {
	return OutputModeType(viper.GetString(OutputMode))
}
