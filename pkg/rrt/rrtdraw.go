package rrt

import (
	"fmt"
	"github.com/johnfercher/go-tree/node"
	"github.com/johnfercher/go-tree/tree"
	"github.com/jung-kurt/gofpdf"
)

type rrtDraw[T any] struct {
	inner              RRT[T]
	drawScale          float64
	collisionCondition func(point T) bool
	file               string
}

func NewRRTWithPDF[T any](stepDistance int, maxTries int, focusOnFinishEveryTry int, drawScale float64, file string) RRT[T] {
	return &rrtDraw[T]{
		inner:     New[T](stepDistance, maxTries, focusOnFinishEveryTry),
		drawScale: drawScale,
		file:      file,
	}
}

func (r *rrtDraw[T]) AddCollisionCondition(condition func(point T) bool) {
	r.collisionCondition = condition
	r.inner.AddCollisionCondition(condition)
}

func (r *rrtDraw[T]) AddStopCondition(condition func(testPoint *Point[T], finish *Point[T]) bool) {
	r.inner.AddStopCondition(condition)
}

func (r *rrtDraw[T]) GetTree() *tree.Tree[*Point[T]] {
	return r.inner.GetTree()
}

func (r *rrtDraw[T]) FindPath(start *Point[any], finish *Point[any], world [][]T) []*Point[T] {
	points := r.inner.FindPath(start, finish, world)
	tree := r.inner.GetTree()

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: float64(len(world)) * r.drawScale,
			Ht: float64(len(world[0])) * r.drawScale,
		},
	})

	pdf.AddPage()

	pdf.SetLineWidth(0.5)
	r.drawObstacles(world, pdf)
	r.drawInterestPoints(start, finish, pdf)

	pdf.SetLineWidth(0.1)
	pdf.SetDrawColor(0, 0, 0)
	r.drawTree(tree, pdf)

	pdf.SetLineWidth(0.5)
	r.drawPath(points, pdf)

	err := pdf.OutputFileAndClose(r.file)
	if err != nil {
		fmt.Println(err.Error())
	}

	return points
}

func (r *rrtDraw[T]) drawTree(tree *tree.Tree[*Point[T]], pdf *gofpdf.Fpdf) {
	root, _ := tree.GetRoot()
	r.drawNodeToParent(root, pdf, 0)
}

func (r *rrtDraw[T]) drawNodeToParent(node *node.Node[*Point[T]], pdf *gofpdf.Fpdf, depth int) {
	nodeData := node.GetData()
	nexts := node.GetNexts()
	for _, next := range nexts {
		//fmt.Printf("depth: %d, next: %d\n", depth, len(nexts))
		nextData := next.GetData()
		pdf.Circle(float64(nodeData.X)*r.drawScale, float64(nodeData.Y)*r.drawScale, 0.5, "")
		pdf.Line(float64(nodeData.X)*r.drawScale, float64(nodeData.Y)*r.drawScale, float64(nextData.X)*r.drawScale, float64(nextData.Y)*r.drawScale)
		r.drawNodeToParent(next, pdf, depth+1)
	}
}

func (r *rrtDraw[T]) drawPath(points []*Point[T], pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(255, 0, 0)
	for i := 0; i < len(points)-1; i++ {
		pdf.Circle(float64(points[i].X)*r.drawScale, float64(points[i].Y)*r.drawScale, 0.5, "")
		pdf.Line(float64(points[i].X)*r.drawScale, float64(points[i].Y)*r.drawScale, float64(points[i+1].X)*r.drawScale, float64(points[i+1].Y)*r.drawScale)
	}
}

func (r *rrtDraw[T]) drawObstacles(world [][]T, pdf *gofpdf.Fpdf) {
	for i, line := range world {
		for j, element := range line {
			if r.collisionCondition(element) {
				pdf.Circle(float64(i)*r.drawScale, float64(j)*r.drawScale, 2, "F")
			}
		}
	}
}

func (r *rrtDraw[T]) drawInterestPoints(start *Point[any], finish *Point[any], pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(0, 0, 255)
	pdf.Circle(float64(start.X)*r.drawScale, float64(start.Y)*r.drawScale, 2, "")
	pdf.SetDrawColor(0, 255, 0)
	pdf.Circle(float64(finish.X)*r.drawScale, float64(finish.Y)*r.drawScale, 2, "")
}
