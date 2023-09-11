package rrt

import "github.com/johnfercher/go-tree/tree"

type RRT[T any] interface {
	AddCollisionCondition(condition func(point T) bool)
	AddStopCondition(condition func(testPoint *Point[T], finish *Point[T]) bool)
	FindPath(start *Point[any], finish *Point[any], world [][]T) []*Point[T]
	GetTree() *tree.Tree[*Point[T]]
}
