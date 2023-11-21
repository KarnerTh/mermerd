package database

import (
	"errors"
	"strings"
)

type connectorFactory struct{}

type ConnectorFactory interface {
	NewConnector(connectionString string) (Connector, error)
}

func NewConnectorFactory() ConnectorFactory {
	return connectorFactory{}
}

func (connectorFactory) NewConnector(connectionString string) (Connector, error) {
	switch {
	case strings.HasPrefix(connectionString, "postgresql") || strings.HasPrefix(connectionString, "postgres"):
		return &postgresConnector{
			dbType:           Postgres,
			connectionString: connectionString,
		}, nil
	case strings.HasPrefix(connectionString, "mysql"):
		return &mySqlConnector{
			dbType:           MySql,
			connectionString: strings.ReplaceAll(connectionString, "mysql://", ""),
		}, nil
	case strings.HasPrefix(connectionString, "sqlserver"):
		return &mssqlConnector{
			dbType:           MsSql,
			connectionString: connectionString,
		}, nil
	case strings.HasPrefix(connectionString, "sqlite3"):
		return &sqliteConnector{
			dbType:           Sqlite3,
			connectionString: connectionString,
		}, nil
	default:
		return nil, errors.New("could not create connector for db")
	}
}
