package internal

import (
	"fmt"
	"github.com/johnfercher/go-rrt/pkg/rrt/math"
	"github.com/johnfercher/go-tree/tree"
	"github.com/jung-kurt/gofpdf"
	builtInMath "math"
	"math/rand"
)

type RapidlyExploringRandomTrees[T any] struct {
	CollisionCondition    func(point T) bool
	StopCondition         func(testPoint *math.Point[T], finish *math.Point[T]) bool
	StepDistance          int
	MaxTries              int
	FocusOnFinishEveryTry int
	DrawScale             float64
	StartPoint            *math.Point[T]
	FinishPoint           *math.Point[T]
}

func New[T any](stepDistance int, maxTries int, focusOnFinishEveryTry int) *RapidlyExploringRandomTrees[T] {
	return &RapidlyExploringRandomTrees[T]{
		StepDistance:          stepDistance,
		MaxTries:              maxTries,
		FocusOnFinishEveryTry: focusOnFinishEveryTry,
		DrawScale:             2.0,
	}
}

func (r *RapidlyExploringRandomTrees[T]) AddCollisionCondition(condition func(point T) bool) {
	r.CollisionCondition = condition
}

func (r *RapidlyExploringRandomTrees[T]) AddStopCondition(condition func(testPoint *math.Point[T], finish *math.Point[T]) bool) {
	r.StopCondition = condition
}

func (r *RapidlyExploringRandomTrees[T]) FindPath(start *math.Coordinate, finish *math.Coordinate, world [][]T) []*math.Point[T] {
	nodeCounter, tr := r.findPath(start, finish, world)

	pathNodes, _ := tr.Backtrack(nodeCounter)
	var points []*math.Point[T]
	for _, pathNode := range pathNodes {
		_, data := pathNode.Get()
		points = append(points, data)
	}

	return points
}

func (r *RapidlyExploringRandomTrees[T]) FindPathAndSavePdf(start *math.Coordinate, finish *math.Coordinate, world [][]T, file string) []*math.Point[T] {
	nodeCounter, tr := r.findPath(start, finish, world)

	pathNodes, _ := tr.Backtrack(nodeCounter)
	var points []*math.Point[T]
	for _, pathNode := range pathNodes {
		_, data := pathNode.Get()
		points = append(points, data)
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: float64(len(world)) * r.DrawScale,
			Ht: float64(len(world[0])) * r.DrawScale,
		},
	})

	pdf.AddPage()

	pdf.SetLineWidth(0.5)
	r.drawObstacles(world, pdf)
	r.drawInterestPoints(start, finish, pdf)

	pdf.SetLineWidth(0.1)
	pdf.SetDrawColor(0, 0, 0)
	r.drawTree(tr, pdf)

	pdf.SetLineWidth(0.5)
	r.drawPath(points, pdf)

	err := pdf.OutputFileAndClose(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	return points
}

func (r *RapidlyExploringRandomTrees[T]) drawTree(tr *tree.Tree[*math.Point[T]], pdf *gofpdf.Fpdf) {
	root, _ := tr.GetRoot()
	r.drawNodeToParent(root, pdf, 0)
}

func (r *RapidlyExploringRandomTrees[T]) drawNodeToParent(node *tree.Node[*math.Point[T]], pdf *gofpdf.Fpdf, depth int) {
	_, nodeData := node.Get()
	nexts := node.GetNexts()
	for _, next := range nexts {
		//fmt.Printf("depth: %d, next: %d\n", depth, len(nexts))
		_, nextData := next.Get()
		pdf.Circle(nodeData.X*r.DrawScale, nodeData.Y*r.DrawScale, 0.5, "")
		pdf.Line(nodeData.X*r.DrawScale, nodeData.Y*r.DrawScale, nextData.X*r.DrawScale, nextData.Y*r.DrawScale)
		r.drawNodeToParent(next, pdf, depth+1)
	}
}

func (r *RapidlyExploringRandomTrees[T]) drawPath(points []*math.Point[T], pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(255, 0, 0)
	for i := 0; i < len(points)-1; i++ {
		pdf.Circle(points[i].X*r.DrawScale, points[i].Y*r.DrawScale, 0.5, "")
		pdf.Line(points[i].X*r.DrawScale, points[i].Y*r.DrawScale, points[i+1].X*r.DrawScale, points[i+1].Y*r.DrawScale)
	}
}

