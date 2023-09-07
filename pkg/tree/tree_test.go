package tree

import (
	"github.com/johnfercher/rrt/pkg/shared"
	"testing"
)

func TestTree_Add(t *testing.T) {
	// Arrange
	tree := New()

	// Act
	tree.AddRoot(NewNode(0, &shared.Vector3D{}))
	tree.Add(0, NewNode(1, &shared.Vector3D{}))
	tree.Add(1, NewNode(3, &shared.Vector3D{}))
	tree.Add(1, NewNode(4, &shared.Vector3D{}))
	tree.Add(4, NewNode(5, &shared.Vector3D{}))

	nodes := tree.Backtrack()
	for _, node := range nodes {
		node.Print()
	}
}
