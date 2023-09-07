package tree

import (
	"fmt"
	"github.com/johnfercher/rrt/pkg/shared"
)

type Tree struct {
	root      *Node
	lastAdded *Node
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) Root() (*shared.Vector3D, bool) {
	if t.root == nil {
		return nil, false
	}

	return t.root.Vector3D, true
}

func (t *Tree) AddRoot(node *Node) bool {
	if t.root == nil {
		//fmt.Printf("Add root(%d)\n", id)
		t.root = node
		return true
	}

	return false
}

func (t *Tree) Add(parentID int, node *Node) bool {
	if t.root == nil {
		return false
	}

	return t.add(parentID, t.root, node)
}

func (t *Tree) Backtrack() []*Node {
	return t.lastAdded.Backtrack()
}

func (t *Tree) Print() {
	if t.root == nil {
		fmt.Println("empty tree")
		return
	}

}

func (t *Tree) add(parentID int, node *Node, newNode *Node) bool {
	if parentID == node.ID {
		newNode.Previous = node
		//fmt.Printf("Add next(%d) to(%d)\n", newNode.ID, node.ID)
		//newNode.Print()
		t.lastAdded = newNode
		node.Nexts = append(node.Nexts, newNode)
		return true
	}

	for _, next := range node.Nexts {
		added := t.add(parentID, next, newNode)
		if added {
			return true
		}
	}

	return false
}

func (t *Tree) findNode(current *Node, parentID int, depth int) (*Node, *Node, bool) {
	if current.ID == parentID {
		return current, nil, true
	}

	for _, next := range current.Nexts {
		if next.ID == parentID {
			return current, next, true
		}
	}

	for _, next := range current.Nexts {
		parent, node, found := t.findNode(next, parentID, depth+1)
		if found {
			return parent, node, true
		}
	}

	return nil, nil, false
}
