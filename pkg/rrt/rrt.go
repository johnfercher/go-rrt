package rrt

import (
	"github.com/johnfercher/go-rrt/pkg/internal"
	"github.com/johnfercher/go-rrt/pkg/rrt/math"
)

type RRT[T any] interface {
	AddCollisionCondition(condition func(point T) bool)
	AddStopCondition(condition func(testPoint *math.Point[T], finish *math.Point[T]) bool)
	FindPath(start *math.Coordinate, finish *math.Coordinate, world [][]T) []*math.Point[T]
	FindPathAndSavePdf(start *math.Coordinate, finish *math.Coordinate, world [][]T, file string) []*math.Point[T]
}

func New[T any](stepDistance int, maxTries int, focusOnFinishEveryTry int) RRT[T] {
	return internal.New[T](stepDistance, maxTries, focusOnFinishEveryTry)
}
