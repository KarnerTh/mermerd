package database

import "database/sql"

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
