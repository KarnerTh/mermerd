package diagram

type ErdRelationType string

const (
	relationOneToOne  ErdRelationType = "|o--||"
	relationManyToOne ErdRelationType = "}o--||"
)

type ErdAttributeKey string

const (
	primaryKey ErdAttributeKey = "PK"
	foreignKey ErdAttributeKey = "FK"
	none       ErdAttributeKey = ""
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
	Name         string
	DataType     string
	EnumValues   string
	AttributeKey ErdAttributeKey
}

type ErdConstraintData struct {
	PkTableName     string
	FkTableName     string
	Relation        ErdRelationType
	ConstraintLabel string
}
