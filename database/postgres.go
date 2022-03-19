package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type postgresConnector baseConnector

func (c *postgresConnector) GetDbType() DbType {
	return c.dbType
}

func (c *postgresConnector) Connect() error {
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

func (c postgresConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c postgresConnector) GetSchemas() ([]string, error) {
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

func (c postgresConnector) GetTables(schemaName string) ([]string, error) {
	rows, err := c.db.Query(`
		select table_name
		from information_schema.tables
		where table_type = 'BASE TABLE'
		  and table_schema = $1
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

func (c postgresConnector) GetColumns(tableName string) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
		select column_name, data_type
		from information_schema.columns
		where table_name = $1
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

func (c postgresConnector) GetConstraints(tableName string) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
		select distinct fk.table_name,
						pk.table_name,
						c.constraint_name,
						coalesce(
								(select tc.constraint_type is not null "isPrimary"
								 from information_schema.key_column_usage kc
										  inner join information_schema.key_column_usage kc2
													 ON kc2.column_name = kc.column_name and kc2.table_name = kc.table_name
										  inner join information_schema.table_constraints tc
													 on kc2.constraint_name = tc.constraint_name and
														tc.constraint_type = 'PRIMARY KEY'
								 where kc.constraint_name = c.constraint_name)
							, false) "isPrimary",
						(
							select COUNT(*) > 1 "hasMultiplePk"
							from information_schema.table_constraints tc
									 -- one constraint can have multiple columns
									 inner join information_schema.key_column_usage kc
												on kc.constraint_name = tc.constraint_name
							where tc.table_name = fk.table_name
							  and tc.constraint_type = 'PRIMARY KEY'
						)
		from information_schema.referential_constraints c
				 inner join information_schema.table_constraints fk on c.constraint_name = fk.constraint_name
				 inner join information_schema.table_constraints pk on c.unique_constraint_name = pk.constraint_name
		where fk.table_name = $1 or pk.table_name = $1
		`, tableName)
	if err != nil {
		return nil, err
	}

	var constraints []ConstraintResult
	for rows.Next() {
		var constraint ConstraintResult
		err = rows.Scan(
			&constraint.FkTable,
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
