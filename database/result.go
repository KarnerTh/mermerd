package database

type Result struct {
	Tables []TableResult
}

type TableResult struct {
	TableName   string
	Columns     []ColumnResult
	Constraints []ConstraintResult
}

type ColumnResult struct {
	Name     string
	DataType string
}

type ConstraintResult struct {
	FkTable        string
	PKTable        string
	ConstraintName string
	IsPrimary      bool
}
