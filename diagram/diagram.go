package diagram

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/KarnerTh/mermerd/config"
	"github.com/KarnerTh/mermerd/database"
)

//go:embed erd_template.gommd
var erdTemplate string

type diagram struct {
	config config.MermerdConfig
}

type Diagram interface {
	Create(result *database.Result) error
}

func NewDiagram(config config.MermerdConfig) Diagram {
	return diagram{config}
}

func (d diagram) Create(result *database.Result) error {
	f, err := os.Create(d.config.OutputFileName())
	if err != nil {
		logrus.Error("Could not create output file", " | ", err)
		return err
	}

	defer f.Close()

	tmpl, err := template.New("erd_template").Parse(erdTemplate)
	if err != nil {
		logrus.Error("Could not load template file", " | ", err)
		return err
	}

	tableData := make([]ErdTableData, len(result.Tables))
	var allConstraints database.ConstraintResultList

	for tableIndex, table := range result.Tables {
		allConstraints = allConstraints.AppendIfNotExists(table.Constraints...)

		columnData := make([]ErdColumnData, len(table.Columns))
		for columnIndex, column := range table.Columns {
			attributeKey := getAttributeKey(column)
			if d.config.OmitAttributeKeys() {
				attributeKey = none
			}

			var enumValues string
			if d.config.ShowEnumValues() {
				enumValues = column.EnumValues
			}

			columnData[columnIndex] = ErdColumnData{
				Name:         column.Name,
				DataType:     column.DataType,
				EnumValues:   enumValues,
				AttributeKey: attributeKey,
			}
		}

		tableData[tableIndex] = ErdTableData{
			Name:    table.Table.Name,
			Columns: columnData,
		}
	}

	var constraints []ErdConstraintData
	for _, constraint := range allConstraints {
		if (!tableNameInSlice(tableData, constraint.PkTable) || !tableNameInSlice(tableData, constraint.FkTable)) && !d.config.ShowAllConstraints() {
			continue
		}

		constraintLabel := constraint.ColumnName
		if d.config.OmitConstraintLabels() {
			constraintLabel = ""
		}

		constraints = append(constraints, ErdConstraintData{
			PkTableName:     constraint.PkTable,
			FkTableName:     constraint.FkTable,
			Relation:        getRelation(constraint),
			ConstraintLabel: constraintLabel,
		})
	}

	diagramData := ErdDiagramData{
		EncloseWithMermaidBackticks: d.config.EncloseWithMermaidBackticks(),
		Tables:                      tableData,
		Constraints:                 constraints,
	}

	if err = tmpl.Execute(f, diagramData); err != nil {
		logrus.Error("Could not create diagram", " | ", err)
		return err
	}
	return nil
}

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
