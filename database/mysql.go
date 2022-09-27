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

func (c *mySqlConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c *mySqlConnector) GetSchemas() ([]string, error) {
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

func (c *mySqlConnector) GetTables(schemaName string) ([]string, error) {
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

		tables = append(tables, SanitizeValue(table))
	}

	return tables, nil
}

func (c *mySqlConnector) GetColumns(tableName string) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
		select c.column_name,
			   c.data_type,
			   (select count(*) > 0
				from information_schema.KEY_COLUMN_USAGE
				where table_name = c.table_name
				  and column_name = c.column_name
				  and constraint_name = 'PRIMARY')        as is_primary,
			   (select count(*) > 0
				from information_schema.key_column_usage cu
						 left join information_schema.table_constraints tc on tc.constraint_name = cu.constraint_name
				where cu.column_name = c.column_name
				  and cu.table_name = c.table_name
				  and tc.constraint_type = 'FOREIGN KEY') as is_foreign
		from information_schema.columns c
		where c.table_name = ?
		order by c.ordinal_position;
		`, tableName)
	if err != nil {
		return nil, err
	}

	var columns []ColumnResult
	for rows.Next() {
		var column ColumnResult
		if err = rows.Scan(&column.Name, &column.DataType, &column.IsPrimary, &column.IsForeign); err != nil {
			return nil, err
		}

		column.Name = SanitizeValue(column.Name)
		column.DataType = SanitizeValue(column.DataType)

		columns = append(columns, column)
	}

	return columns, nil
}

func (c *mySqlConnector) GetConstraints(tableName string) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
		select c.TABLE_NAME,
			   c.REFERENCED_TABLE_NAME,
			   c.CONSTRAINT_NAME,
       		   kcu.COLUMN_NAME,
			   (
				   select kc2.CONSTRAINT_NAME is not null "isPrimary"
				   from information_schema.KEY_COLUMN_USAGE kc
							left join information_schema.KEY_COLUMN_USAGE kc2
									  ON kc.COLUMN_NAME = kc2.COLUMN_NAME AND kc2.CONSTRAINT_NAME = 'PRIMARY' AND
										 kc2.TABLE_NAME = kc.TABLE_NAME
           		    where kc.CONSTRAINT_NAME = c.CONSTRAINT_NAME and kc.COLUMN_NAME = kcu.COLUMN_NAME
			   ) "isPrimary",
			   (
				   select COUNT(*) > 1
				   from information_schema.KEY_COLUMN_USAGE kc
				   where kc.TABLE_NAME = c.TABLE_NAME
					 and kc.CONSTRAINT_NAME = 'PRIMARY'
			   ) "hasMultiplePk"
		from information_schema.REFERENTIAL_CONSTRAINTS c
    		inner join information_schema.KEY_COLUMN_USAGE kcu on c.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
		where c.TABLE_NAME = ? or c.REFERENCED_TABLE_NAME = ?
		`, tableName, tableName)
	if err != nil {
		return nil, err
	}

	var constraints []ConstraintResult
	for rows.Next() {
		var constraint ConstraintResult
		err = rows.Scan(
			&constraint.FkTable,
			&constraint.PkTable,
			&constraint.ConstraintName,
			&constraint.ColumnName,
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
