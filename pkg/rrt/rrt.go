package rrt

import (
	"fmt"
	"github.com/johnfercher/go-tree/tree"
	"github.com/jung-kurt/gofpdf"
	"math"
	"math/rand"
)

type RRT[T any] struct {
	collisionCondition    func(point T) bool
	stopCondition         func(testPoint *Point[T], finish *Point[T]) bool
	stepDistance          float64
	maxTries              int
	focusOnFinishEveryTry int
	startPoint            *Point[T]
	finishPoint           *Point[T]
}

func New[T any](stepDistance float64, maxTries int, focusOnFinishEveryTry int) *RRT[T] {
	return &RRT[T]{
		stepDistance:          stepDistance,
		maxTries:              maxTries,
		focusOnFinishEveryTry: focusOnFinishEveryTry,
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
	nodeCounter, tr := r.findPath(start, finish, world)

	pathNodes, _ := tr.Backtrack(nodeCounter)
	var points []*Point[T]
	for _, pathNode := range pathNodes {
		_, data := pathNode.Get()
		points = append(points, data)
	}

	return points
}

func (r *RRT[T]) FindPathAndSavePdf(start *Coordinate, finish *Coordinate, world [][]T, file string) []*Point[T] {
	scale := 2.0
	nodeCounter, tr := r.findPath(start, finish, world)

	pathNodes, _ := tr.Backtrack(nodeCounter)
	var points []*Point[T]
	for _, pathNode := range pathNodes {
		_, data := pathNode.Get()
		points = append(points, data)
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: float64(len(world)) * scale,
			Ht: float64(len(world[0])) * scale,
		},
	})

	pdf.AddPage()

	for i, line := range world {
		for j, element := range line {
			if r.collisionCondition(element) {
				pdf.Circle(float64(i)*scale, float64(j)*scale, 2, "F")
			}
		}
	}

	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(0, 0, 255)
	pdf.Circle(start.X*scale, start.Y*scale, 2, "")
	pdf.SetDrawColor(0, 255, 0)
	pdf.Circle(finish.X*scale, finish.Y*scale, 2, "")
	pdf.SetDrawColor(255, 0, 0)
	for i := 0; i < len(points)-1; i++ {
		pdf.Circle(points[i].X*scale, points[i].Y*scale, 1, "")
		pdf.Line(points[i].X*scale, points[i].Y*scale, points[i+1].X*scale, points[i+1].Y*scale)
	}

	fmt.Println("pdf")
	err := pdf.OutputFileAndClose(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	return points
}

func (r *RRT[T]) findPath(start *Coordinate, finish *Coordinate, world [][]T) (int, *tree.Tree[*Point[T]]) {
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
	var minNode *tree.Node[*Point[T]]

	for try < r.maxTries {
		minDistance = math.MaxFloat64
		minNode = nil
		try++

		_, lastAdded := lastNodeAdded.Get()
		if r.stopCondition(lastAdded, r.finishPoint) {
			break
		}

		newPoint := r.getRandomPoint(world, try)
		for _, point := range nodes {
			_, vector := point.Get()
			distance := Distance(vector, newPoint)
			if distance < minDistance {
				minDistance = distance
				minNode = point
			}
		}

		minID, minPoint := minNode.Get()
		fixedPoint := r.getFixedPoint(minPoint, newPoint, world)

		if !r.collisionCondition(fixedPoint.Data) {
			key := r.getKey(fixedPoint)
			if _, ok := nodes[key]; ok {
				continue
			}

			fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minPoint.GetString(), newPoint.GetString(), fixedPoint.GetString(), Distance(minPoint, fixedPoint))

			nodeCounter++
			newNode := tree.NewNode(nodeCounter, fixedPoint)
			ok := tr.Add(minID, newNode)
			if !ok {
				fmt.Println("Could not add to tree")
			}

			nodes[key] = newNode
			lastNodeAdded = newNode
		}
	}

	return nodeCounter, tr
}

func (r *RRT[T]) getRandomPoint(world [][]T, try int) *Point[T] {
	if try%r.focusOnFinishEveryTry == 0 {
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

	fixed := &Point[T]{
		X:    float64(x),
		Y:    float64(y),
		Data: world[x][y],
	}

	//fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minDistancePoint.GetString(), newPoint.GetString(), fixed.GetString(), Distance(minDistancePoint, fixed))
	return fixed
}

func (r *RRT[T]) getKey(point *Point[T]) string {
	return fmt.Sprintf("%d-%d", int(point.X), int(point.Y))
}
