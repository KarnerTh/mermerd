package diagram

import (
	_ "embed"
	"io"
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
	Create(wr io.Writer, result *database.Result) error
}

func NewDiagram(config config.MermerdConfig) Diagram {
	return diagram{config}
}

func (d diagram) Create(wr io.Writer, result *database.Result) error {
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
			columnData[columnIndex] = getColumnData(d.config, column)
		}

		tableData[tableIndex] = ErdTableData{
			Name:    getTableName(d.config, table.Table),
			Columns: columnData,
		}
	}

	var constraints []ErdConstraintData
	relationshipLabelMap := BuildRelationshipLabelMapFromConfig(d.config)
	for _, constraint := range allConstraints {
		if shouldSkipConstraint(d.config, tableData, constraint) {
			continue
		}

		constraints = append(constraints, getConstraintData(d.config, relationshipLabelMap, constraint))
	}

	diagramData := ErdDiagramData{
		EncloseWithMermaidBackticks: d.config.EncloseWithMermaidBackticks(),
		Tables:                      tableData,
		Constraints:                 constraints,
	}

	if err = tmpl.Execute(wr, diagramData); err != nil {
		logrus.Error("Could not create diagram", " | ", err)
		return err
	}
	return nil
}
