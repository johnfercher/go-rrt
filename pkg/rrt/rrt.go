package rrt

import (
	"fmt"
	"github.com/johnfercher/go-tree/tree"
	"math"
	"math/rand"
)

type RRT[T any] struct {
	collisionCondition func(point T) bool
	stopCondition      func(testPoint *Point[T], finish *Point[T]) bool
	stepDistance       float64
	maxTries           int
	startPoint         *Point[T]
	finishPoint        *Point[T]
}

func New[T any](stepDistance float64, maxTries int) *RRT[T] {
	return &RRT[T]{
		stepDistance: stepDistance,
		maxTries:     maxTries,
	}
}

func (r *RRT[T]) AddCollisionCondition(condition func(point T) bool) *RRT[T] {
	r.collisionCondition = condition
	return r
}

func (r *RRT[T]) AddStopCondition(condition func(testPoint *Point[T], finish *Point[T]) bool) *RRT[T] {
	r.stopCondition = condition
	return r
}

func (r *RRT[T]) FindPath(start *Coordinate, finish *Coordinate, world [][]T) []*Point[T] {
	r.startPoint = &Point[T]{
		X: start.X,
		Y: start.Y,
	}
	r.finishPoint = &Point[T]{
		X: finish.X,
		Y: finish.Y,
	}

	nodes := make(map[string]*tree.Node[*Point[T]])
	var lastNodeAdded *tree.Node[*Point[T]]
	tr := tree.New[*Point[T]]()

	maxDistance := Distance(&Point[T]{X: 0, Y: 0}, &Point[T]{X: float64(len(world) - 1), Y: float64(len(world[0]) - 1)})

	fmt.Printf("Max Distance: %f\n", maxDistance)

	nodeCounter := 0
	node := tree.NewNode[*Point[T]](nodeCounter, r.startPoint)
	tr.AddRoot(node)
	_, point := node.Get()
	nodes[r.getKey(point)] = node
	lastNodeAdded = node

	try := 0
	minDistance := math.MaxFloat64
	var minDistancePoint *Point[T]

	goTofinish := true
	for try < r.maxTries {
		minDistance = math.MaxFloat64
		minDistancePoint = nil
		try++
		goTofinish = !goTofinish

		_, lastAdded := lastNodeAdded.Get()
		if r.stopCondition(lastAdded, r.finishPoint) {
			break
		}

		newPoint := r.getRandomPoint(world, goTofinish)
		for _, point := range nodes {
			_, vector := point.Get()
			distance := Distance(vector, newPoint)
			if distance < minDistance {
				minDistance = distance
				minDistancePoint = vector
			}
		}

		fixedPoint := r.getFixedPoint(minDistancePoint, newPoint, world)

		if !r.collisionCondition(fixedPoint.Data) {
			key := r.getKey(fixedPoint)
			if _, ok := nodes[key]; ok {
				continue
			}

			nodeCounter++
			newNode := tree.NewNode(nodeCounter, fixedPoint)
			ok := tr.Add(nodeCounter-1, newNode)
			if !ok {
				fmt.Println("Could not add to tree")
			}

			nodes[key] = newNode
			lastNodeAdded = newNode
		}
	}

	pathNodes, _ := tr.Backtrack(nodeCounter)
	var points []*Point[T]
	for _, pathNode := range pathNodes {
		_, data := pathNode.Get()
		points = append(points, data)
	}

	return points
}

func (r *RRT[T]) getRandomPoint(world [][]T, goToFinish bool) *Point[T] {
	if goToFinish {
		return r.finishPoint
	}

	x := rand.Int() % len(world)
	y := rand.Int() & len(world[0])

	if x > len(world)-1 {
		x = len(world) - 1
	}

	if y > len(world[0])-1 {
		y = len(world[0]) - 1
	}

	return &Point[T]{
		X:    float64(x),
		Y:    float64(y),
		Data: world[x][y],
	}
}

func (r *RRT[T]) getFixedPoint(minDistancePoint *Point[T], newPoint *Point[T], world [][]T) *Point[T] {

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
	return &Point[T]{
		X:    float64(x),
		Y:    float64(y),
		Data: world[x][y],
	}
}

func (r *RRT[T]) getKey(point *Point[T]) string {
	return fmt.Sprintf("%d-%d", int(point.X), int(point.Y))
}
