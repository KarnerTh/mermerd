package diagram

import (
	"bufio"
	"fmt"
	"mermerd/config"
	"mermerd/database"
	"os"
	"strings"
)

func Create(result *database.Result) error {
	f, err := os.Create("result.mmd")
	if err != nil {
		return err
	}

	defer f.Close()

	buffer := bufio.NewWriter(f)
	_, err = buffer.WriteString("erDiagram\n")
	if err != nil {
		return err
	}

	var tableNames []string
	var allConstraints []database.ConstraintResult
	for _, table := range result.Tables {
		tableNames = append(tableNames, table.TableName)
		allConstraints = appendConstraintsIfNotExists(allConstraints, table.Constraints...)
	}

	for _, table := range result.Tables {
		_, err := buffer.WriteString(fmt.Sprintf("    %s {\n", table.TableName))
		if err != nil {
			return err
		}

		for _, column := range table.Columns {
			_, err := buffer.WriteString(fmt.Sprintf("        %s %s\n", column.DataType, column.Name))
			if err != nil {
				return err
			}
		}
		_, err = buffer.WriteString("    }")
		if err != nil {
			return err
		}
		_, err = buffer.WriteString("\n\n")

	}

	constraints := strings.Builder{}
	for _, constraint := range allConstraints {
		if (!sliceContainsItem(tableNames, constraint.PKTable) || !sliceContainsItem(tableNames, constraint.FkTable)) && !config.ShowAllConstraints {
			continue
		}

		relation := getRelation(constraint)
		constraints.WriteString(fmt.Sprintf("    %s %s %s : \"\"\n", constraint.FkTable, relation, constraint.PKTable))
	}

	_, err = buffer.WriteString(constraints.String())
	if err != nil {
		return err
	}

	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}

func getRelation(constraint database.ConstraintResult) string {
	if constraint.IsPrimary {
		return "|o--||"
	} else {
		return "}o--||"
	}
}

// ensure that only unique items are appended to the list of constraints
func appendConstraintsIfNotExists(list []database.ConstraintResult, items ...database.ConstraintResult) []database.ConstraintResult {
	result := list
	for _, item := range items {
		if !sliceContainsConstraint(result, item) {
			result = append(result, item)
		}
	}

	return result
}

func sliceContainsConstraint(slice []database.ConstraintResult, item database.ConstraintResult) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}

	return false
}

func sliceContainsItem(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}

	return false
}
