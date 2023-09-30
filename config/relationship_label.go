package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type RelationshipLabel struct {
	PkName string
	FkName string
	Label  string
}

func parseLabels(labels []string) []RelationshipLabel {
	var relationshipLabels []RelationshipLabel
	for _, label := range labels {
		parsed, err := parseLabel(label)
		if err != nil {
			logrus.Warnf("label '%s' is not in the correct format", label)
			continue
		}
		relationshipLabels = append(relationshipLabels, parsed)
	}
	return relationshipLabels
}

func parseLabel(label string) (RelationshipLabel, error) {
	label = strings.Trim(label, " \t")
	matched, groups := match(label)
	if !matched {
		return RelationshipLabel{}, errors.New("invalid relationship label")
	}

	return RelationshipLabel{
		PkName: string(groups[1]),
		FkName: string(groups[2]),
		Label:  string(groups[3]),
	}, nil
}

// The regex works by creating three capture groups
// Each group allows for all word characters, `.`, `_` and `-` any number of times
// The first two groups (the table names) are separated by any amount of whitespace characters
// The table names and label are are separated by
// - any number of whitespace characters
// - a `:`
// - and then any other number of whitespace characters
// The string must start with the first table name and it must end with the label
var labelRegex = regexp.MustCompile(`^([\w\._-]+)[\s]+([\w\._-]+)[\s]+:[\s]+([\w._-]+)$`)

func match(label string) (bool, [][]byte) {
	groups := labelRegex.FindSubmatch([]byte(label))
	if groups == nil {
		return false, [][]byte{}
	}
	return true, groups
}
