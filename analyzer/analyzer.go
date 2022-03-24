package analyzer

import (
	"errors"
	"mermerd/config"
	"mermerd/database"
	"mermerd/util"
)

type analyzer struct {
	config           config.MermerdConfig
	connectorFactory database.ConnectorFactory
	questioner       Questioner
}

type Analyzer interface {
	Analyze() (*database.Result, error)
}

func NewAnalyzer(config config.MermerdConfig, connectorFactory database.ConnectorFactory, questioner Questioner) Analyzer {
	return analyzer{config, connectorFactory, questioner}
}

func (a analyzer) Analyze() (*database.Result, error) {
	loading, err := util.NewLoadingSpinner()
	if err != nil {
		return nil, err
	}

	connectionString := a.config.ConnectionString()
	if connectionString == "" {
		connectionString, err = a.questioner.AskConnectionQuestion(a.config.ConnectionStringSuggestions())
		if err != nil {
			return nil, err
		}
	}

	loading.Start("Connecting to database and getting schemas")
	db, err := a.connectorFactory.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectedSchema := a.config.Schema()
	if selectedSchema == "" {
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
			selectedSchema, err = a.questioner.AskSchemaQuestion(schemas)
			if err != nil {
				return nil, err
			}
		}
	}

	// get tables
	var selectedTables []string
	loading.Start("Getting tables")
	tables, err := db.GetTables(selectedSchema)
	if err != nil {
		return nil, err
	}
	loading.Stop()

	if a.config.UseAllTables() {
		selectedTables = tables
	} else if len(a.config.SelectedTables()) > 0 {
		selectedTables = a.config.SelectedTables()
	} else {
		selectedTables, err = a.questioner.AskTableQuestion(tables)
		if err != nil {
			return nil, err
		}
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
