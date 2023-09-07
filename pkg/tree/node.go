package tree

import (
	"fmt"
	"github.com/johnfercher/rrt/pkg/shared"
)

type Node struct {
	ID       int
	Vector3D *shared.Vector3D
	Previous *Node
	Nexts    []*Node
}

func NewNode(ID int, vector3D *shared.Vector3D) *Node {
	return &Node{
		ID:       ID,
		Vector3D: vector3D,
	}
}

func (n *Node) IsLeaf() bool {
	if n == nil {
		return false
	}

	return len(n.Nexts) == 0
}

func (n *Node) Backtrack() []*Node {
	var nodes []*Node
	//n.Print()

	current := n
	for current != nil {
		nodes = append(nodes, current)
		current = current.Previous
	}

	return nodes
}

func (n *Node) Print(label string) {
	if n.Previous == nil {
		fmt.Printf("%s - ID(%d), Vector(%f, %f, %f), NextSize(%d), HasPrevious(false)\n", label, n.ID, n.Vector3D.X, n.Vector3D.Y, n.Vector3D.Z, len(n.Nexts))
	} else {
		fmt.Printf("%s - ID(%d), Vector(%f, %f, %f), NextSize(%d), Previous(%d)\n", label, n.ID, n.Vector3D.X, n.Vector3D.Y, n.Vector3D.Z, len(n.Nexts), n.Previous.ID)
	}
}