func (r *RapidlyExploringRandomTrees[T]) drawObstacles(world [][]T, pdf *gofpdf.Fpdf) {
	for i, line := range world {
		for j, element := range line {
			if r.CollisionCondition(element) {
				pdf.Circle(float64(i)*r.DrawScale, float64(j)*r.DrawScale, 2, "F")
			}
		}
	}
}

func (r *RapidlyExploringRandomTrees[T]) drawInterestPoints(start *math.Coordinate, finish *math.Coordinate, pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(0, 0, 255)
	pdf.Circle(start.X*r.DrawScale, start.Y*r.DrawScale, 2, "")
	pdf.SetDrawColor(0, 255, 0)
	pdf.Circle(finish.X*r.DrawScale, finish.Y*r.DrawScale, 2, "")
}

func (r *RapidlyExploringRandomTrees[T]) findPath(start *math.Coordinate, finish *math.Coordinate, world [][]T) (int, *tree.Tree[*math.Point[T]]) {
	r.StartPoint = &math.Point[T]{
		X:    start.X,
		Y:    start.Y,
		Data: world[int(start.X)][int(start.Y)],
	}
	r.FinishPoint = &math.Point[T]{
		X:    finish.X,
		Y:    finish.Y,
		Data: world[int(finish.X)][int(finish.Y)],
	}

	nodes := make(map[string]*tree.Node[*math.Point[T]])
	var lastNodeAdded *tree.Node[*math.Point[T]]
	tr := tree.New[*math.Point[T]]()

	nodeCounter := 0
	node := tree.NewNode[*math.Point[T]](nodeCounter, r.StartPoint)
	tr.AddRoot(node)
	_, point := node.Get()
	key := r.GetKey(int(point.X), int(point.Y))
	nodes[key] = node

	lastNodeAdded = node

	try := 0
	minDistance := builtInMath.MaxFloat64
	var minNode *tree.Node[*math.Point[T]]

	for try < r.MaxTries {
		minDistance = builtInMath.MaxFloat64
		minNode = nil
		try++

		_, lastAdded := lastNodeAdded.Get()
		if r.StopCondition(lastAdded, r.FinishPoint) {
			break
		}

		newPoint, ok := r.GetRandomPoint(world, try)
		if !ok {
			continue
		}

		for _, point := range nodes {
			_, vector := point.Get()
			distance := math.Distance(vector, newPoint)
			if distance < minDistance {
				minDistance = distance
				minNode = point
			}
		}

		minID, minPoint := minNode.Get()
		fixedPoint := r.GetFixedPoint(minPoint, newPoint, world)

		if !r.CollisionCondition(fixedPoint.Data) {
			key := r.GetKey(int(fixedPoint.X), int(fixedPoint.Y))
			if _, ok := nodes[key]; ok {
				continue
			}

			//fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minPoint.GetString(), newPoint.GetString(), fixedPoint.GetString(), math.Distance(minPoint, fixedPoint))

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

func (r *RapidlyExploringRandomTrees[T]) GetRandomPoint(world [][]T, try int) (*math.Point[T], bool) {
	if try%r.FocusOnFinishEveryTry == 0 {
		return r.FinishPoint, true
	}

	x, y := r.GetRandomXY(world)

	if x > len(world) {
		x = len(world)
	}

	if y > len(world[0]) {
		y = len(world[0])
	}

	return &math.Point[T]{
		X:    float64(x),
		Y:    float64(y),
		Data: world[x][y],
	}, true
}

func (r *RapidlyExploringRandomTrees[T]) GetRandomXY(world [][]T) (int, int) {
	x := rand.Int() % len(world)
	y := rand.Int() % len(world[0])

	return x, y
}

func (r *RapidlyExploringRandomTrees[T]) GetFixedPoint(minDistancePoint *math.Point[T], newPoint *math.Point[T], world [][]T) *math.Point[T] {
	radian := math.Radian(minDistancePoint, newPoint)

	deltaX := builtInMath.Sin(radian) * float64(r.StepDistance)
	deltaY := builtInMath.Cos(radian) * float64(r.StepDistance)

	x := int(minDistancePoint.X) + int(deltaX)
	y := int(minDistancePoint.Y) + int(deltaY)

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

	fixed := &math.Point[T]{
		X:    float64(x),
		Y:    float64(y),
		Data: world[x][y],
	}

	//fmt.Printf("Min %s, New %s, Fix %s, D: %f\n", minDistancePoint.GetString(), newPoint.GetString(), fixed.GetString(), Distance(minDistancePoint, fixed))

	return fixed
}

func (r *RapidlyExploringRandomTrees[T]) GetKey(x int, y int) string {
	return fmt.Sprintf("%d-%d", x, y)
}
