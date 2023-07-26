package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type mssqlConnector baseConnector

func (c *mssqlConnector) GetDbType() DbType {
	return c.dbType
}

func (c *mssqlConnector) Connect() error {
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

func (c *mssqlConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c *mssqlConnector) GetSchemas() ([]string, error) {
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

func (c *mssqlConnector) GetTables(schemaNames []string) ([]TableDetail, error) {
	args := make([]any, len(schemaNames))
	searchPlaceholder := make([]string, len(schemaNames))
	for i, schemaName := range schemaNames {
		args[i] = schemaName
		searchPlaceholder[i] = fmt.Sprintf("@p%d", i+1)
	}
	rows, err := c.db.Query(`
		select table_schema, table_name
		from information_schema.tables
		where table_type = 'BASE TABLE'
		  and table_schema in(`+strings.Join(searchPlaceholder, ",")+`) 
		`, args...)
	if err != nil {
		return nil, err
	}

	var tables []TableDetail
	for rows.Next() {
		var table TableDetail
		if err = rows.Scan(&table.Schema, &table.Name); err != nil {
			return nil, err
		}

		table.Name = SanitizeValue(table.Name)

		tables = append(tables, table)
	}

	return tables, nil
}

func (c *mssqlConnector) GetColumns(tableName TableDetail) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
		select c.column_name,
			   c.data_type,
			   (select IIF(count(*) > 0, 1, 0)
				from information_schema.key_column_usage cu
						 left join information_schema.table_constraints tc on tc.constraint_name = cu.constraint_name
				where cu.column_name = c.column_name
				  and cu.table_name = c.table_name
				  and tc.constraint_type = 'PRIMARY KEY') as is_primary,
			   (select IIF(count(*) > 0, 1, 0)
				from information_schema.key_column_usage cu
						 left join information_schema.table_constraints tc on tc.constraint_name = cu.constraint_name
				where cu.column_name = c.column_name
				  and cu.table_name = c.table_name
				  and tc.constraint_type = 'FOREIGN KEY') as is_foreign,
		    	case when c.is_nullable = 'YES' then 1 else 0 end as is_nullable,
			   (select ISNULL(ep.value, '') from sys.tables t
			      inner join sys.columns col on col.object_id = t.object_id and col.name = c.column_name
				  left join sys.extended_properties ep on ep.major_id = t.object_id and ep.minor_id = col.column_id
				  where t.name = c.table_name and SCHEMA_NAME(t.schema_id) = c.TABLE_SCHEMA) as comment
		from information_schema.columns c
		where c.table_name = @p1 and c.TABLE_SCHEMA = @p2
		order by c.ordinal_position;
		`, tableName.Name, tableName.Schema)
	if err != nil {
		return nil, err
	}

	var columns []ColumnResult
	for rows.Next() {
		var column ColumnResult
		if err = rows.Scan(&column.Name, &column.DataType, &column.IsPrimary, &column.IsForeign, &column.IsNullable, &column.Comment); err != nil {
			return nil, err
		}

		column.Name = SanitizeValue(column.Name)
		column.DataType = SanitizeValue(column.DataType)

		columns = append(columns, column)
	}

	return columns, nil
}

func (c *mssqlConnector) GetConstraints(tableName TableDetail) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
select fk.table_name,
       fk.table_schema,
       pk.table_name,
       pk.table_schema,
       c.constraint_name,
       kcu.column_name,
       coalesce(
               (select IIF(tc.constraint_type is not null, 'true', 'false')
                from information_schema.key_column_usage kc
                         inner join information_schema.key_column_usage kc2
                                    ON kc2.column_name = kc.column_name and kc2.table_name = kc.table_name
                         inner join information_schema.table_constraints tc
                                    on kc2.constraint_name = tc.constraint_name and
                                       tc.constraint_type = 'PRIMARY KEY'
                where kc.constraint_name = c.constraint_name
                  and kc.column_name = kcu.column_name), 'false') "isPrimary",
       (select IIF(COUNT(*) > 1, 'true', 'false')
        from information_schema.table_constraints tc
                 -- one constraint can have multiple columns
                 inner join information_schema.key_column_usage kc
                            on kc.constraint_name = tc.constraint_name
        where tc.table_name = fk.table_name
          and tc.constraint_type = 'PRIMARY KEY') "hasMultiplePk"
from information_schema.referential_constraints c
         inner join information_schema.table_constraints fk on c.constraint_name = fk.constraint_name
         inner join information_schema.table_constraints pk on c.unique_constraint_name = pk.constraint_name
         inner join information_schema.key_column_usage kcu on c.constraint_name = kcu.constraint_name
where c.CONSTRAINT_SCHEMA = @p1 and (fk.table_name = @p2 or pk.table_name = @p2);
		`, tableName.Schema, tableName.Name)
	if err != nil {
		return nil, err
	}

	var constraints []ConstraintResult
	for rows.Next() {
		var constraint ConstraintResult
		err = rows.Scan(
			&constraint.FkTable,
			&constraint.FkSchema,
			&constraint.PkTable,
			&constraint.PkSchema,
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
