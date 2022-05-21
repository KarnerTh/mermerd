package analyzer

import (
	"testing"

	"github.com/KarnerTh/mermerd/database"
	"github.com/KarnerTh/mermerd/mocks"
	"github.com/stretchr/testify/assert"
)

func getAnalyzerWithMocks() (Analyzer, *mocks.MermerdConfig, *mocks.ConnectorFactory, *mocks.Questioner) {
	configMock := mocks.MermerdConfig{}
	connectionFactoryMock := mocks.ConnectorFactory{}
	questionerMock := mocks.Questioner{}
	return NewAnalyzer(&configMock, &connectionFactoryMock, &questionerMock), &configMock, &connectionFactoryMock, &questionerMock
}

func TestAnalyzer_GetConnectionString(t *testing.T) {
	t.Run("Use value from config", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		configMock.On("ConnectionString").Return("configuredConnectionString").Once()

		// Act
		result, err := analyzer.GetConnectionString()

		// Assert
		configMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "configuredConnectionString", result)
	})

	t.Run("Use value from questioner", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, questionerMock := getAnalyzerWithMocks()
		configMock.On("ConnectionString").Return("").Once()
		configMock.On("ConnectionStringSuggestions").Return([]string{"suggestion"})
		questionerMock.On("AskConnectionQuestion", []string{"suggestion"}).Return("validConnectionString", nil)

		// Act
		result, err := analyzer.GetConnectionString()

		// Assert
		configMock.AssertExpectations(t)
		questionerMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "validConnectionString", result)
	})
}

func TestAnalyzer_GetSchema(t *testing.T) {
	t.Run("Use value from config", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("Schema").Return("configuredSchema").Once()

		// Act
		result, err := analyzer.GetSchema(&connectorMock)

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "configuredSchema", result)
	})

	t.Run("No schema available return error", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("Schema").Return("").Once()
		connectorMock.On("GetSchemas").Return([]string{}, nil).Once()

		// Act
		result, err := analyzer.GetSchema(&connectorMock)

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Empty(t, result)
	})

	t.Run("Use the only returned schema", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("Schema").Return("").Once()
		connectorMock.On("GetSchemas").Return([]string{"onlyItem"}, nil).Once()

		// Act
		result, err := analyzer.GetSchema(&connectorMock)

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "onlyItem", result)
	})

	t.Run("Use value from questioner", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, questionerMock := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("Schema").Return("").Once()
		connectorMock.On("GetSchemas").Return([]string{"first", "second"}, nil).Once()
		questionerMock.On("AskSchemaQuestion", []string{"first", "second"}).Return("first", nil).Once()

		// Act
		result, err := analyzer.GetSchema(&connectorMock)

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		questionerMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "first", result)
	})
}

func TestAnalyzer_GetTables(t *testing.T) {
	t.Run("Use value from config", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("SelectedTables").Return([]string{"configuredTable"}).Once()

		// Act
		result, err := analyzer.GetTables(&connectorMock, "validSchema")

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"configuredTable"}, result)
	})

	t.Run("Use all available tables", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, _ := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("SelectedTables").Return([]string{}).Once()
		connectorMock.On("GetTables", "validSchema").Return([]string{"tableA", "tableB"}, nil).Once()
		configMock.On("UseAllTables").Return(true).Once()

		// Act
		result, err := analyzer.GetTables(&connectorMock, "validSchema")

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"tableA", "tableB"}, result)
	})

	t.Run("Use value from questioner", func(t *testing.T) {
		// Arrange
		analyzer, configMock, _, questionerMock := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("SelectedTables").Return([]string{}).Once()
		connectorMock.On("GetTables", "validSchema").Return([]string{"tableA", "tableB"}, nil).Once()
		configMock.On("UseAllTables").Return(false).Once()
		questionerMock.On("AskTableQuestion", []string{"tableA", "tableB"}).Return([]string{"tableA"}, nil).Once()

		// Act
		result, err := analyzer.GetTables(&connectorMock, "validSchema")

		// Assert
		configMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		questionerMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"tableA"}, result)
	})
}

func TestAnalyzer_Analyze(t *testing.T) {
	t.Run("Existing run configuration does not ask for input", func(t *testing.T) {
		// Arrange
		analyzer, configMock, connectionFactoryMock, questionerMock := getAnalyzerWithMocks()
		connectorMock := mocks.Connector{}
		configMock.On("ConnectionString").Return("validConnectionString").Once()
		connectionFactoryMock.On("NewConnector", "validConnectionString").Return(&connectorMock, nil).Once()
		connectorMock.On("Connect").Return(nil).Once()
		connectorMock.On("Close").Return().Once()
		configMock.On("Schema").Return("validSchema").Once()
		configMock.On("SelectedTables").Return([]string{"tableA", "tableB"}).Once()
		connectorMock.On("GetColumns", "tableA").Return([]database.ColumnResult{
			{
				Name:     "fieldA",
				DataType: "int",
			},
			{
				Name:     "fieldB",
				DataType: "string",
			},
		}, nil).Once()
		connectorMock.On("GetColumns", "tableB").Return([]database.ColumnResult{
			{
				Name:     "fieldC",
				DataType: "int",
			},
			{
				Name:     "fieldD",
				DataType: "string",
			},
		}, nil).Once()
		connectorMock.On("GetConstraints", "tableA").Return([]database.ConstraintResult{{
			FkTable:        "tableA",
			PkTable:        "tableB",
			ConstraintName: "testConstraint",
			IsPrimary:      false,
			HasMultiplePK:  false,
		}}, nil).Once()
		connectorMock.On("GetConstraints", "tableB").Return([]database.ConstraintResult{{
			FkTable:        "tableA",
			PkTable:        "tableB",
			ConstraintName: "testConstraint",
			IsPrimary:      false,
			HasMultiplePK:  false,
		}}, nil).Once()

		// Act
		result, err := analyzer.Analyze()

		// Assert
		configMock.AssertExpectations(t)
		connectionFactoryMock.AssertExpectations(t)
		questionerMock.AssertExpectations(t)
		connectorMock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
}
