package config

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestYamlConfig(t *testing.T) {
	// Arrange
	config := NewConfig()
	var configYaml = []byte(`
# Connection properties
connectionString: "connectionStringExample"
schema: "public"

# Define what tables should be used
useAllTables: false
selectedTables:
  - city
  - customer

# Additional flags
showAllConstraints: true
outputFileName: "my-db.mmd"
encloseWithMermaidBackticks: false
debug: true
omitConstraintLabels: true
omitAttributeKeys: true
showDescriptions:
  - enumValues
  - columnComments
  - notNull
relationshipLabels:
  - "schema.table1 schema.table2 : is_a"
  - "table-name  another-table-name  : has_many"
  - "incorrect format"
useAllSchemas: true
showSchemaPrefix: true
schemaPrefixSeparator: "_"

# These connection strings are available as suggestions in the cli (use tab to access)
connectionStringSuggestions:
  - suggestion1
  - suggestion2
`)

	// Act
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(configYaml))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "connectionStringExample", config.ConnectionString())
	assert.ElementsMatch(t, []string{"public"}, config.Schemas())
	assert.Equal(t, false, config.UseAllTables())
	assert.ElementsMatch(t, []string{"city", "customer"}, config.SelectedTables())
	assert.Equal(t, true, config.ShowAllConstraints())
	assert.Equal(t, "my-db.mmd", config.OutputFileName())
	assert.Equal(t, false, config.EncloseWithMermaidBackticks())
	assert.ElementsMatch(t, []string{"suggestion1", "suggestion2"}, config.ConnectionStringSuggestions())
	assert.True(t, config.Debug())
	assert.True(t, config.OmitConstraintLabels())
	assert.True(t, config.OmitAttributeKeys())
	assert.ElementsMatch(t, []string{"notNull", "enumValues", "columnComments"}, config.ShowDescriptions())
	assert.True(t, config.UseAllSchemas())
	assert.True(t, config.ShowSchemaPrefix())
	assert.Equal(t, "_", config.SchemaPrefixSeparator())
	assert.ElementsMatch(t,
		config.RelationshipLabels(),
		[]RelationshipLabel{
			{
				PkName: "schema.table1",
				FkName: "schema.table2",
				Label:  "is_a",
			},
			{
				PkName: "table-name",
				FkName: "another-table-name",
				Label:  "has_many",
			},
		},
	)
}
