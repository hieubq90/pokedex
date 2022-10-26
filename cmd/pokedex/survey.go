package pokedex

import "github.com/AlecAivazis/survey/v2"

type FileTypeAnswer struct {
	FileType string `survey:"file_type"`
}

var fileTypeQuestion = []*survey.Question{
	{
		Name: "file_type",
		Prompt: &survey.Select{
			Message: "Choose a file type:",
			Options: []string{"csv", "json"},
			Default: "json",
		},
	},
}
