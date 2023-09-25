package rrt

import (
	"fmt"
	"github.com/johnfercher/go-tree/node"
	"github.com/johnfercher/go-tree/tree"
	builtInMath "math"
	"math/rand"
)

type rrt[T any] struct {
	collisionCondition    func(point T) bool
	stopCondition         func(testPoint *Point[T], finish *Point[T]) bool
	stepDistance          int
	maxTries              int
	focusOnFinishEveryTry int
	StartPoint            *Point[T]
	FinishPoint           *Point[T]
	tree                  *tree.Tree[*Point[T]]
	nodeCounter           int
}

func New[T any](stepDistance int, maxTries int, focusOnFinishEveryTry int) *rrt[T] {
	return &rrt[T]{
		stepDistance:          stepDistance,
		maxTries:              maxTries,
		focusOnFinishEveryTry: focusOnFinishEveryTry,
		tree:                  tree.New[*Point[T]](),
	}
}

func (r *rrt[T]) AddCollisionCondition(condition func(point T) bool) {
	r.collisionCondition = condition
}

func (r *rrt[T]) AddStopCondition(condition func(testPoint *Point[T], finish *Point[T]) bool) {
	r.stopCondition = condition
}

func (r *rrt[T]) GetTree() *tree.Tree[*Point[T]] {
	return r.tree
}

func (r *rrt[T]) FindPath(start *Point[any], finish *Point[any], world [][]T) []*Point[T] {
	nodeCounter := r.findLastNode(start, finish, world)

	pathNodes, _ := r.tree.Backtrack(nodeCounter)
	var points []*Point[T]
	for _, pathNode := range pathNodes {
		data := pathNode.GetData()
		points = append(points, data)
	}

	return points
}

func (r *rrt[T]) findLastNode(start *Point[any], finish *Point[any], world [][]T) int {
	r.StartPoint = &Point[T]{
		X:    start.X,
		Y:    start.Y,
		Data: world[start.X][start.Y],
	}
	r.FinishPoint = &Point[T]{
		X:    finish.X,
		Y:    finish.Y,
		Data: world[finish.X][finish.Y],
	}

	nodes := make(map[string]*node.Node[*Point[T]])
	var lastNodeAdded *node.Node[*Point[T]]

	nodeCounter := 0
	n := node.New[*Point[T]](r.StartPoint).WithID(nodeCounter)
	r.tree.AddRoot(n)
	point := n.GetData()
	key := r.GetKey(point.X, point.Y)
	nodes[key] = n

	lastNodeAdded = n

	try := 0
	minDistance := builtInMath.MaxFloat64
	var minNode *node.Node[*Point[T]]

	for try < r.maxTries {
		minDistance = builtInMath.MaxFloat64
		minNode = nil
		try++

		lastAdded := lastNodeAdded.GetData()
		if r.stopCondition(lastAdded, r.FinishPoint) {
			break
		}

		newPoint, ok := r.GetRandomPoint(world, try)
		if !ok {
			continue
		}

		for _, point := range nodes {
			vector := point.GetData()
			distance := vector.DistanceTo(newPoint)
			if distance < minDistance {
				minDistance = distance
				minNode = point
			}
		}

		minID := minNode.GetID()
		minPoint := minNode.GetData()
		fixedPoint := r.GetFixedPoint(minPoint, newPoint, world)

		if !r.collisionCondition(fixedPoint.Data) {
			key := r.GetKey(fixedPoint.X, fixedPoint.Y)
			if _, ok := nodes[key]; ok {
				continue
			}

			//fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minPoint.GetString(), newPoint.GetString(), fixedPoint.GetString(), math.Distance(minPoint, fixedPoint))

			nodeCounter++
			newNode := node.New(fixedPoint).WithID(nodeCounter)
			ok := r.tree.Add(minID, newNode)
			if !ok {
				fmt.Println("Could not add to tree")
			}

			nodes[key] = newNode
			lastNodeAdded = newNode
		}
	}

	return nodeCounter
}

func (r *rrt[T]) GetRandomPoint(world [][]T, try int) (*Point[T], bool) {
	if try%r.focusOnFinishEveryTry == 0 {
		return r.FinishPoint, true
	}

	x, y := r.GetRandomXY(world)

	if x > len(world) {
		x = len(world)
	}

	if y > len(world[0]) {
		y = len(world[0])
	}

	return &Point[T]{
		X:    x,
		Y:    y,
		Data: world[x][y],
	}, true
}

func (r *rrt[T]) GetRandomXY(world [][]T) (int, int) {
	x := rand.Int() % len(world)
	y := rand.Int() % len(world[0])

	return x, y
}

func (r *rrt[T]) GetFixedPoint(minDistancePoint *Point[T], newPoint *Point[T], world [][]T) *Point[T] {
	radian := minDistancePoint.Atan2To(newPoint)

	deltaX := builtInMath.Sin(radian) * float64(r.stepDistance)
	deltaY := builtInMath.Cos(radian) * float64(r.stepDistance)

	x := minDistancePoint.X + int(deltaX)
	y := minDistancePoint.Y + int(deltaY)

	if x < 0 {
		x = 0
	} else if x > len(world)-1 {
		x = len(world) - 1
	}

	if y < 0 {
		y = 0
	} else if y > len(world[0])-1 {
		y = len(world[0]) - 1
	}

	fixed := &Point[T]{
		X:    x,
		Y:    y,
		Data: world[x][y],
	}

	//fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minDistancePoint.GetString(), newPoint.GetString(), fixed.GetString(), Distance(minDistancePoint, fixed))

	return fixed
}

func (r *rrt[T]) GetKey(x int, y int) string {
	return fmt.Sprintf("%d-%d", x, y)
}
