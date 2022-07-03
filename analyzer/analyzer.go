package analyzer

import (
	"errors"
	"github.com/sirupsen/logrus"

	"github.com/KarnerTh/mermerd/config"
	"github.com/KarnerTh/mermerd/database"
	"github.com/KarnerTh/mermerd/util"
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
		logrus.Error("Getting schemas failed", " | ", err)
		return "", err
	}

	logrus.WithField("count", len(schemas)).Info("Got schemas")

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
		return selectedTables, nil
	}

	a.loadingSpinner.Start("Getting tables")
	tables, err := db.GetTables(selectedSchema)
	a.loadingSpinner.Stop()
	if err != nil {
		logrus.Error("Getting tables failed", " | ", err)
		return nil, err
	}

	logrus.WithField("count", len(tables)).Info("Got tables")

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
			logrus.Error("Getting columns failed", " | ", err)
			return nil, err
		}

		constraints, err := db.GetConstraints(table)
		if err != nil {
			logrus.Error("Getting constraints failed", " | ", err)
			return nil, err
		}

		tableResults = append(tableResults, database.TableResult{TableName: table, Columns: columns, Constraints: constraints})
	}
	a.loadingSpinner.Stop()
	columnCount, constraintCount := getTableResultStats(tableResults)
	logrus.WithFields(logrus.Fields{"columns": columnCount, "constraints": constraintCount}).Info("Got columns and constraints")
	return tableResults, nil
}

func getTableResultStats(tableResults []database.TableResult) (columnCount int, constraintCount int) {
	for _, tableResult := range tableResults {
		columnCount += len(tableResult.Columns)
		constraintCount += len(tableResult.Constraints)
	}

	return columnCount, constraintCount
}
