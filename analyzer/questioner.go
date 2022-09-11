package analyzer

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type questioner struct{}

type Questioner interface {
	AskConnectionQuestion(suggestions []string) (string, error)
	AskSchemaQuestion(schemas []string) (string, error)
	AskTableQuestion(tables []string) ([]string, error)
}

func NewQuestioner() Questioner {
	return questioner{}
}

func (q questioner) AskConnectionQuestion(suggestions []string) (string, error) {
	var result string
	question := survey.Input{
		Message: "Connection string:",
		Suggest: func(toComplete string) []string {
			return suggestions
		},
	}

	err := survey.AskOne(&question, &result, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return os.ExpandEnv(result), nil
}

func (q questioner) AskSchemaQuestion(schemas []string) (string, error) {
	var result string
	question := &survey.Select{
		Message: "Choose a schema:",
		Options: schemas,
	}

	err := survey.AskOne(question, &result)
	return result, err
}

func (q questioner) AskTableQuestion(tables []string) ([]string, error) {
	var result []string
	question := &survey.MultiSelect{
		Message:  "Choose tables:",
		Options:  tables,
		PageSize: 15,
	}

	err := survey.AskOne(question, &result, survey.WithValidator(survey.MinItems(1)))
	return result, err
}
