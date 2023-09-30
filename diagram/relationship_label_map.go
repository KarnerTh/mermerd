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
	key := r.buildMapKey(pkName, fkName)
	label, found = r.mapping[key]
	return
}

func (r *relationshipLabelMap) buildMapKey(pkName, fkName string) string {
	return fmt.Sprintf("%s-%s", pkName, fkName)
}

func BuildRelationshipLabelMap(c config.MermerdConfig) RelationshipLabelMap {
	labelMap := &relationshipLabelMap{}
	for _, label := range c.RelationshipLabels() {
		labelMap.AddRelationshipLabel(label)
	}
	return labelMap
}
