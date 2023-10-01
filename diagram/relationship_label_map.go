package diagram

import (
	"fmt"

	"github.com/KarnerTh/mermerd/config"
)

type RelationshipLabelMap interface {
	AddRelationshipLabel(label config.RelationshipLabel)
	LookupRelationshipLabel(pkName, fkName string) (label config.RelationshipLabel, found bool)
}

type relationshipLabelMap struct {
	mapping map[string]config.RelationshipLabel
}

func (r *relationshipLabelMap) AddRelationshipLabel(label config.RelationshipLabel) {
	if r.mapping == nil {
		r.mapping = make(map[string]config.RelationshipLabel)
	}
	key := r.buildMapKey(label.PkName, label.FkName)
	r.mapping[key] = label
}

func (r *relationshipLabelMap) LookupRelationshipLabel(pkName, fkName string) (label config.RelationshipLabel, found bool) {
	if r.mapping == nil {
		return config.RelationshipLabel{}, false
	}
	key := r.buildMapKey(pkName, fkName)
	label, found = r.mapping[key]
	return
}

func (r *relationshipLabelMap) buildMapKey(pkName, fkName string) string {
	return fmt.Sprintf("%s-%s", pkName, fkName)
}

func BuildRelationshipLabelMapFromConfig(c config.MermerdConfig) RelationshipLabelMap {
	return BuildRelationshipLabelMap(c.RelationshipLabels())
}

func BuildRelationshipLabelMap(labels []config.RelationshipLabel) RelationshipLabelMap {
	labelMap := &relationshipLabelMap{}
	for _, label := range labels {
		labelMap.AddRelationshipLabel(label)
	}
	return labelMap
}
