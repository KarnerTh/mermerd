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
	GetTables(schemaNames []string) ([]TableNameResult, error)
	GetColumns(tableName TableNameResult) ([]ColumnResult, error)
	GetConstraints(tableName TableNameResult) ([]ConstraintResult, error)
}
