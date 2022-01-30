package analyzer

import (
	"github.com/AlecAivazis/survey/v2"
	"mermerd/config"
)

func ConnectionQuestion() survey.Prompt {
	return &survey.Input{
		Message: "Connection string",
		Suggest: func(toComplete string) []string {
			return config.ConnectionStringSuggestions
		},
	}
}

func SchemaQuestion(schemas []string) survey.Prompt {
	return &survey.Select{
		Message: "Choose a schema:",
		Options: schemas,
	}
}

func TableQuestion(tables []string) survey.Prompt {
	return &survey.MultiSelect{
		Message:  "Choose tables:",
		Options:  tables,
		PageSize: 15,
	}
}
