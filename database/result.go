package database

type Result struct {
	Tables []TableResult
}

type TableResult struct {
	TableName   string
	Columns     []ColumnResult
	Constraints ConstraintResultList
}

type ColumnResult struct {
	Name     string
	DataType string
}

type ConstraintResultList []ConstraintResult
type ConstraintResult struct {
	FkTable        string
	PkTable        string
	ConstraintName string
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
