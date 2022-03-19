package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type mySqlConnector baseConnector

func (c *mySqlConnector) GetDbType() DbType {
	return c.dbType
}

func (c *mySqlConnector) Connect() error {
	db, err := sql.Open(c.dbType.String(), c.connectionString)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	c.db = db
	return nil
}

func (c mySqlConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c mySqlConnector) GetSchemas() ([]string, error) {
	rows, err := c.db.Query("select schema_name from information_schema.schemata")
	if err != nil {
		return nil, err
	}

	var schemas []string
	for rows.Next() {
		var schema string
		if err = rows.Scan(&schema); err != nil {
			return nil, err
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (c mySqlConnector) GetTables(schemaName string) ([]string, error) {
	rows, err := c.db.Query(`
		select table_name
		from information_schema.tables
		where table_type = 'BASE TABLE'
		  and table_schema = ?
		`, schemaName)
	if err != nil {
		return nil, err
	}

	var tables []string
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, SanitizeTableName(table))
	}

	return tables, nil
}

func (c mySqlConnector) GetColumns(tableName string) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
		select column_name, data_type
		from information_schema.columns
		where table_name = ?
		order by ordinal_position
		`, tableName)
	if err != nil {
		return nil, err
	}

	var columns []ColumnResult
	for rows.Next() {
		var column ColumnResult
		if err = rows.Scan(&column.Name, &column.DataType); err != nil {
			return nil, err
		}

		column.Name = SanitizeColumnName(column.Name)
		column.DataType = SanitizeColumnType(column.DataType)

		columns = append(columns, column)
	}

	return columns, nil
}

func (c mySqlConnector) GetConstraints(tableName string) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
		select c.TABLE_NAME,
			   c.REFERENCED_TABLE_NAME,
			   c.CONSTRAINT_NAME,
			   (
				   select kc2.CONSTRAINT_NAME is not null "isPrimary"
				   from information_schema.KEY_COLUMN_USAGE kc
							left join information_schema.KEY_COLUMN_USAGE kc2
									  ON kc.COLUMN_NAME = kc2.COLUMN_NAME AND kc2.CONSTRAINT_NAME = 'PRIMARY' AND
										 kc2.TABLE_NAME = c.TABLE_NAME
				   where kc.CONSTRAINT_NAME = c.CONSTRAINT_NAME
			   ) "isPrimary",
			   (
				   select COUNT(*) > 1
				   from information_schema.KEY_COLUMN_USAGE kc
				   where kc.TABLE_NAME = c.TABLE_NAME
					 and kc.CONSTRAINT_NAME = 'PRIMARY'
			   ) "hasMultiplePk"
		from information_schema.REFERENTIAL_CONSTRAINTS c
		where c.TABLE_NAME = ? or c.REFERENCED_TABLE_NAME = ?
		`, tableName, tableName)
	if err != nil {
		return nil, err
	}

	var constraints []ConstraintResult
	for rows.Next() {
		var constraint ConstraintResult
		err = rows.Scan(&constraint.FkTable,
			&constraint.PKTable,
			&constraint.ConstraintName,
			&constraint.IsPrimary,
			&constraint.HasMultiplePK,
		)

		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}
