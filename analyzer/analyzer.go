package analyzer

import (
	"errors"
	"mermerd/config"
	"mermerd/database"
	"mermerd/util"
)

type analyzer struct {
	loadingSpinner   util.LoadingSpinner
	config           config.MermerdConfig
	connectorFactory database.ConnectorFactory
	questioner       Questioner
}

type Analyzer interface {
	Analyze() (*database.Result, error)

	GetConnectionString() (string, error)
	GetSchema(db database.Connector) (string, error)
	GetTables(db database.Connector, selectedSchema string) ([]string, error)
	GetColumnsAndConstraints(db database.Connector, selectedTables []string) ([]database.TableResult, error)
}

func NewAnalyzer(config config.MermerdConfig, connectorFactory database.ConnectorFactory, questioner Questioner) Analyzer {
	loadingSpinner := util.NewLoadingSpinner()
	return analyzer{loadingSpinner, config, connectorFactory, questioner}
}

func (a analyzer) Analyze() (*database.Result, error) {
	connectionString, err := a.GetConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := a.connectorFactory.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}

	a.loadingSpinner.Start("Connecting to database")
	if err = db.Connect(); err != nil {
		return nil, err
	}
	defer db.Close()
	a.loadingSpinner.Stop()

	selectedSchema, err := a.GetSchema(db)
	if err != nil {
		return nil, err
	}

	selectedTables, err := a.GetTables(db, selectedSchema)
	if err != nil {
		return nil, err
	}

	tableResults, err := a.GetColumnsAndConstraints(db, selectedTables)
	if err != nil {
		return nil, err
	}

	return &database.Result{Tables: tableResults}, nil
}

func (a analyzer) GetConnectionString() (string, error) {
	if connectionString := a.config.ConnectionString(); connectionString != "" {
		return connectionString, nil
	}

	return a.questioner.AskConnectionQuestion(a.config.ConnectionStringSuggestions())
}

func (a analyzer) GetSchema(db database.Connector) (string, error) {
	if selectedSchema := a.config.Schema(); selectedSchema != "" {
		return selectedSchema, nil
	}

	a.loadingSpinner.Start("Getting schemas")
	schemas, err := db.GetSchemas()
	a.loadingSpinner.Stop()

	if err != nil {
		return "", err
	}

	switch len(schemas) {
	case 0:
		return "", errors.New("no schemas available")
	case 1:
		return schemas[0], nil
	default:
		return a.questioner.AskSchemaQuestion(schemas)
	}
}

func (a analyzer) GetTables(db database.Connector, selectedSchema string) ([]string, error) {
	if selectedTables := a.config.SelectedTables(); len(selectedTables) > 0 {
		tableNames := make([]string, 0, len(selectedTables))

		for tableName := range selectedTables {
			tableNames = append(tableNames, tableName)
		}

		return tableNames, nil
	}

	a.loadingSpinner.Start("Getting tables")
	tables, err := db.GetTables(selectedSchema)
	if err != nil {
		return nil, err
	}
	a.loadingSpinner.Stop()

	if a.config.UseAllTables() {
		return tables, nil
	} else {
		return a.questioner.AskTableQuestion(tables)
	}
}

func (a analyzer) GetColumnsAndConstraints(db database.Connector, selectedTables []string) ([]database.TableResult, error) {
	var tableResults []database.TableResult
	a.loadingSpinner.Start("Getting columns and constraints")
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
	a.loadingSpinner.Stop()
	return tableResults, nil
}
