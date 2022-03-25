package diagram

import (
	"bufio"
	"fmt"
	"mermerd/config"
	"mermerd/database"
	"os"
	"strings"
)

const (
	relationOneToOne  = "|o--||"
	relationManyToOne = "}o--||"
)

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
		return err
	}

	defer f.Close()

	buffer := bufio.NewWriter(f)
	if d.config.EncloseWithMermaidBackticks() {
		if _, err = buffer.WriteString("```mermaid\n"); err != nil {
			return err
		}
	}

	if _, err = buffer.WriteString("erDiagram\n"); err != nil {
		return err
	}

	var tableNames []string
	var allConstraints database.ConstraintResultList
	for _, table := range result.Tables {
		tableNames = append(tableNames, table.TableName)
		allConstraints = allConstraints.AppendIfNotExists(table.Constraints...)
	}

	for _, table := range result.Tables {
		if _, err := buffer.WriteString(fmt.Sprintf("    %s {\n", table.TableName)); err != nil {
			return err
		}

		for _, column := range table.Columns {
			if _, err := buffer.WriteString(fmt.Sprintf("        %s %s\n", column.DataType, column.Name)); err != nil {
				return err
			}
		}

		if _, err = buffer.WriteString("    }"); err != nil {
			return err
		}

		if _, err = buffer.WriteString("\n\n"); err != nil {
			return err
		}
	}

	constraints := strings.Builder{}
	for _, constraint := range allConstraints {
		if (!sliceContainsItem(tableNames, constraint.PKTable) || !sliceContainsItem(tableNames, constraint.FkTable)) && !d.config.ShowAllConstraints() {
			continue
		}

		relation := getRelation(constraint)
		constraints.WriteString(fmt.Sprintf("    %s %s %s : \"\"\n", constraint.FkTable, relation, constraint.PKTable))
	}

	if _, err = buffer.WriteString(constraints.String()); err != nil {
		return err
	}

	if d.config.EncloseWithMermaidBackticks() {
		_, err = buffer.WriteString("```\n")
	}

	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}

func getRelation(constraint database.ConstraintResult) string {
	if constraint.IsPrimary && !constraint.HasMultiplePK {
		return relationOneToOne
	} else {
		return relationManyToOne
	}
}

func sliceContainsItem(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}

	return false
}
