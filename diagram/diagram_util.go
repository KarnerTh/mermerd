package diagram

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

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

func getAttributeKeys(column database.ColumnResult) []ErdAttributeKey {
	var attributeKeys []ErdAttributeKey
	if column.IsPrimary {
		attributeKeys = append(attributeKeys, primaryKey)
	}

	if column.IsForeign {
		attributeKeys = append(attributeKeys, foreignKey)
	}

	return attributeKeys
}

func getColumnData(config config.MermerdConfig, column database.ColumnResult) ErdColumnData {
	attributeKeys := getAttributeKeys(column)
	if config.OmitAttributeKeys() {
		attributeKeys = []ErdAttributeKey(nil)
	}

	return ErdColumnData{
		Name:          column.Name,
		DataType:      column.DataType,
		Description:   getDescription(config.ShowDescriptions(), column),
		AttributeKeys: attributeKeys,
	}
}

func getDescription(options []string, column database.ColumnResult) string {
	var description []string
	for _, option := range options {
		switch option {
		case "notNull":
			if !column.IsNullable {
				description = append(description, "{NOT_NULL}")
			}
		case "enumValues":
			if column.EnumValues != "" {
				description = append(description, "<"+column.EnumValues+">")
			}
		case "columnComments":
			description = append(description, escapeComments(column.Comment))
		default:
			logrus.Errorf("Could not parse option %q", option)
		}
	}
	return strings.TrimSpace(strings.Join(description, " "))
}

// escapeComments escapes invalid characters in comments.
// mermaid does not support quote marks ("), as it uses this to mark the comment, as such these need to be replaced
// with `#quot;`.
func escapeComments(str string) string {
	return strings.ReplaceAll(str, `"`, "#quot;")
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
		PkTableName:     getTableName(config, database.TableDetail{Schema: constraint.PkSchema, Name: constraint.PkTable}),
		FkTableName:     getTableName(config, database.TableDetail{Schema: constraint.FkSchema, Name: constraint.FkTable}),
		Relation:        getRelation(constraint),
		ConstraintLabel: constraintLabel,
	}
}

func getTableName(config config.MermerdConfig, table database.TableDetail) string {
	if !config.ShowSchemaPrefix() {
		return table.Name
	}

	separator := config.SchemaPrefixSeparator()
	name := fmt.Sprintf("%s%s%s", table.Schema, separator, table.Name)

	// if fullstop is used the table name needs to be escaped with quote marks
	if separator == "." {
		return fmt.Sprintf("\"%s\"", name)
	}

	return name
}
