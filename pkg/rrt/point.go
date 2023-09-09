package rrt

import (
	"fmt"
	"math"
)

type Coordinate struct {
	X float64
	Y float64
}

type Point[T any] struct {
	X    float64
	Y    float64
	Data T
}

func (p *Point[T]) DistanceTo(other *Point[T]) float64 {
	return math.Sqrt(((p.X - other.X) * (p.X - other.X)) + ((p.Y - other.Y) * (p.Y - other.Y)))
}

func (p *Point[T]) Atan2To(other *Point[T]) float64 {
	return math.Atan2(other.X-p.X, other.Y-p.Y)
}

func (p *Point[T]) Println() {
	fmt.Println(p.GetString())
}

func (p *Point[T]) GetString() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}
