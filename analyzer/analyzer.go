package analyzer

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"mermerd/config"
	"mermerd/database"
	"mermerd/util"
)

func Analyze() (*database.Result, error) {
	loading, err := util.NewLoadingSpinner()
	if err != nil {
		return nil, err
	}

	var connectionString string
	if config.ConnectionString == "" {
		err = survey.AskOne(ConnectionQuestion(), &connectionString, survey.WithValidator(survey.Required))
		if err != nil {
			return nil, err
		}
	} else {
		connectionString = config.ConnectionString
	}

	loading.Start("Connecting to database and getting schemas")
	db, err := database.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var selectedSchema string
	if config.Schema == "" {
		schemas, err := db.GetSchemas()
		if err != nil {
			return nil, err
		}
		loading.Stop()

		switch len(schemas) {
		case 0:
			return nil, errors.New("no schemas available")
		case 1:
			selectedSchema = schemas[0]
			break
		default:
			err = survey.AskOne(SchemaQuestion(schemas), &selectedSchema)
			if err != nil {
				return nil, err
			}
		}
	} else {
		selectedSchema = config.Schema
	}

	// get tables
	var selectedTables []string
	loading.Start("Getting tables")
	tables, err := db.GetTables(selectedSchema)
	if err != nil {
		return nil, err
	}
	loading.Stop()

	err = survey.AskOne(TableQuestion(tables), &selectedTables, survey.WithValidator(survey.MinItems(1)))
	if err != nil {
		return nil, err
	}

	// get columns and constraints
	var tableResults []database.TableResult
	loading.Start("Getting columns and constraints")
	for _, table := range selectedTables {
		columns, err := db.GetColumns(table)
		if err != nil {
			return nil, err
		}

		constraints, err := db.GetConstraints(table)
		if err != nil {
			return nil, err
		}

		tableResults = append(tableResults, database.TableResult{TableName: table, Columns: columns, Constraints: constraints})
	}
	loading.Stop()

	return &database.Result{Tables: tableResults}, nil
}
