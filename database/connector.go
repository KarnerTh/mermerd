package database

import (
	"database/sql"
	"errors"
	"strings"
)

type baseConnector struct {
	dbType         DbType
	dataSourceName string
	db             *sql.DB
}

type Connector interface {
	Connect() error
	Close()
	GetSchemas() ([]string, error)
	GetTables(schemaName string) ([]string, error)
	GetColumns(tableName string) ([]ColumnResult, error)
	GetConstraints(tableName string) ([]ConstraintResult, error)
}

func NewConnector(dataSourceName string) (Connector, error) {
	switch {
	case strings.HasPrefix(dataSourceName, "postgresql") || strings.HasPrefix(dataSourceName, "postgres"):
		return &postgresConnector{
			dbType:         Postgres,
			dataSourceName: dataSourceName,
		}, nil
	case strings.HasPrefix(dataSourceName, "mysql"):
		return &mysqlConnector{
			dbType:         MySql,
			dataSourceName: strings.ReplaceAll(dataSourceName, "mysql://", ""),
		}, nil
	default:
		return nil, errors.New("could not create connector for db")
	}
}
