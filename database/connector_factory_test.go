package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {
	connectorFactory := NewConnectorFactory()
	testCases := []struct {
		connectionString string
		expectedDbType   DbType
	}{
		{connectionString: "postgresql://user:password@localhost:5432/yourDb", expectedDbType: Postgres},
		{connectionString: "postgres://user:password@localhost:5432/yourDb", expectedDbType: Postgres},
		{connectionString: "mysql://root:password@tcp(127.0.0.1:3306)/yourDb", expectedDbType: MySql},
		{connectionString: "sqlserver://sa:securePassword1!@localhost:1433?database=mermerd_test", expectedDbType: MsSql},
		{connectionString: "sqlite3://mermerd_test.db", expectedDbType: Sqlite3},
	}

	for index, testCase := range testCases {
		t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
			// Arrange
			connectionString := testCase.connectionString

			// Act
			connector, err := connectorFactory.NewConnector(connectionString)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, connector.GetDbType(), testCase.expectedDbType)
		})
	}
}

func TestUnsupportedConnector(t *testing.T) {
	// Arrange
	connectorFactory := NewConnectorFactory()
	connectionString := "notSupported://user:password@localhost:5432/yourDb"

	// Act
	connector, err := connectorFactory.NewConnector(connectionString)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, connector)
}
