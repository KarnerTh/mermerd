package database

type Result struct {
	Tables []TableResult
}

type TableResult struct {
	Table       TableDetail
	Columns     []ColumnResult
	Constraints ConstraintResultList
}

type TableDetail struct {
	Schema string
	Name   string
}

type ColumnResult struct {
	Name       string
	DataType   string
	IsPrimary  bool
	IsForeign  bool
	IsUnique   bool
	IsNullable bool
	EnumValues string
	Comment    string
}

type ConstraintResultList []ConstraintResult
type ConstraintResult struct {
	FkTable        string
	FkSchema       string
	PkTable        string
	PkSchema       string
	ConstraintName string
	ColumnName     string
	IsPrimary      bool
	HasMultiplePK  bool
}

// AppendIfNotExists ensures that only unique items are appended to the list of constraints
func (source ConstraintResultList) AppendIfNotExists(items ...ConstraintResult) ConstraintResultList {
	result := source
	for _, item := range items {
		if !sliceContainsConstraint(result, item) {
			result = append(result, item)
		}
	}

	return result
}

func sliceContainsConstraint(slice []ConstraintResult, item ConstraintResult) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}

	return false
}
