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

func ParseLabels(labels []string) []RelationshipLabel {
	var relationshipLabels []RelationshipLabel
	for _, label := range labels {
		parsed, err := ParseLabel(label)
		if err != nil {
			logrus.Warnf("label '%s' is not in the correct format", label)
			continue
		}
		relationshipLabels = append(relationshipLabels, parsed)
	}
	return relationshipLabels
}

func ParseLabel(label string) (RelationshipLabel, error) {
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

var labelRegex = regexp.MustCompile(`([\w\._-]+)[\s]+([\w\._-]+)[\s]+:[\s]+([\w._-]+)`)

func match(label string) (bool, [][]byte) {
	groups := labelRegex.FindSubmatch([]byte(label))
	if groups == nil {
		return false, [][]byte{}
	}
	return true, groups
}
