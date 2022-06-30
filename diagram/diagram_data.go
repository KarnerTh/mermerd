package diagram

type ErdRelationType string

const (
	relationOneToOne  ErdRelationType = "|o--||"
	relationManyToOne ErdRelationType = "}o--||"
)

type ErdDiagramData struct {
	EncloseWithMermaidBackticks bool
	Tables                      []ErdTableData
	Constraints                 []ErdConstraintData
}

type ErdTableData struct {
	Name    string
	Columns []ErdColumnData
}

type ErdColumnData struct {
	Name     string
	DataType string
}

type ErdConstraintData struct {
	PkTableName string
	FkTableName string
	Relation    ErdRelationType
	ColumnName  string
}
