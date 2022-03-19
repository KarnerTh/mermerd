package database

import (
	"database/sql"
	"errors"
	"strings"
)

type baseConnector struct {
	dbType           DbType
	connectionString string
	db               *sql.DB
}

type Connector interface {
	Connect() error
	Close()
	GetDbType() DbType
	GetSchemas() ([]string, error)
	GetTables(schemaName string) ([]string, error)
	GetColumns(tableName string) ([]ColumnResult, error)
	GetConstraints(tableName string) ([]ConstraintResult, error)
}

func NewConnector(connectionString string) (Connector, error) {
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
	default:
		return nil, errors.New("could not create connector for db")
	}
}
