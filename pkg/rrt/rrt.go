package rrt

import (
	"fmt"
	"github.com/johnfercher/go-tree/tree"
	"math"
	"math/rand"
)

type RRT[T any] struct {
	collisionCondition func(point *Point[T]) bool
	stopCondition      func(testPoint *Point[T], finish *Point[T]) bool
	stepDistance       float64
	start              *Point[T]
	finish             *Point[T]
}

func New[T any](stepDistance float64) *RRT[T] {
	return &RRT[T]{
		stepDistance: stepDistance,
	}
}

func (r *RRT[T]) AddCollisionCondition(condition func(point *Point[T]) bool) *RRT[T] {
	r.collisionCondition = condition
	return r
}

func (r *RRT[T]) AddStopCondition(condition func(testPoint *Point[T], finish *Point[T]) bool) *RRT[T] {
	r.stopCondition = condition
	return r
}

func (r *RRT[T]) FindPath(start *Point[T], finish *Point[T], world [][]*Point[T]) []*Point[T] {
	r.start = start
	r.finish = finish
	var nodesArray []*tree.Node[*Point[T]]
	tr := tree.New[*Point[T]]()

	maxDistance := Distance(world[0][0], world[len(world)-1][len(world[0])-1])

	fmt.Printf("Max Distance: %f\n", maxDistance)

	nodeCounter := 0
	node := tree.NewNode(nodeCounter, start)
	tr.AddRoot(node)
	nodesArray = append(nodesArray, node)

	maxGenerations := 100000
	generation := 0
	minDistance := math.MaxFloat64
	var minDistancePoint *Point[T]

	goTofinish := true
	for generation < maxGenerations {
		minDistance = math.MaxFloat64
		minDistancePoint = nil
		generation++
		goTofinish = !goTofinish

		_, lastAdded := nodesArray[len(nodesArray)-1].Get()
		if r.stopCondition(lastAdded, finish) {
			break
		}

		newPoint := r.getRandomPoint(world, goTofinish)
		for _, point := range nodesArray {
			_, vector := point.Get()
			distance := Distance(vector, newPoint)
			if distance < minDistance {
				minDistance = distance
				minDistancePoint = vector
			}
		}

		fixedPoint := r.getFixedPoint(minDistancePoint, newPoint, world)

		if !r.collisionCondition(fixedPoint) {
			nodeCounter++
			newNode := tree.NewNode(nodeCounter, fixedPoint)
			ok := tr.Add(nodeCounter-1, newNode)
			if !ok {
				fmt.Println("Could not add to tree")
			}
			nodesArray = append(nodesArray, newNode)
		}
	}

	nodes, _ := tr.Backtrack(nodeCounter)
	var points []*Point[T]
	for _, node := range nodes {
		_, data := node.Get()
		points = append(points, data)
	}

	return points
}

func (r *RRT[T]) getRandomPoint(world [][]*Point[T], goToFinish bool) *Point[T] {
	if goToFinish {
		return r.finish
	}

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

func (r *RRT[T]) getFixedPoint(minDistancePoint *Point[T], newPoint *Point[T], world [][]*Point[T]) *Point[T] {

	radian := Radian(minDistancePoint, newPoint)
	x := int(minDistancePoint.X + (math.Sin(radian)*minDistancePoint.X)*r.stepDistance)
	y := int(minDistancePoint.Y + (math.Cos(radian)*minDistancePoint.Y)*r.stepDistance)

	if x > len(world)-1 {
		x = len(world) - 1
	}

	if y > len(world[0])-1 {
		y = len(world[0]) - 1
	}

	//fmt.Printf("Min %s, New %s, Fix %s\n", minDistancePoint.GetString(), newPoint.GetString(), world[x][y].GetString())

	return world[x][y]
}
