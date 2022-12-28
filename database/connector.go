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
	GetTables(schemaNames []string) ([]TableDetail, error)
	GetColumns(tableName TableDetail) ([]ColumnResult, error)
	GetConstraints(tableName TableDetail) ([]ConstraintResult, error)
}
