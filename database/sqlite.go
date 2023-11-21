package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type sqliteConnector baseConnector

func (c *sqliteConnector) GetDbType() DbType {
	return c.dbType
}

func getFilenameFromConnectionString(connectionString string) string {
	return strings.Replace(connectionString, "sqlite3://", "", 1)
}

func (c *sqliteConnector) Connect() error {
	db, err := sql.Open(c.dbType.String(), getFilenameFromConnectionString(c.connectionString))
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	c.db = db
	return nil
}

func (c *sqliteConnector) Close() {
	err := c.db.Close()
	if err != nil {
		fmt.Println("could not close database connection", err)
	}
}

func (c *sqliteConnector) GetSchemas() ([]string, error) {
	fileName := getFilenameFromConnectionString(c.connectionString)
	schema := strings.Replace(fileName, ".db", "", 1)
	return []string{schema}, nil
}

func (c *sqliteConnector) GetTables(schemaNames []string) ([]TableDetail, error) {
	rows, err := c.db.Query(`
    select distinct tbl_name 
    from sqlite_schema
		`)
	if err != nil {
		return nil, err
	}

	var tables []TableDetail
	for rows.Next() {
		table := TableDetail{Schema: schemaNames[0]}
		if err = rows.Scan(&table.Name); err != nil {
			return nil, err
		}

		table.Name = SanitizeValue(table.Name)
		tables = append(tables, table)
	}

	return tables, nil
}

func (c *sqliteConnector) GetColumns(tableName TableDetail) ([]ColumnResult, error) {
	rows, err := c.db.Query(`
 select 
    t.name, 
    type,
    pk > 0, -- first pk has 1, second has 2 ...
    (case when fk.id is not null then 1 else 0 end) "isForeign",
    (case when t."notnull" = true then false else true end) "nullable",
    coalesce(
      (select (case when i."unique" = true and i.origin = "u" then true else false end) "isUnique"
        from pragma_index_list(:tableName) i
      left join pragma_index_info(i.name) ii
      where ii.name = t.name),
        0
   ) "isUnique"
    from pragma_table_info(:tableName) t 
    left join pragma_foreign_key_list(:tableName) fk on t.name = fk."from";
		`, sql.Named("tableName", tableName.Name))
	if err != nil {
		return nil, err
	}

	var columns []ColumnResult
	for rows.Next() {
		var column ColumnResult
		if err = rows.Scan(&column.Name, &column.DataType, &column.IsPrimary, &column.IsForeign, &column.IsNullable, &column.IsUnique); err != nil {
			return nil, err
		}

		column.Name = SanitizeValue(column.Name)
		column.DataType = SanitizeValue(column.DataType)

		columns = append(columns, column)
	}

	return columns, nil
}

func (c *sqliteConnector) GetConstraints(tableName TableDetail) ([]ConstraintResult, error) {
	rows, err := c.db.Query(`
select 
  pk."table" "pkTableName",
  fk.name "pkTableName",
  pk."from", 
    coalesce((select pk > 0 from pragma_table_info(fk.name) ti where ti.name = pk."from"), 0) "isPrimary", -- first pk has 1, second has 2 ...
    coalesce((select count(*) > 1 from pragma_table_info(fk.name) ti where pk > 0), 0) "hasMultiplePK"
    from sqlite_master fk
      join pragma_foreign_key_list(fk.name) pk on pk."table" != fk.name
  where fk.name = :tableName or pk."table" = :tableName
		`, sql.Named("tableName", tableName.Name))
	if err != nil {
		return nil, err
	}

	var constraints []ConstraintResult
	for rows.Next() {
		constraint := ConstraintResult{
			PkSchema: tableName.Schema,
			FkSchema: tableName.Schema, // sqlite only has one schema
		}
		err = rows.Scan(
			&constraint.PkTable,
			&constraint.FkTable,
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
