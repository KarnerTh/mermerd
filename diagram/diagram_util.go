package diagram

import (
	"github.com/KarnerTh/mermerd/config"
	"github.com/KarnerTh/mermerd/database"
)

func getRelation(constraint database.ConstraintResult) ErdRelationType {
	if constraint.IsPrimary && !constraint.HasMultiplePK {
		return relationOneToOne
	} else {
		return relationManyToOne
	}
}

func tableNameInSlice(slice []ErdTableData, tableName string) bool {
	for _, sliceItem := range slice {
		if sliceItem.Name == tableName {
			return true
		}
	}

	return false
}

func getAttributeKey(column database.ColumnResult) ErdAttributeKey {
	if column.IsPrimary {
		return primaryKey
	}

	if column.IsForeign {
		return foreignKey
	}

	return none
}

func getColumnData(config config.MermerdConfig, column database.ColumnResult) ErdColumnData {
	attributeKey := getAttributeKey(column)
	if config.OmitAttributeKeys() {
		attributeKey = none
	}

	var enumValues string
	if config.ShowEnumValues() {
		enumValues = column.EnumValues
	}

	return ErdColumnData{
		Name:         column.Name,
		DataType:     column.DataType,
		EnumValues:   enumValues,
		AttributeKey: attributeKey,
	}
}

func shouldSkipConstraint(config config.MermerdConfig, tables []ErdTableData, constraint database.ConstraintResult) bool {
	if config.ShowAllConstraints() {
		return false
	}

	// if config for all constraints is not set, only show constraints of selected tables
	return !(tableNameInSlice(tables, constraint.PkTable) && tableNameInSlice(tables, constraint.FkTable))
}

func getConstraintData(config config.MermerdConfig, constraint database.ConstraintResult) ErdConstraintData {
	constraintLabel := constraint.ColumnName
	if config.OmitConstraintLabels() {
		constraintLabel = ""
	}

	return ErdConstraintData{
		PkTableName:     constraint.PkTable,
		FkTableName:     constraint.FkTable,
		Relation:        getRelation(constraint),
		ConstraintLabel: constraintLabel,
	}
}
