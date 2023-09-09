package rrt

import (
	"fmt"
	"github.com/johnfercher/go-rrt/pkg/shared"
	"github.com/johnfercher/go-tree/tree"
	"math"
	"math/rand"
)

type CollisionCondition func(vector3D *shared.Vector3D) bool
type StopCondition func(testPoint *shared.Vector3D, finish *shared.Vector3D) bool

type RRT struct {
	collisionCondition CollisionCondition
	stopCondition      StopCondition
	stepDistance       float64
}

func New(stepDistance float64) *RRT {
	return &RRT{
		stepDistance: stepDistance,
	}
}

func (r *RRT) AddCollisionCondition(condition CollisionCondition) *RRT {
	r.collisionCondition = condition
	return r
}

func (r *RRT) AddStopCondition(condition StopCondition) *RRT {
	r.stopCondition = condition
	return r
}

func (r *RRT) FindPath(start *shared.Vector3D, finish *shared.Vector3D, world [][]*shared.Vector3D) []*tree.Node[*shared.Vector3D] {
	var nodesArray []*tree.Node[*shared.Vector3D]
	tr := tree.New[*shared.Vector3D]()

	maxDistance := shared.Distance(world[0][0], world[len(world)-1][len(world[0])-1])

	fmt.Printf("Max Distance: %f\n", maxDistance)

	nodeCounter := 0
	node := tree.NewNode(nodeCounter, start)
	tr.AddRoot(node)
	nodesArray = append(nodesArray, node)

	maxGenerations := 100000
	generation := 0
	minDistance := math.MaxFloat64
	var minDistancePoint *shared.Vector3D

	for generation < maxGenerations {
		minDistance = math.MaxFloat64
		generation++

		_, vector := nodesArray[len(nodesArray)-1].Get()
		if r.stopCondition(vector, finish) {
			break
		}

		newPoint := r.getRandomPoint(world)
		for _, point := range nodesArray {
			_, vector := point.Get()
			distance := shared.Distance(vector, newPoint)
			if distance < minDistance {
				minDistance = distance
				minDistancePoint = vector
			}
		}

		if minDistance < maxDistance*r.stepDistance && !r.collisionCondition(newPoint) {
			nodeCounter++
			newNode := tree.NewNode(nodeCounter, newPoint)
			//newNode.Print("New 1")
			tr.Add(nodeCounter-1, newNode)
			nodesArray = append(nodesArray, newNode)
			continue
		}

		fixedPoint := r.getFixedPoint(minDistancePoint, newPoint, world)

		if !r.collisionCondition(fixedPoint) {
			nodeCounter++
			newNode := tree.NewNode(nodeCounter, fixedPoint)
			//newNode.Print("New 2")
			tr.Add(nodeCounter-1, newNode)
			nodesArray = append(nodesArray, newNode)
		}
	}

	nodes, _ := tr.Backtrack(nodeCounter)
	return nodes
}

func (r *RRT) getRandomPoint(world [][]*shared.Vector3D) *shared.Vector3D {
	x := rand.Int() % len(world)
	y := rand.Int() & len(world[0])

	if x > len(world)-1 {
		x = len(world) - 1
	}

	if y > len(world[0])-1 {
		y = len(world[0]) - 1
	}
	//fmt.Printf("%d, %d\n", x, y)
	return world[x][y]
}

func (r *RRT) getFixedPoint(minDistancePoint *shared.Vector3D, newPoint *shared.Vector3D, world [][]*shared.Vector3D) *shared.Vector3D {
	radian := shared.Radian(minDistancePoint, newPoint)
	x := int(minDistancePoint.X + (math.Sin(radian)*minDistancePoint.X)*r.stepDistance*2)
	y := int(minDistancePoint.Y + (math.Cos(radian)*minDistancePoint.Y)*r.stepDistance*2)

	if x > len(world)-1 {
		x = len(world) - 1
	}

	if y > len(world[0])-1 {
		y = len(world[0]) - 1
	}

	return world[x][y]
}
