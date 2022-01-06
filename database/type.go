package database

type DbType string

const (
	Postgres DbType = "pgx"
	MySql    DbType = "mysql"
)

func (c DbType) String() string {
	return string(c)
}
