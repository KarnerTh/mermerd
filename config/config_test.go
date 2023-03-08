package config

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
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
showEnumValues: true
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
	assert.True(t, config.ShowEnumValues())
	assert.True(t, config.UseAllSchemas())
	assert.True(t, config.ShowSchemaPrefix())
	assert.Equal(t, "_", config.SchemaPrefixSeparator())
}
