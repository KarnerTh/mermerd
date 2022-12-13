package database

import (
	"database/sql"
	"fmt"
	"strings"

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

func (c *postgresConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c *postgresConnector) GetSchemas() ([]string, error) {
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

func (c *postgresConnector) GetTables(schemaNames []string) ([]TableNameResult, error) {
	// possible sql injection
	schemaSearch := strings.Join(schemaNames, ",")
	rows, err := c.db.Query(`
		select table_schema, table_name
		from information_schema.tables
		where table_type = 'BASE TABLE'
		  and table_schema in ($1)
		`, schemaSearch)
	if err != nil {
		return nil, err
	}

	var tables []TableNameResult
	for rows.Next() {
		var table TableNameResult
		if err = rows.Scan(&table.Schema, &table.Name); err != nil {
			return nil, err
		}

		table.Name = SanitizeValue(table.Name)

		tables = append(tables, table)
	}

	return tables, nil
}

func (c *postgresConnector) GetColumns(tableName TableNameResult) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
        select c.column_name,
               (case
                    when c.data_type = 'USER-DEFINED'
                        then c.udt_name
                    else c.data_type
                   end)                                                        as data_type,
               (select count(*) > 0
                from information_schema.key_column_usage cu
                         left join information_schema.table_constraints tc on tc.constraint_name = cu.constraint_name
                where cu.column_name = c.column_name
                  and cu.table_name = c.table_name
                  and tc.constraint_type = 'PRIMARY KEY')                      as is_primary,
               (select count(*) > 0
                from information_schema.key_column_usage cu
                         left join information_schema.table_constraints tc on tc.constraint_name = cu.constraint_name
                where cu.column_name = c.column_name
                  and cu.table_name = c.table_name
                  and tc.constraint_type = 'FOREIGN KEY')                      as is_foreign,
               coalesce(string_agg(enumlabel, ',' order by enumsortorder), '') as enum_values
        from information_schema.columns c
                 left join pg_type typ on c.udt_name = typ.typname
                 left join pg_enum enu on typ.oid = enu.enumtypid
        where c.table_name = $1
        group by c.column_name,
                 c.table_name,
                 c.data_type,
                 c.udt_name,
                 c.ordinal_position
        order by c.ordinal_position;
		`, tableName.Name)
	if err != nil {
		return nil, err
	}

	var columns []ColumnResult
	for rows.Next() {
		var column ColumnResult
		if err = rows.Scan(&column.Name, &column.DataType, &column.IsPrimary, &column.IsForeign, &column.EnumValues); err != nil {
			return nil, err
		}

		column.Name = SanitizeValue(column.Name)
		column.DataType = SanitizeValue(column.DataType)

		columns = append(columns, column)
	}

	return columns, nil
}

func (c *postgresConnector) GetConstraints(tableName TableNameResult) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
	select fk.table_name,
		   pk.table_name,
		   c.constraint_name,
		   kcu.column_name,
		   coalesce(
				   (select tc.constraint_type is not null "isPrimary"
					from information_schema.key_column_usage kc
							 inner join information_schema.key_column_usage kc2
										ON kc2.column_name = kc.column_name and kc2.table_name = kc.table_name
							 inner join information_schema.table_constraints tc
										on kc2.constraint_name = tc.constraint_name and
										   tc.constraint_type = 'PRIMARY KEY'
					where kc.constraint_name = c.constraint_name
					  and kc.column_name = kcu.column_name)
			   , false) "isPrimary",
		   (select COUNT(*) > 1 "hasMultiplePk"
			from information_schema.table_constraints tc
					 -- one constraint can have multiple columns
					 inner join information_schema.key_column_usage kc
								on kc.constraint_name = tc.constraint_name
			where tc.table_name = fk.table_name
			  and tc.constraint_type = 'PRIMARY KEY')
	from information_schema.referential_constraints c
			 inner join information_schema.table_constraints fk on c.constraint_name = fk.constraint_name
			 inner join information_schema.table_constraints pk on c.unique_constraint_name = pk.constraint_name
			 inner join information_schema.key_column_usage kcu on c.constraint_name = kcu.constraint_name
	where fk.table_name = $1 or pk.table_name = $1;
		`, tableName.Name)
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
