package database

type DbType string

const (
	Postgres DbType = "pgx"
	MySql    DbType = "mysql"
	MsSql    DbType = "sqlserver"
	Sqlite3  DbType = "sqlite"
)

func (c DbType) String() string {
	return string(c)
}
