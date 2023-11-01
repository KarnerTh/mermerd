package diagram_test

import (
	"testing"

	"github.com/KarnerTh/mermerd/config"
	"github.com/KarnerTh/mermerd/diagram"
	"github.com/stretchr/testify/assert"
)

func TestEmptyRelationshipMapDoesNotError(t *testing.T) {
	relationshipMap := diagram.BuildRelationshipLabelMap([]config.RelationshipLabel{})

	_, found := relationshipMap.LookupRelationshipLabel("pk", "fk")

	assert.False(t, found)
}

func TestRelationshipMapCanAddAndLookupLabel(t *testing.T) {
	relationshipMap := diagram.BuildRelationshipLabelMap([]config.RelationshipLabel{})

	exampleLabel := config.RelationshipLabel{
		PkName: "name",
		FkName: "another-name",
		Label:  "a-label",
	}
	relationshipMap.AddRelationshipLabel(exampleLabel)

	actual, found := relationshipMap.LookupRelationshipLabel("name", "another-name")

	assert.True(t, found)
	assert.Equal(t, actual, exampleLabel)
}
