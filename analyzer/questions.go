package analyzer

import "github.com/AlecAivazis/survey/v2"

func ConnectionQuestion() survey.Prompt {
	return &survey.Input{
		Message: "Connection string",
		Suggest: func(toComplete string) []string {
			return []string{
				"postgresql://user:password@localhost:5432/dvdrental",
				"mysql://root:password@tcp(127.0.0.1:3306)/db",
			}
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
